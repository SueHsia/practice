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

	"github.com/garyburd/redigo/redis"

	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/gorilla/mux"
)

const (
	upload_path = "."
)

var (
	db        *sql.DB
	pubUserid int
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
	resPage := r.URL.Query()["page"][0]
	if resPage == "" {
		resPage = "1"
	}
	dict := dic{}
	temGoods := Goodsinfo{}
	page, _ := strconv.Atoi(resPage)
	rows, err := db.Query("select * from lost_goods where status = ? limit ?,8", 1, (int(page)-1)*8)
	errorHandle(err, w)
	db.QueryRow("select count(*) from lost_goods where status = ?", 1).Scan(&dict.Total_count)
	dict.Total_page = int(math.Ceil(float64(dict.Total_count) / 8.0))
	errorHandle(err, w)
	for rows.Next() {
		rows.Scan(&temGoods.Goodsid, &temGoods.Goodsname, &temGoods.Address, &temGoods.Pic, &temGoods.Phone, &temGoods.Des, &temGoods.Userid, &temGoods.Create_time, &temGoods.Update_time, &temGoods.View_count, &temGoods.Status, &temGoods.Is_return)
		dict.Data = append(dict.Data, temGoods)
	}
	dictJSON, _ := json.Marshal(dict)
	result := string(dictJSON)
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
	temData := Userinfo{}
	var result string
	if r.Method == "POST" {
		r.ParseForm()
		username := r.Form["username"][0]
		password := r.Form["password"][0]
		rows := db.QueryRow("select * from pi_user where username = ?", username)
		rows.Scan(&temData.Userid, &temData.Username, &temData.Password, &temData.Create_time)
		if temData.Username == "" {
			dict.Msg = "用户名或密码错误"
			dict.Flag = 0
		} else {
			newPass := md5V(password)
			if newPass != temData.Password {
				dict.Msg = "用户名或密码错误"
				dict.Flag = 0
			} else {
				dict.Flag = 1
			}
			dict.Username = temData.Username
			dict.Userid = temData.Userid
			dict.Token = enbase(temData.Username + "," + strconv.Itoa(temData.Userid) + "," + strconv.FormatInt(time.Now().Unix(), 10))
		}
		dictJSON, _ := json.Marshal(dict)
		result = string(dictJSON)
	}
	fmt.Fprintf(w, string(result))
}
func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	dict := dic{}
	temData := Userinfo{}
	var result string
	if r.Method == "POST" {
		r.ParseForm()
		username := strings.Replace(r.Form["username"][0], " ", "", -1)
		password := strings.Replace(r.Form["password"][0], " ", "", -1)
		repassword := strings.Replace(r.Form["repassword"][0], " ", "", -1)
		rows := db.QueryRow("select * from pi_user where username = ?", username)
		rows.Scan(&temData.Userid, &temData.Username, &temData.Password, &temData.Create_time)
		if password != repassword {
			dict.Msg = "密码不一致，注册失败！"
			dict.Flag = 0
		} else if temData.Username != "" {
			dict.Msg = "用户已存在，注册失败"
			dict.Flag = 0
		} else {
			password = md5V(password)
			tm := time.Unix(time.Now().Unix(), 0)
			createTime := tm.Format("2006-01-02 15:04:05")
			stmt, err := db.Prepare("INSERT pi_user SET username=?,password=?,create_time=?")
			errorHandle(err, w)
			stmt.Exec(username, password, createTime)
			dict.Msg = "注册成功！"
			dict.Flag = 1
		}
		dictJSON, err := json.Marshal(dict)
		errorHandle(err, w)
		result = string(dictJSON)
	}
	fmt.Fprintf(w, result)
}

