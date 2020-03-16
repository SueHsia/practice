package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/garyburd/redigo/redis"
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

func errorPrint(err error) {
	if err != nil {
		fmt.Println(err)
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

func mysql_redis(goodsId int, viewCount int) string {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
	}
	defer c.Close()

	// db, err = sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/lost_and_found?charset=utf8")
	errorPrint(err)
	curName := fmt.Sprintf("viewCount-%v", goodsId)

	is_key_exit, err := redis.Bool(c.Do("EXISTS", curName))

	if is_key_exit == false {
		curResult := fmt.Sprintf("%v-%v", viewCount, time.Now().Unix())
		_, err = c.Do("SET", curName, curResult)
	} else {
		curResult, _ := redis.String(c.Do("GET", curName))
		result := strings.Split(curResult, "-")
		viewCount, err = strconv.Atoi(result[0])
		newResult := fmt.Sprintf("%v-%v", viewCount+1, result[1])
		_, err = c.Do("SET", curName, newResult)
	}
	redisCount, _ := redis.String(c.Do("GET", curName))
	return redisCount
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
// 	mysql_redis(1)
// 	c, _ := redis.Dial("tcp", "127.0.0.1:6379")
// 	curResult, _ := redis.String(c.Do("GET", "view_count"))
// 	fmt.Println(curResult)
// }
