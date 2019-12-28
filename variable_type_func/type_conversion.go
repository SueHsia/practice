package main

import (
	"fmt"
	"math"
)

func main(){
	x,y :=3,7
	i:=float64(x*x+y*y)
	f:=float64(math.Sqrt(i))
	z:=uint(f)
	fmt.Println(x,y,z)
	var a,b int = 3,7
	var c =math.Sqrt(float64(a*a+b*b))
	var d =int(c)
	fmt.Println(a,b,d)
}