func publish(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "multipart/form-data")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	dict := dic{}
	temUser := Userinfo{}
	temGoods := Goodsinfo{}
	var (
		ext         string
		isValid     bool
		result      string
		curUsername string
	)
	r.ParseForm()
	if r.Method == "OPTIONS" {
		fmt.Fprintf(w, "")
		return
	} else if r.Method == "GET" {
		token := r.Form["token"][0]
		curUsername, pubUserid, isValid = validate(token)
		if !isValid {
			dict.Msg = "登录失效，请重新登录"
			dict.Flag = 0
			dictJSON, err := json.Marshal(dict)
			errorHandle(err, w)
			result = string(dictJSON)
			fmt.Fprintf(w, result)
			return
		}
	} else if r.Method == "POST" {
		uploadFile, handle, err := r.FormFile("file")
		errorHandle(err, w)
		tm := time.Unix(time.Now().Unix(), 0)
		temGoods.Address = r.Form["address"][0]
		temUser.Username = curUsername
		temGoods.Goodsname = r.Form["name"][0]
		temGoods.Userid = pubUserid
		temGoods.Phone = r.Form["phone"][0]
		temGoods.Des = r.Form["des"][0]
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
		temGoods.Pic = filename[1:]
		stmt, err := db.Prepare("INSERT lost_goods SET name=?,address=?,pic=?,phone=?,des=?,user_id=?,create_time=?,update_time=?,view_count=?,status=?,is_return=?")
		errorHandle(err, w)
		stmt.Exec(temGoods.Goodsname, temGoods.Address, temGoods.Pic, temGoods.Phone, temGoods.Des, temGoods.Userid, tm.Format("2006-01-02 15:04:05"), tm.Format("2006-01-02 15:04:05"), 0, 1, 0)
		dict.Msg = "发布成功"
		dict.Flag = 1
		// 上传图片成功
		dictJSON, err := json.Marshal(dict)
		errorHandle(err, w)
		result = string(dictJSON)
		fmt.Fprintf(w, result)
	}
}

func detail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	con, err := redis.Dial("tcp", "127.0.0.1:6379")
	errorPrint(err)
	temGoods := Goodsinfo{}
	dict := dic{}
	var goodsID string
	if r.Method == "OPTIONS" {
		fmt.Fprintf(w, "")
		return
	} else if r.Method == "GET" {
		goodsID = r.Form["goods_id"][0]
		rows := db.QueryRow("select * from lost_goods where id = ?", goodsID)
		rows.Scan(&temGoods.Goodsid, &temGoods.Goodsname, &temGoods.Address, &temGoods.Pic, &temGoods.Phone, &temGoods.Des, &temGoods.Userid, &temGoods.Create_time, &temGoods.Update_time, &temGoods.View_count, &temGoods.Status, &temGoods.Is_return)
		stmt, err := db.Prepare("update lost_goods set view_count=? where id=?")
		errorHandle(err, w)
		redisCount := strings.Split(mysql_redis(temGoods.Goodsid, temGoods.View_count), "-")
		curName := fmt.Sprintf("viewCount-%v", goodsID)
		oriTime, err := strconv.ParseInt(redisCount[1], 10, 64)
		errorPrint(err)
		count, _ := strconv.Atoi(redisCount[0])
		temGoods.View_count = count
		curTime := time.Now().Unix()
		if curTime-oriTime >= 20 {
			stmt.Exec(count, temGoods.Goodsid)
			redisContent := fmt.Sprintf("%v-%v", count, curTime)
			_, err = con.Do("SET", curName, redisContent)
		}
		dict.Data = append(dict.Data, temGoods)
		dictJSON, err := json.Marshal(dict)
		errorHandle(err, w)
		result := string(dictJSON)
		fmt.Fprintf(w, result)
	}
}

func edit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	dict := dic{}
	temGoods := Goodsinfo{}
	r.ParseForm()
	var (
		token     string
		curUserid int
		isValid   bool
		result    string
	)
	if r.Method == "GET" {
		token = r.Form["token"][0]
		_, curUserid, isValid = validate(token)
		if !isValid {
			dict.Msg = "登录失效，请重新登录"
			dict.Flag = 0
		} else {
			goodsID, _ := strconv.Atoi(r.Form["goods_id"][0])
			rows := db.QueryRow("select * from lost_goods where id = ?", goodsID)
			rows.Scan(&temGoods.Goodsid, &temGoods.Goodsname, &temGoods.Address, &temGoods.Pic, &temGoods.Phone, &temGoods.Des, &temGoods.Userid, &temGoods.Create_time, &temGoods.Update_time, &temGoods.View_count, &temGoods.Status, &temGoods.Is_return)
			dict.Data = append(dict.Data, temGoods)
			if temGoods.Is_return == 1 {
				dict.Msg = "失物已归还禁止编辑"
				dict.Flag = 0
			} else {
				dict.Flag = 1
			}
		}
		dictJSON, err := json.Marshal(dict)
		errorHandle(err, w)
		result = string(dictJSON)
		fmt.Fprintf(w, result)
		return
	} else if r.Method == "POST" && temGoods.Userid == curUserid {
		var ext string
		uploadFile, handle, err := r.FormFile("file")
		errorHandle(err, w)
		tm := time.Unix(time.Now().Unix(), 0)
		goodsID, _ := strconv.Atoi(r.Form["goods_id"][0])
		temGoods.Address = r.Form["address"][0]
		temGoods.Goodsname = r.Form["name"][0]
		temGoods.Phone = r.Form["phone"][0]
		temGoods.Des = r.Form["des"][0]
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
		temGoods.Pic = filename[1:]
		stmt, err := db.Prepare("UPDATE lost_goods SET name=?,address=?,pic=?,phone=?,des=?,update_time=? WHERE id=?")
		errorHandle(err, w)
		stmt.Exec(temGoods.Goodsname, temGoods.Address, temGoods.Pic, temGoods.Phone, temGoods.Des, tm.Format("2006-01-02 15:04:05"), goodsID)
		dict.Msg = "修改成功"
		dict.Flag = 1
		// 上传图片成功
		dictJSON, err := json.Marshal(dict)
		errorHandle(err, w)
		result = string(dictJSON)
		fmt.Fprintf(w, result)
		return
	} else {
		dict.Msg = "用户无权编辑"
		dict.Flag = 0
		dictJSON, err := json.Marshal(dict)
		errorHandle(err, w)
		result = string(dictJSON)
		fmt.Fprintf(w, result)
		return
	}
}

