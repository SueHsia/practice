package main
import(
	"fmt"
	"math/rand"
	"time"
)
func main(){
	fmt.Println(rand.Intn(10))			//每次返回一样的数
	rand.Seed(time.Now().Unix())	//将时间作为随机数种子，则生成不同的随机数
	fmt.Println(rand.Intn(10))			//每次返回不同的数
	fmt.Println(time.Now().UnixNano())
}