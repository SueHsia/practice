package main

import "fmt"

func main(){
	i,j := 42,2701
	p:=&i
	fmt.Println("*p:",*p)
	*p=21
	fmt.Println("i:",i)

	p=&j
	x:=*p
	*p=*p/37
	fmt.Printf("原始j:%d,修改后的j:%d\n",x,j)
}
