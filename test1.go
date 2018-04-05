package main
import "fmt"
import "time"
func f1(ch chan int) {
	time.Sleep(time.Second * 10)
	ch <- 1
}
func f2(ch chan int) {
	for {
		time.Sleep(time.Second * 9)
		ch <- 1
	}
}
func main3() {
	var ch1 = make(chan int)
	var ch2 = make(chan int)
	go f1(ch1)
	go f2(ch2)
	var val,val1 int
	for {
		select {
		case val = <-ch1:
			fmt.Println("The first case is selected."+string(val) )
		case val1 = <-ch2:
			fmt.Println(val1)
			fmt.Println("The second case is selected." )
		}
	}
}