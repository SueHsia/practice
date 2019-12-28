package main
import "fmt"

type Vertex struct{
	X,Y int
}

var m=map[string]Vertex{
	"xiaxu":{123,456},
	"mingjue":{456,789},
}

func main(){
	fmt.Println(m)
	dict:=make(map[string][]int)
	dict["xiaxu"]=append(dict["xiaxu"],1,2,3)
	fmt.Println(dict["xiaxu"][0])
}