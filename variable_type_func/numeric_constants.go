package main

import (
	"fmt"
)
//二进制1左移100位，得出的十进制数
//1<<10就是2^10即1024
const big=1<<100
const small=big>>99

func needInt(x int)int{
	return x*10+1
}
func needFloat(x float64)float64{
	return x*0.1
}

func main(){
	// fmt.Println(big)
	// fmt.Println(small)
	fmt.Println(needInt(small))
	fmt.Println(needFloat(small))
	fmt.Println(needFloat(big))
	fmt.Println(needInt(big))
}