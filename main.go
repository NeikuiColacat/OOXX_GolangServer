package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Score    int    `json:"score"`
}
type UserScoreUpdate struct {
	Username string `json:"username"`
	Score    int    `json:"score"`
}
type UserQuery struct {
	Username string `json:"username"`
}
type UserScoreResponse struct {
	Score int `json:"score"`
}
type Response struct {
	Username string `json:"username"`
	Match    string `json:"match"`
}

type Player struct {
	Username string `json:"username"`
}

var (
	matchQueue   []Player
	matchResults = make(map[string]string)
	queueLock    sync.Mutex
)

func handleLogin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 登录逻辑
	var dbUser User
	err = db.QueryRow("SELECT username, password, score FROM users WHERE username = ?", user.Username).Scan(&dbUser.Username, &dbUser.Password, &dbUser.Score)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if dbUser.Password != user.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Login successful for user: %s, score: %d\n", user.Username, dbUser.Score)
}

func handleRegister(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 注册逻辑
	_, err = db.Exec("INSERT INTO users (username, password, score) VALUES (?, ?, ?)", user.Username, user.Password, 0)
	if err != nil {
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}

	fmt.Fprintf(w, "Registration successful for user: %s\n", user.Username)
}

func handleGetUsers(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT username, password, score FROM users ORDER BY score DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Username, &user.Password, &user.Score)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}
func handleUpdateScore(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Decode the incoming request
	var update UserScoreUpdate
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement
	stmt, err := db.Prepare("UPDATE users SET score = score + ? WHERE username = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(update.Score, update.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Score updated successfully"))
}
func handleQueryScore(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Decode the incoming request
	var query UserQuery
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT score FROM users WHERE username = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	row := stmt.QueryRow(query.Username)

	// Read the result
	var score int
	err = row.Scan(&score)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the response
	response := UserScoreResponse{
		Score: score,
	}

	// Send a success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleJoinMatchQueue(w http.ResponseWriter, r *http.Request) {
	var player Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queueLock.Lock()
	matchQueue = append(matchQueue, player)
	if len(matchQueue) >= 2 {
		// 弹出两个玩家进行配对
		player1 := matchQueue[0]
		player2 := matchQueue[1]
		matchQueue = matchQueue[2:]

		// 存储配对结果
		matchResults[player1.Username] = player2.Username
		matchResults[player2.Username] = player1.Username
		fmt.Println(player1.Username + " match with " + player2.Username)
	}
	queueLock.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%s joined the match queue successfully", player.Username)))
}
func handleQueryMatchResult(w http.ResponseWriter, r *http.Request) {
	var player Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queueLock.Lock()
	matchUsername, exists := matchResults[player.Username]
	queueLock.Unlock()

	if !exists {
		fmt.Println("no match found with" + player.Username)
		http.Error(w, "No match found", http.StatusNotFound)
		return
	}

	response := Response{
		Username: player.Username,
		Match:    matchUsername,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
func handleRemovePlayerFromMatch(w http.ResponseWriter, r *http.Request) {
	var player Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queueLock.Lock()
	defer queueLock.Unlock()

	matchUsername, exists := matchResults[player.Username]
	if !exists {
		http.Error(w, "Player not found in match queue", http.StatusNotFound)
		return
	}

	// 从匹配结果中删除这两个玩家
	fmt.Println(player.Username + " and " + matchUsername + "delete from match")
	delete(matchResults, player.Username)
	delete(matchResults, matchUsername)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Player %s and their match %s have been removed from the match queue", player.Username, matchUsername)))
}
func setupServer() (*sql.DB, *mux.Router) {
	// 连接数据库
	dsn := "root:maomao@tcp(127.0.0.1:3306)/game_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// 创建路由器
	r := mux.NewRouter()

	// 定义路由
	r.HandleFunc("/user/login", func(w http.ResponseWriter, r *http.Request) {
		handleLogin(db, w, r)
	}).Methods("POST")

	r.HandleFunc("/user/register", func(w http.ResponseWriter, r *http.Request) {
		handleRegister(db, w, r)
	}).Methods("POST")

	r.HandleFunc("/user/rank", func(w http.ResponseWriter, r *http.Request) {
		handleGetUsers(db, w, r)
	}).Methods("GET")

	r.HandleFunc("/user/score", func(w http.ResponseWriter, r *http.Request) {
		handleUpdateScore(db, w, r)
	}).Methods("POST")

	r.HandleFunc("/user/query", func(w http.ResponseWriter, r *http.Request) {
		handleQueryScore(db, w, r)
	}).Methods("POST")

	r.HandleFunc("/user/query_queue", handleJoinMatchQueue).Methods("POST")
	r.HandleFunc("/user/query_match", handleQueryMatchResult).Methods("POST")
	r.HandleFunc("/user/query_remove", handleRemovePlayerFromMatch).Methods("POST")
	return db, r
}

func main() {
	db, r := setupServer()
	defer db.Close()

	// 启动服务器
	log.Fatal(http.ListenAndServe(":8080", r))
}
