package main
import "fmt"

func main(){
	var a []int
	prinSlice("a",a)
	for i:=0;i<5;i++{
		a=append(a,i)
	}
	prinSlice("a",a)
	a=append(a,8,9,10)
	prinSlice("a",a)
}
func prinSlice(s string,x []int){
	fmt.Printf("%s len=%d cap=%d %v\n",s,len(x),cap(x),x)
}