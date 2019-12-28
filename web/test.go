package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	a := "abc dce"
	b := strings.Replace(a, " ", "", -1)
	fmt.Println(b)
	a = "abc"
	fmt.Println(a)
	a = ""
	if a == "" {
		fmt.Println("nil")
	}
	tm := time.Unix(time.Now().Unix(), 0)
	fmt.Println(tm.Format("2006-01-02 15:04:05"))
}
