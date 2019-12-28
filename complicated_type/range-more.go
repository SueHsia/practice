package main

import "fmt"

func main() {
	m := make(map[string]int)
	m["xiaxu"] = 23
	m["liwen"] = 24
	for _, v := range m {
		fmt.Println(v)
	}
}