func manage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	temGoods := Goodsinfo{}
	dict := dic{}
	var result string
	r.ParseForm()
	if r.Method == "OPTIONS" {
		fmt.Fprintf(w, "")
		return
	} else {
		token := r.Form["token"][0]
		_, curUserid, isValid := validate(token)
		if !isValid {
			dict.Msg = "登录失效，请重新登录"
			dict.Flag = 0
			dictJSON, err := json.Marshal(dict)
			errorHandle(err, w)
			result = string(dictJSON)
			fmt.Fprintf(w, result)
			return
		} else if r.Method == "GET" {
			rows, _ := db.Query("select * from lost_goods where status=? and user_id=?", 1, curUserid)
			for rows.Next() {
				rows.Scan(&temGoods.Goodsid, &temGoods.Goodsname, &temGoods.Address, &temGoods.Pic, &temGoods.Phone, &temGoods.Des, &temGoods.Userid, &temGoods.Create_time, &temGoods.Update_time, &temGoods.View_count, &temGoods.Status, &temGoods.Is_return)
				dict.Data = append(dict.Data, temGoods)
			}
			dictJSON, err := json.Marshal(dict)
			errorHandle(err, w)
			result = string(dictJSON)
			fmt.Fprintf(w, result)
			return
		} else if r.Method == "PUT" {
			goodsID := r.Form["goods_id"][0]
			doReturn := r.Form["return_code"][0]
			db.QueryRow("select user_id,is_return from lost_goods where id=?", goodsID).Scan(&temGoods.Userid, &temGoods.Is_return)
			if curUserid == temGoods.Userid && doReturn == "1" {
				stmt, err := db.Prepare("update lost_goods set is_return=? where id=?")
				errorHandle(err, w)
				stmt.Exec(1, goodsID)
				dict.Flag = 1
				dict.Msg = "归还成功"
			} else if curUserid == temGoods.Userid && temGoods.Is_return == 0 {
				dict.Msg = "失物可编辑"
				dict.Is_edit = 1
			} else {
				dict.Msg = "用户无权编辑"
				dict.Is_edit = 0
			}
			dictJSON, err := json.Marshal(dict)
			errorHandle(err, w)
			result := string(dictJSON)
			fmt.Fprintf(w, result)
			return
		} else if r.Method == "DELETE" {
			goodsID := r.Form["goods_id"][0]
			deleteCode := r.Form["delete_code"][0]
			db.QueryRow("select user_id,status from lost_goods where id=?", goodsID).Scan(&temGoods.Userid, &temGoods.Status)
			if curUserid == temGoods.Userid && deleteCode == "1" {
				stmt, err := db.Prepare("update lost_goods set status=? where id=?")
				errorHandle(err, w)
				stmt.Exec(0, goodsID)
				dict.Flag = 1
				dict.Msg = "删除成功"
			} else {
				dict.Msg = "删除失败"
			}
			dictJSON, err := json.Marshal(dict)
			errorHandle(err, w)
			result := string(dictJSON)
			fmt.Fprintf(w, result)
			return
		}
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
