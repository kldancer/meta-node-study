package __2_golang_advanced

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/*
指针

题目1 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
考察点 ：指针的使用、值传递与引用传递的区别。
*/

func changeValue(p *int) {
	*p += 10
}

/*
题目2 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
考察点 ：指针运算、切片操作。
*/
func doubleSliceElements(slice *[]int) {
	for i := range *slice {
		(*slice)[i] *= 2
	}
}

/*
Goroutine

题目1 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
考察点 ： go 关键字的使用、协程的并发执行。
*/

func oddAndEven() {
	oddCh := make(chan struct{})
	evenCh := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)

	// 打印奇数的 goroutine
	go func() {
		defer wg.Done()
		for i := 1; i <= 9; i += 2 {
			<-oddCh
			fmt.Println(i)
			evenCh <- struct{}{}
		}
	}()

	// 打印偶数的 goroutine
	go func() {
		defer wg.Done()
		for i := 2; i <= 10; i += 2 {
			<-evenCh
			fmt.Println(i)
			if i < 10 {
				oddCh <- struct{}{}
			}
		}
	}()

	// 启动从 1（奇数）开始
	oddCh <- struct{}{}

	// 等待两边都完成
	wg.Wait()

	fmt.Println("Done")
}

/*
题目2 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
考察点 ：协程原理、并发任务调度。
*/

type Task struct {
	Fn       func() // 任务函数
	Name     string // 任务名称（用于标识）
	Duration int64  // 执行时间，单位为纳秒
}

func taskScheduler(tasks []Task) {
	var wg sync.WaitGroup

	for i := range tasks {
		wg.Add(1)
		go func(task Task) {
			defer wg.Done()
			start := time.Now()
			task.Fn()
			duration := time.Since(start).Nanoseconds()
			task.Duration = duration
			fmt.Printf("任务 [%s] 执行完成，耗时 %d 纳秒\n", task.Name, duration)
		}(tasks[i])
	}

	wg.Wait()
	fmt.Println("所有任务执行完成")
}

/*
面向对象,

题目1 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
考察点 ：接口的定义与实现、面向对象编程风格。
*/

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width, Height float64
}

type Circle struct {
	Radius float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (c Circle) Area() float64 {
	return 3.14 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * 3.14 * c.Radius
}

/*
2.题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
考察点 ：组合的使用、方法接收者。
*/

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeID int
}

func (e Employee) PrintInfo() {
	fmt.Printf("员工信息：Name=%s, Age=%d, EmployeeID=%d\n", e.Name, e.Age, e.EmployeeID)
}

/*
Channel
1. 编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
考察点 ：通道的基本使用、协程间通信。
*/

func channelCommunication() {
	ch := make(chan int)

	go func() {
		for i := 1; i <= 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	go func() {
		for num := range ch {
			fmt.Println("接收到数字:", num)
		}
	}()

	time.Sleep(time.Second * 2)
	fmt.Println("通信完成")
}

/*
2. 实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
考察点 ：通道的缓冲机制。
*/

func bufferedChannel() {
	ch := make(chan int, 10)

	go func() {
		for i := 1; i <= 100; i++ {
			ch <- i
			fmt.Printf("生产者发送：%d\n", i)
		}
		close(ch) // 发送完成后关闭通道
	}()

	go func() {
		for num := range ch {
			fmt.Printf("消费者接收到：%d\n", num)
		}
		fmt.Println("所有数字已消费完毕")
	}()

	time.Sleep(2 * time.Second)
}

/*
锁机制
1. 题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点 ： sync.Mutex 的使用、并发数据安全。
*/

func mutexCounter() {
	var (
		counter int
		mutex   sync.Mutex
		wg      sync.WaitGroup
	)

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mutex.Lock()
				counter++
				mutex.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("最终计数器的值: %d\n", counter)
}

/*
2. 题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点 ：原子操作、并发数据安全。
*/

func atomicCounter() {
	var counter int64
	var wg sync.WaitGroup

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}

	wg.Wait()
	fmt.Printf("最终计数器的值: %d\n", counter)
}
