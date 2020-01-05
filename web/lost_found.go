package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/gorilla/mux"
)

const (
	upload_path = "."
)

var (
	db *sql.DB
)

type dic struct {
	Flag        int
	Msg         string
	Token       string
	Username    string
	Userid      int
	Total_count int
	Data        []Goodsinfo
	Total_page  int
	Is_edit     int
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
	dict := dic{}
	tem_goods := Goodsinfo{}
	page, _ := strconv.Atoi(res_page)
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
			dict.Username = tem_data.Username
			dict.Userid = tem_data.Userid
		}
		dict_json, _ := json.Marshal(dict)
		result = string(dict_json)
	}
	fmt.Fprintf(w, string(result))
}
func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	dict := dic{}
	tem_data := Userinfo{}
	var result string
	if r.Method == "POST" {
		r.ParseForm()
		username := strings.Replace(r.Form["username"][0], " ", "", -1)
		password := strings.Replace(r.Form["password"][0], " ", "", -1)
		repassword := strings.Replace(r.Form["repassword"][0], " ", "", -1)
		rows := db.QueryRow("select * from pi_user where username = ?", username)
		rows.Scan(&tem_data.Userid, &tem_data.Username, &tem_data.Password, &tem_data.Create_time)
		if password != repassword {
			dict.Msg = "密码不一致，注册失败！"
			dict.Flag = 0
		} else if tem_data.Username != "" {
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
	}
	fmt.Fprintf(w, result)
}

func publish(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "multipart/form-data")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	dict := dic{}
	tem_user := Userinfo{}
	tem_goods := Goodsinfo{}
	var ext string
	r.ParseForm()
	uploadFile, handle, err := r.FormFile("file")
	errorHandle(err, w)
	var result string
	if r.Method == "POST" {
		tm := time.Unix(time.Now().Unix(), 0)
		tem_goods.Address = r.Form["address"][0]
		tem_user.Username = r.Form["username"][0]
		tem_goods.Goodsname = r.Form["name"][0]
		tem_goods.Userid, _ = strconv.Atoi(r.Form["userid"][0])
		tem_goods.Phone = r.Form["phone"][0]
		tem_goods.Des = r.Form["des"][0]
		name := strings.Split(handle.Filename, ".")
		ext = strings.ToLower(name[len(name)-1])
		fileDir := fmt.Sprintf("./media/file/%v/", tm.Format("2006-01-02"))
		second := strconv.FormatInt(time.Now().Unix(), 10)
		filename := fileDir + second + "." + ext
		//保存图片
		os.Mkdir(fileDir, 0777)
		saveFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
		errorHandle(err, w)
		io.Copy(saveFile, uploadFile)
		defer uploadFile.Close()
		defer saveFile.Close()
		tem_goods.Pic = filename[1:]
		stmt, err := db.Prepare("INSERT lost_goods SET name=?,address=?,pic=?,phone=?,des=?,user_id=?,create_time=?,update_time=?,view_count=?,status=?,is_return=?")
		errorHandle(err, w)
		stmt.Exec(tem_goods.Goodsname, tem_goods.Address, tem_goods.Pic, tem_goods.Phone, tem_goods.Des, tem_goods.Userid, tm.Format("2006-01-02 15:04:05"), tm.Format("2006-01-02 15:04:05"), 0, 1, 0)
		dict.Msg = "发布成功"
		dict.Flag = 1
		// 上传图片成功
		dict_json, err := json.Marshal(dict)
		errorHandle(err, w)
		result = string(dict_json)
	}
	fmt.Fprintf(w, result)
}

func detail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	if r.Method == "GET" {
		tem_goods := Goodsinfo{}
		goods_id, _ := strconv.Atoi(r.URL.Query()["goods_id"][0])
		rows := db.QueryRow("select * from lost_goods where id = ?", goods_id)
		rows.Scan(&tem_goods.Goodsid, &tem_goods.Goodsname, &tem_goods.Address, &tem_goods.Pic, &tem_goods.Phone, &tem_goods.Des, &tem_goods.Userid, &tem_goods.Create_time, &tem_goods.Update_time, &tem_goods.View_count, &tem_goods.Status, &tem_goods.Is_return)
		stmt, err := db.Prepare("update lost_goods set view_count=? where id=?")
		errorHandle(err, w)
		stmt.Exec(tem_goods.View_count+1, tem_goods.Goodsid)
		dict := dic{}
		dict.Data = append(dict.Data, tem_goods)
		dict_json, err := json.Marshal(dict)
		errorHandle(err, w)
		result := string(dict_json)
		fmt.Fprintf(w, result)
	}
	if r.Method == "PUT" {
		tem_goods := Goodsinfo{}
		dict := dic{}
		do_return := r.Form["return_goods"][0]
		userid, _ := strconv.Atoi(r.Form["userid"][0])
		goods_id := r.Form["goods_id"][0]
		db.QueryRow("select user_id,is_return form lost_goods where id=?", goods_id).Scan(tem_goods.Userid, tem_goods.Is_return)
		if userid == tem_goods.Userid && do_return == "1" {
			stmt, err := db.Prepare("update lost_goods set is_return=? where id=?")
			errorHandle(err, w)
			stmt.Exec(1, goods_id)
			dict.Flag = 1
			dict.Msg = "归还成功"
		} else if userid == tem_goods.Userid && tem_goods.Is_return == 0 {
			dict.Msg = "失物可编辑"
			dict.Is_edit = 1
		} else {
			dict.Msg = "用户无权编辑"
			dict.Is_edit = 0
		}
		dict_json, err := json.Marshal(dict)
		errorHandle(err, w)
		result := string(dict_json)
		fmt.Fprintf(w, result)
	}
	if r.Method == "DELETE" {
		tem_goods := Goodsinfo{}
		dict := dic{}
		delete_code := r.Form["delete_code"][0]
		userid, _ := strconv.Atoi(r.Form["userid"][0])
		goods_id := r.Form["goods_id"][0]
		db.QueryRow("select user_id,status form lost_goods where id=?", goods_id).Scan(tem_goods.Userid, tem_goods.Status)
		if userid == tem_goods.Userid && delete_code == "1" {
			stmt, err := db.Prepare("update lost_goods set status=? where id=?")
			errorHandle(err, w)
			stmt.Exec(0, goods_id)
			dict.Flag = 1
			dict.Msg = "删除成功"
		} else {
			dict.Msg = "删除失败"
		}
		dict_json, err := json.Marshal(dict)
		errorHandle(err, w)
		result := string(dict_json)
		fmt.Fprintf(w, result)
	}
}

