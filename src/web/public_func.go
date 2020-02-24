package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func enbase(str string) string {
	encodestr := []byte(str)
	for i := 0; i < len(encodestr); i++ {
		encodestr[i] += 124
	}
	encodeString := base64.StdEncoding.EncodeToString(encodestr)
	return encodeString
}

func debase(str string) string {
	decodeBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < len(decodeBytes); i++ {
		decodeBytes[i] -= 124
	}
	return string(decodeBytes)
}

func errorHandle(err error, w http.ResponseWriter) {
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func validate(str string) (username string, userid int, is_valid bool) {
	content := strings.Split(debase(str), ",")
	token_time, _ := strconv.ParseInt(content[len(content)-1], 10, 64)
	cur_time := time.Now().Unix()
	if cur_time-token_time >= 7200 {
		is_valid = false
	} else {
		is_valid = true
	}
	username = content[0]
	userid, _ = strconv.Atoi(content[1])
	return username, userid, is_valid
}

// func main() {
// 	tm1 := time.Now().Unix()
// 	input := enbase("hello")
// 	fmt.Println(input)
// 	output := debase(input)
// 	fmt.Println(output)
// 	tm2 := time.Now().Unix()
// 	fmt.Println("间隔", tm2+100-tm1)
// 	fmt.Println(debase("9OXd9PGoraitsbO0r7Owsq60"))
// }
