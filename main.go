package main

import (
	"errors"
	"fmt"
	"sync"
)

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	wg := sync.WaitGroup{}
	ec := make(chan error, 1)
	mc := make(chan string, len(nums))

	wg.Add(len(nums))
	for _, num := range nums {
		go addNums(num, &wg, ec, mc)
	}

	wg.Wait()
	close(ec)
	close(mc)

	if err := <-ec; err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < len(nums); i++ {
		select {
		case m := <-mc:
			fmt.Println("this is the message:", m)
		default:
			fmt.Println("df")
		}
	}
}

func addNums(n int, wg *sync.WaitGroup, ec chan error, mc chan string) {
	defer wg.Done()
	m := fmt.Sprintf("that's cool, %d + 2 = 12", n)
	if n+2 != 12 {
		m = fmt.Sprintf("%d + 2 does not equal 12", n)
	}
	if n+2 == 13 {
		select {
		case ec <- errors.New("adds up to 13, unlucky"):
		default:
		}
	}
	select {
	case mc <- m:
	default:
		fmt.Println("oops")
	}
}
