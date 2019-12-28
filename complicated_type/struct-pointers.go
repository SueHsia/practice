package main
import "fmt"

type Vertex struct{
	X int
	Y int
}

func main(){
	v:=Vertex{1,2}
	p:=&v
	t:=p
	v.Y=1e9
	p.X=20
	//*p=*t=v
	fmt.Println(v,*p,*t)
}