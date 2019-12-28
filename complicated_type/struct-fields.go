package main
import "fmt"

//结构体中，首字母小写的属性为private，首字母大写的属性为public，但是如果在同一个GO文件内部则大小写都能访问，在不同文件内则小写不能访问。
type Vertex struct{
	x int
	y int
	X int
	Y int
}

func main(){
	//下面这句话等同于 v:=&Vertex{1,2,3,4}
	v:=Vertex{1,2,3,4}
	v.x=9
	fmt.Println(v.x,v.y,v.X,v.Y)
}