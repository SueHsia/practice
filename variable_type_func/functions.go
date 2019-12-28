package main
import (
	"fmt"
)

func add(x float32,y float32)float32{
	return x+y
}
func main(){
	fmt.Println(add(1.0,2.1))
}