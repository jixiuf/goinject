package inject

import (
	"fmt"
	"testing"
)

func TestInject(t *testing.T) {
	Register("run", &Runner{}) //以run为名注册 Runner
	Register("fly", &Bird{})
	people := &People{}
	err := Inject(people) //注入带有Tag:"auto"的所有对象
	if err != nil {
		fmt.Println(err)
	}

	if people.Run == nil {
		t.Error("people.Run注入失败")
	} else {
		people.Run.Run()
	}
	if people.fly == nil {
		t.Error("people.fly注入失败")
	} else {
		people.fly.Fly()
	}

}

func TestInjectAll(t *testing.T) {
	Register("run", &Runner{}) //以run为名注册 Runner
	Register("fly", &Bird{})
	people := &People{}
	Register("people", people)
	InjectAll() //注入带有Tag:"auto"的所有对象

	if people.Run == nil {
		t.Error("people.Run注入失败")
	} else {
		people.Run.Run()
	}
	if people.fly == nil {
		t.Error("people.fly注入失败")
	} else {
		people.fly.Fly()
	}

}

type IRun interface {
	Run()
}
type IFly interface {
	Fly()
}
type Runner struct {
}

func (r Runner) Run() {
	fmt.Println("run")
}

type Bird struct {
}

func (r Bird) Fly() {
	fmt.Println("bird.Fly()")
}

type People struct {
	Run IRun `auto:"run"`
	fly IFly `auto:"fly"` // 私有变量也可以自动注入
}
