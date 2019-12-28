package main

import (
	"sort"
	"fmt"
)
func main(){
	// 1.对整形数组排序
	s1 := []int{1,3,43,3,5,6,32,56,7,8,53,}
	sort.Ints(s1[:])  //切片本质上是对底层数组的一个view，不能直接修改view，通过视图来修改
	fmt.Println(s1) //[1 3 3 5 6 7 8 32 43 53 56]

	// 2.对字符串类型数组进行排序
	s2:=[]string{"mashiro","satori","hatsunemiku"}
	sort.Strings(s2[:])
	fmt.Println(s2) //[hatsunemiku mashiro satori]

	// 3.对浮点数类型数组进行排序
	s3:=[]float64{12.3,31.31,312.2,1.1,43,2}
	sort.Float64s(s3[:])
	fmt.Println(s3) //[1.1 2 12.3 31.31 43 312.2]

	// 4.查找一个int在数组中的索引
	s4:=[]int{2,3,452,1,3,4,5,67}
	index:=sort.SearchInts(s4[:],3)
	fmt.Println(index) //如果有多个，只返回第一个

	// 5.查找一个string在数组中的位置
	s5:=[]string{"mashiro","satori","hatsunemiku"}
	index1 := sort.SearchStrings(s5[:],"mashiro")
	fmt.Println(index1) //0
}