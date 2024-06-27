package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"   // 修改为您的MySQL用户名
	dbPass := "maomao" // 修改为您的MySQL密码
	dbName := "game_db"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func handler(w http.ResponseWriter, r *http.Request) {
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
		tableRows += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%d</td></tr>", username, password, score)
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

func main() {
	http.HandleFunc("/", handler)
	log.Println("Server started on: http://localhost:8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal(err)
	}
}
