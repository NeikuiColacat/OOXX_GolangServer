package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// User 结构体表示用户对象
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// JSON文件路径
const usersFilePath = "users.json"

func dbConn() *sql.DB {
	// MySQL数据库连接信息
	dbDriver := "mysql"
	dbUser := "root"   // 修改为您的MySQL用户名
	dbPass := "maomao" // 修改为您的MySQL密码
	dbName := "game_db"

	// 连接数据库
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// 读取本地JSON文件中的用户信息
		users, err := readUsersFromJSON(usersFilePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 遍历JSON中的用户信息，进行验证
		valid := false
		for _, user := range users {
			if user.Username == username && user.Password == password {
				valid = true
				break
			}
		}

		// 验证成功则设置Cookie并重定向到/dashboard，否则返回未授权错误
		if valid {
			cookie := http.Cookie{
				Name:  "username",
				Value: username,
				Path:  "/",
			}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else {
			// 返回登录页面并显示错误消息
			http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
		}
		return
	}

	// 读取登录表单页面
	html, err := ioutil.ReadFile("login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 检查是否有错误消息
	var pageContent string
	if r.URL.Query().Get("error") == "1" {
		errorMsg := "<p style='color: red;'>Invalid username or password. Please try again.</p>"
		pageContent = strings.Replace(string(html), "{{ERROR_MSG}}", errorMsg, 1)
	} else {
		pageContent = strings.Replace(string(html), "{{ERROR_MSG}}", "", 1)
	}

	// 将页面内容写入响应
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, pageContent)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// 检查用户是否已登录
	username, err := getLoggedInUsername(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if username == "" {
		// 未登录，重定向到登录页面
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	db := dbConn()
	defer db.Close()

	// 查询时按score列降序排列
	rows, err := db.Query("SELECT username, password, score FROM users ORDER BY score DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tableRows string
	for rows.Next() {
		var username, password string
		var score int
		err = rows.Scan(&username, &password, &score)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tableRows += fmt.Sprintf("<tr><td>%s</td><td class='editable' data-field='password'>%s</td><td class='editable' data-field='score'>%d</td></tr>", username, password, score)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 读取模板文件
	template, err := ioutil.ReadFile("template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 替换占位符
	pageContent := strings.Replace(string(template), "{{TABLE_ROWS}}", tableRows, -1)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, pageContent)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var updateData struct {
		Username string `json:"username"`
		Field    string `json:"field"`
		Value    string `json:"value"`
	}

	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := dbConn()
	defer db.Close()

	query := fmt.Sprintf("UPDATE users SET %s = ? WHERE username = ?", updateData.Field)
	stmt, err := db.Prepare(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(updateData.Value, updateData.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func readUsersFromJSON(filename string) ([]User, error) {
	// 读取并解析JSON文件
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var users []User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func getLoggedInUsername(r *http.Request) (string, error) {
	cookie, err := r.Cookie("username")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", nil // 没有找到Cookie，视为未登录
		}
		return "", err
	}
	return cookie.Value, nil
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/dashboard", handler)
	http.HandleFunc("/update", updateHandler)
	log.Println("Server started on: http://localhost:8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal(err)
	}
}
