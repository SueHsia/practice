package main
import (
	"fmt"
)

type Myfloat float64

func(f Myfloat) abs() float64{
	if f<0{
		return float64(-f)
	}else if f==0{
		return 0
	}else{
		return float64(f)
	}
}

func main(){
	f:=Myfloat(-123)
	fmt.Println(f.abs())
}