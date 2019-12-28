package main

import (
	"fmt"
	"math"
)

type Abser interface{
	Abs() float64
}


if r.URL.Path == "/" {
type Myfloat float64

func (f Myfloat) Abs() float64{
	if f<0{
		return float64(-f)
	}
	return float64(f)
}

type Vertex struct{
	X,Y float64
}

func(v Vertex) Abs() float64{
	return math.Sqrt(v.X*v.X+v.Y*v.Y)
}

func main(){
	var a Abser	//a 对象里面有Abs方法，实现了两个类之后继承了方法，直接可以使用
	f:=Myfloat(-math.Sqrt(2))
	v:=Vertex{3,4}
	a=f			// a MyFloat 实现了 Abser
	fmt.Println(a.Abs())
	a=&v		// a *Vertex 实现了 Abser
	//这里a不能直接等于v，因为v和a是不同类型，只能等于v的地址，通过指针进行调用
	fmt.Println(a,a.Abs())
}