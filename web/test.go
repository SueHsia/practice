package main

import (
	"fmt"
	"strconv"
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
	a = fmt.Sprintf("abc-%v-abc", tm.Format("2006-01-02 15:04:05"))
	fmt.Println(a)
	// second := time.Now().Unix()
	second1 := strconv.FormatInt(time.Now().Unix(), 10)
	fmt.Printf("%T\n", second1)
}
