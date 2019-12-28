package main

import (
	"fmt"
	"strconv"
)

type Person struct{
	Name string
	Age int
}

func (p Person) String() string{
	if p.Age<18{
		return p.Name+"的年龄是"+strconv.Itoa(p.Age)+"岁"
	}else{
		return fmt.Sprintf("%v的年龄是%v岁",p.Name,p.Age)
	}
}
func main(){
	p:=Person{"夏旭",23}
	fmt.Println(p.String())
}