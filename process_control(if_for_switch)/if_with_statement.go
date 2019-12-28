package main

import "fmt"
import "math"

func pow(x,n,lim float64)float64{
	if math.Pow(x,n)<=lim {
		return math.Pow(x,n)
	}else{
		return lim
	}
}

func main(){
	fmt.Println(pow(2,2,3),pow(3,3,4))
}