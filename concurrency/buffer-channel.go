package main

import "fmt"

func main() {
	// 2为参数，表示缓冲区长度。
	// 向缓冲 channel 发送数据的时候，只有在缓冲区满的时候才会阻塞。当缓冲区清空的时候接受阻塞。
	c := make(chan int, 2)
	c <- 1
	c <- 2
	// c <- 3
	fmt.Println(<-c)
	fmt.Println(<-c)
}
