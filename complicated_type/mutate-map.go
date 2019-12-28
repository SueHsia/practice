package main
import "fmt"

func main(){
	m:=make(map[string]int)
	m["xiaxu"]=19
	m["zhaopan"]=20
	fmt.Println(m)
	m["xiaxu"]=21
	fmt.Println(m)
	delete(m,"xiaxu")
	fmt.Println(m)
	//如果 key为zhaopan 在 m 中，`v` 为 true 。否则， v 为 `false`,并且i为value值
	i,v:=m["zhaopan"]
	fmt.Println(i,v)
}