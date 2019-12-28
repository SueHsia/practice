package main
import "fmt"

type IPAddr [4] byte

func(IP IPAddr) String() string{
	return fmt.Sprintf("%v.%v.%v.%v",IP[0],IP[1],IP[2],IP[3])
}

func main(){
	addres:=map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for n,a:=range addres{
		fmt.Printf("%v:%v\n",n,a.String())
	}
}