func edit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	dict := dic{}
	tem_goods := Goodsinfo{}
	r.ParseForm()
	if r.Method == "GET" {
		var result string
		goods_id, _ := strconv.Atoi(r.Form["goods_id"][0])
		rows := db.QueryRow("select * from lost_goods where id = ?", goods_id)
		rows.Scan(&tem_goods.Goodsid, &tem_goods.Goodsname, &tem_goods.Address, &tem_goods.Pic, &tem_goods.Phone, &tem_goods.Des, &tem_goods.Userid, &tem_goods.Create_time, &tem_goods.Update_time, &tem_goods.View_count, &tem_goods.Status, &tem_goods.Is_return)
		dict.Data = append(dict.Data, tem_goods)
		if tem_goods.Is_return == 1 {
			dict.Msg = "失物已归还禁止编辑"
			dict.Flag = 0
		} else {
			dict.Flag = 1
		}
		dict_json, err := json.Marshal(dict)
		errorHandle(err, w)
		result = string(dict_json)
		fmt.Fprintf(w, result)
	}
	if r.Method == "POST" {
		var ext string
		uploadFile, handle, err := r.FormFile("file")
		errorHandle(err, w)
		var result string
		tm := time.Unix(time.Now().Unix(), 0)
		goods_id, _ := strconv.Atoi(r.Form["goods_id"][0])
		tem_goods.Address = r.Form["address"][0]
		tem_goods.Goodsname = r.Form["name"][0]
		tem_goods.Phone = r.Form["phone"][0]
		tem_goods.Des = r.Form["des"][0]
		name := strings.Split(handle.Filename, ".")
		ext = strings.ToLower(name[len(name)-1])
		fileDir := fmt.Sprintf("./media/file/%v/", tm.Format("2006-01-02"))
		second := strconv.FormatInt(time.Now().Unix(), 10)
		filename := fileDir + second + "." + ext
		//保存图片
		os.Mkdir(fileDir, 0777)
		saveFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
		errorHandle(err, w)
		io.Copy(saveFile, uploadFile)
		defer uploadFile.Close()
		defer saveFile.Close()
		tem_goods.Pic = filename[1:]
		stmt, err := db.Prepare("update lost_goods set name=?,address=?,pic=?,phone=?,des=?,update_time=? where id=?")
		errorHandle(err, w)
		stmt.Exec(tem_goods.Goodsname, tem_goods.Address, tem_goods.Pic, tem_goods.Phone, tem_goods.Des, tm.Format("2006-01-02 15:04:05"), goods_id)
		dict.Msg = "修改成功"
		dict.Flag = 1
		// 上传图片成功
		dict_json, err := json.Marshal(dict)
		errorHandle(err, w)
		result = string(dict_json)
		fmt.Fprintf(w, result)
	}
}

func manage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	tem_goods := Goodsinfo{}
	dict := dic{}
	r.ParseForm()
	userid := r.URL.Query()["userid"][0]
	if r.Method == "GET" {
		rows, _ := db.Query("select * from lost_goods where status=? and user_id=?", 1, userid)
		for rows.Next() {
			rows.Scan(&tem_goods.Goodsid, &tem_goods.Goodsname, &tem_goods.Address, &tem_goods.Pic, &tem_goods.Phone, &tem_goods.Des, &tem_goods.Userid, &tem_goods.Create_time, &tem_goods.Update_time, &tem_goods.View_count, &tem_goods.Status, &tem_goods.Is_return)
			dict.Data = append(dict.Data, tem_goods)
		}
		dict_json, err := json.Marshal(dict)
		errorHandle(err, w)
		result := string(dict_json)
		fmt.Fprintf(w, result)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/login", login)
	r.HandleFunc("/register", register)
	// r.HandleFunc("^/media/file/", showPicHandle)
	r.PathPrefix("/media").HandlerFunc(showPicHandle)
	r.HandleFunc("/publish", publish)
	r.HandleFunc("/detail", detail)
	r.HandleFunc("/manage", manage)
	r.HandleFunc("/edit", edit)
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	var err error
	db, err = sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/lost_and_found?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Fatal(srv.ListenAndServe())
	http.Handle("/", r)
}
