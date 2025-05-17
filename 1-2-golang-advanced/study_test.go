package __2_golang_advanced

import (
	"fmt"
	"testing"
	"time"
)

func TestChangeValue(t *testing.T) {
	num := 5
	fmt.Println("Before function call:", num) // 输出初始值

	changeValue(&num) // 传递num的地址给函数

	fmt.Println("After function call:", num) // 输出修改后的值
}

func TestDoubleSliceElements(t *testing.T) {
	num := []int{1, 2, 3, 4, 5}
	fmt.Println("Before function call:", num)
	doubleSliceElements(&num)
	fmt.Println("After function call:", num)
}

func TestOddAndEven(t *testing.T) {
	oddAndEven()
}

func TestTaskScheduler(t *testing.T) {
	tasks := []Task{
		{
			Name: "Task1",
			Fn: func() {
				time.Sleep(10 * time.Millisecond)
			},
		},
		{
			Name: "Task2",
			Fn: func() {
				time.Sleep(20 * time.Millisecond)
			},
		},
	}

	taskScheduler(tasks)
}

func TestShape(t *testing.T) {
	re := Rectangle{Width: 5, Height: 10}
	ci := Circle{Radius: 5}
	fmt.Println(re.Area())
	fmt.Println(ci.Area())
	fmt.Println(re.Perimeter())
	fmt.Println(ci.Perimeter())
}

func TestPerson(t *testing.T) {
	e := Employee{
		Person: Person{
			Name: "John",
			Age:  30,
		},
		EmployeeID: 123,
	}
	fmt.Println(e.Name)
}

func TestChannelCommunication(t *testing.T) {
	channelCommunication()
}

func TestBufferedChannel(t *testing.T) {
	bufferedChannel()
}

func TestMutexCounter(t *testing.T) {
	mutexCounter()
}

func TestAtomicCounter(t *testing.T) {
	atomicCounter()
}
