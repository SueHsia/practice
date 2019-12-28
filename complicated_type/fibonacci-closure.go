package main

import "fmt"

// fibonacci 函数会返回一个返回 int 的函数。
func fibonacci() func()int{
	x,y:=0,1
	return func()int{
		temp:=x
		x,y=y,(x+y)
		return temp
	}
}
func main(){
	f:=fibonacci()
	for i:=0;i<10;i++{
		fmt.Println(f())
	}
}