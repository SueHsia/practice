package main

import "fmt"

func reverse(x,y string)(string,string){
	return y,x
}

func main(){
	a, b := reverse("a", "b")
	fmt.Println(a,b)
}