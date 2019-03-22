package inject

// 反射实现依赖注入
// 用法见单元测试
import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

const (
	InjectorTag = "auto"
)

var objs sync.Map = sync.Map{}

// Register 注册对象
// 注册进来的对象尽量是指针类型，这样才可以修改它的属性，为它注入
func Register(name string, v interface{}) {
	objs.Store(name, reflect.ValueOf(v))
}

// // AutoRegister 注册对象
// func AutoRegister(v interface{}) {
// 	rv := reflect.ValueOf(v)
// 	Register(rv.Type().String(), rv)
// }

func get(key string) (value reflect.Value, ok bool) {
	i, ok := objs.Load(key)
	if !ok {
		return
	}
	return i.(reflect.Value), true
}

// Get 获取注册对象
func Get(key string) interface{} {
	value, ok := get(key)
	if !ok {
		return nil
	}
	return value.Interface()
}

// Remove 删除注册对象
func Remove(key string) {
	objs.Delete(key)
}

type ErrMissedInjectField struct {
	InjectorTag string
}

func (e ErrMissedInjectField) Error() string {
	return fmt.Sprintf("MirredInjectField:%s", e.InjectorTag)
}

func Inject(v interface{}) (err error) {
	value := reflect.ValueOf(v)
	return inject(value)
}

func InjectAll() {
	f := func(key interface{}, v interface{}) bool {
		value := v.(reflect.Value)
		err := inject(value)
		if err != nil {
			fmt.Println(err)
		}
		return true
	}
	objs.Range(f)

}
func inject(value reflect.Value) (err error) {
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if !value.CanSet() {
		return errors.New(value.Type().String() + " CanSet() is false,please give an pointer of " + value.Type().String())
	}

	for i := 0; i < value.NumField(); i++ {
		name := value.Type().Field(i).Tag.Get(InjectorTag)
		temp, ok := get(name)
		if ok {
			field := value.Field(i)
			if field.CanSet() {
				field.Set(temp)
			} else {
				field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
				field.Set(temp)
			}
		} else {
			err = ErrMissedInjectField{InjectorTag: name}
		}
	}
	return err
}
