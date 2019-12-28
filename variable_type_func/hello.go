package main
import "fmt"
var a="xiaxu"
var b="aaabc"
var e, f = 123, "hello"
var g = 111

func main(){
	var Book1 book       /* 声明 Book1 为 Books 类型 */

	/* book 1 描述 */
	Book1.title = "Go 语言"
	Book1.author = "www.runoob.com"
	Book1.subject = "Go 语言教程"
	Book1.book_id = 6495407
	
	pringBook(Book1)
	println("=======================")
	// print("aaa")
	// println("bbb")
	var numbers []int
	printSlice(numbers)

	numbers=append(numbers,1)
	printSlice(numbers)

	numbers=append(numbers,2,3,4)
	printSlice(numbers)

	fmt.Println(numbers[0:2])
	println("=======================")
	var agemap map[string]int //括号内是key类型，外面是value类型
	agemap = make(map[string]int)
	agemap [ "France" ] = 112
    agemap [ "Italy" ] = 223
    agemap [ "Japan" ] = 444
	agemap [ "India " ] = 555
	
	for country := range agemap {
        println(country, "年龄是", agemap[country])
	}
	
	var i int = 15
	// println(uint64(i))
    fmt.Printf("%d 的阶乘是 %d\n", i, Factorial(uint64(i)))
	println("==================================")
	var new_phone Phone
	new_phone=new(iPhone)
	new_phone.call()
	new_phone=new(Android)
	new_phone.call()
}

func Factorial(n uint64) uint64{
    if (n > 0) {
        result:= n * Factorial(n-1)
        return result
    }
    return 1
}


func printSlice(x []int){
	fmt.Printf("len=%d cap=%d slice=%v\n",len(x),cap(x),x)
 }

func pringBook(Books book){
	println("title:",Books.title)
	println("author:",Books.author)
	println("subject:",Books.subject)
	println("book_id:",Books.book_id)
}

type Phone interface{
	call()
}

type Android struct{

}

type iPhone struct{

}

func (xiaomi Android) call(){
	println("i am xiaomi")
}

func (iphone8 iPhone)call(){
	println("i am iphone8")
}

type book struct{
	title string
	author string
	subject string
	book_id int
}

func max(a int,b string) int {
	println(a,b)
	if a>1{
		return 1
	}else{
		return 0
	}
}

func getSequence() func() int {
	i:=0
	return func() int {
	   i+=1
	  return i  
	}
 }