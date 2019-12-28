package main

import "fmt"

var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}

func main(){
	//range相当于迭代，i为键值，v表示value值
	for i,v:=range pow{
		fmt.Printf("key=%d,value=%d\n",i,v)
	}
}