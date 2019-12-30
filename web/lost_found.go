package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/Go-SQL-Driver/MySQL"
)

const (
	upload_path = "."
)

type dic struct {
	Flag        int
	Msg         string
	Token       string
	Username    string
	Total_count int
	Data        []Goodsinfo
	Total_page  int
}

type Goodsinfo struct {
	Goodsid     int
	Goodsname   string
	Address     string
	Pic         string
	Phone       string
	Des         string
	Userid      int
	Create_time string
	Update_time string
	View_count  int
	Status      int
	Is_return   int
}

type Userinfo struct {
	Userid      int
	Username    string
	Password    string
	Create_time string
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	res_page := r.URL.Query()["page"][0]
	if res_page == "" {
		res_page = "1"
	}
	// fmt.Println("page", res_page)
	dict := dic{}
	tem_goods := Goodsinfo{}
	page, _ := strconv.Atoi(res_page)
	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/lost_and_found?charset=utf8")
	errorHandle(err, w)
	rows, err := db.Query("select * from lost_goods where status = ? limit ?,?", 1, (int(page)-1)*8, int(page)*8)
	errorHandle(err, w)
	db.QueryRow("select count(*) from lost_goods where status = ?", 1).Scan(&dict.Total_count)
	dict.Total_page = int(math.Ceil(float64(dict.Total_count) / 8.0))
	errorHandle(err, w)
	for rows.Next() {
		rows.Scan(&tem_goods.Goodsid, &tem_goods.Goodsname, &tem_goods.Address, &tem_goods.Pic, &tem_goods.Phone, &tem_goods.Des, &tem_goods.Userid, &tem_goods.Create_time, &tem_goods.Update_time, &tem_goods.View_count, &tem_goods.Status, &tem_goods.Is_return)
		dict.Data = append(dict.Data, tem_goods)
	}
	dict_json, _ := json.Marshal(dict)
	result := string(dict_json)
	fmt.Fprint(w, result)
}

func showPicHandle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	file, err := os.Open(upload_path + req.URL.Path)
	errorHandle(err, w)
	defer file.Close()
	buff, err := ioutil.ReadAll(file)
	errorHandle(err, w)
	w.Write(buff)
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	dict := dic{}
	tem_data := Userinfo{}
	var result string
	if r.Method == "POST" {
		r.ParseForm()
		username := r.Form["username"][0]
		password := r.Form["password"][0]
		db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/lost_and_found?charset=utf8")
		errorHandle(err, w)
		rows := db.QueryRow("select * from pi_user where username = ?", username)
		rows.Scan(&tem_data.Userid, &tem_data.Username, &tem_data.Password, &tem_data.Create_time)
		if tem_data.Username == "" {
			dict.Msg = "用户名或密码错误"
			dict.Flag = 0
		} else {
			new_pass := md5V(password)
			if new_pass != tem_data.Password {
				dict.Msg = "用户名或密码错误"
				dict.Flag = 0
			} else {
				dict.Token = enbase([]byte(username + string(time.Now().Unix())))
				dict.Flag = 1
			}
			// dict.data = append(dict.data, tem_data)
			// fmt.Println(tem_data.id, tem_data.username, tem_data.password, tem_data.create_time)

			dict.Username = tem_data.Username
			// io.WriteString(w, string(dict_json))
		}
		dict_json, _ := json.Marshal(dict)
		result = string(dict_json)
		db.Close()
	}
	fmt.Fprintf(w, string(result))
}
func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/lost_and_found?charset=utf8")
	errorHandle(err, w)
	dict := dic{}
	tem_data := Userinfo{}
	var result string
	if r.Method == "POST" {
		r.ParseForm()
		username := strings.Replace(r.Form["username"][0], " ", "", -1)
		password := strings.Replace(r.Form["password"][0], " ", "", -1)
		repassword := strings.Replace(r.Form["repassword"][0], " ", "", -1)
		fmt.Println(username, password)
		rows := db.QueryRow("select * from pi_user where username = ?", username)
		rows.Scan(&tem_data.Userid, &tem_data.Username, &tem_data.Password, &tem_data.Create_time)
		if password != repassword {
			fmt.Println("111")
			dict.Msg = "密码不一致，注册失败！"
			dict.Flag = 0
		} else if tem_data.Username != "" {
			fmt.Println("222")
			dict.Msg = "用户已存在，注册失败"
			dict.Flag = 0
		} else {
			password = md5V(password)
			tm := time.Unix(time.Now().Unix(), 0)
			create_time := tm.Format("2006-01-02 15:04:05")
			stmt, err := db.Prepare("INSERT pi_user SET username=?,password=?,create_time=?")
			errorHandle(err, w)
			stmt.Exec(username, password, create_time)
			dict.Msg = "注册成功！"
			dict.Flag = 1
		}
		dict_json, err := json.Marshal(dict)
		errorHandle(err, w)
		result = string(dict_json)
		fmt.Println(result)
	}
	db.Close()
	fmt.Fprintf(w, result)
}
func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/media/", showPicHandle)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
