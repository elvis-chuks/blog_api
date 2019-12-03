package main

import (
	"database/sql"
	"fmt"
	"log"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/gorilla/handlers"
)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "password"
	dbname = "lilly"
)
type res map[string]interface{}

// struct used for login and register endpoints
type Person struct{
	Email string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Username  string `json:"username,omitempty"`
}
// struct used for likes
type Likes struct{
	Email string `json:"email,omitempty"`
	postID string `json:"id,omitempty"`
}
// struct used for comments
type Comments struct{
	Email string `json:"email,omitempty"`
	postID string `json:"id,omitempty"`
	Comment string `json:"comment,omitempty"`
}
func Register(w http.ResponseWriter, r *http.Request){
	setupResponse(&w, r)
	w.Header().Set("content-type","application/json")
	var person Person
	fmt.Println(r)
	_ = json.NewDecoder(r.Body).Decode(&person)// decodes the request body and parses it to the person variable
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
	" password=%s dbname=%s sslmode=disable",
	host,port,user,password,dbname)
	db, err := sql.Open("postgres",psqlInfo)
	if err != nil{
		panic(err)
	}
	defer db.Close()
	query := fmt.Sprintf("INSERT into users(email,password,username) VALUES('%s','%s','%s');",person.Email,person.Password,person.Username)
	_, err1 := db.Exec(query)
	if err1 != nil{
		fmt.Println(err1.Error())
		resp := res{"status":"error"}
		json.NewEncoder(w).Encode(resp)
	}else{
		result := res{"status":"success","msg":"proceed"}
		json.NewEncoder(w).Encode(result)
	}
}

func Login(w http.ResponseWriter,r *http.Request){
	setupResponse(&w, r)
	w.Header().Set("content-type","application/json")
	var person Person
	fmt.Println(r)
	_ = json.NewDecoder(r.Body).Decode(&person)
	fmt.Println(person)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
	" password=%s dbname=%s sslmode=disable",
	host,port,user,password,dbname)
	db, err := sql.Open("postgres",psqlInfo)
	if err != nil{
		panic(err)
	}
	defer db.Close()
	query := fmt.Sprintf("select email, password from users where email ='%s';",person.Email)
	rows,err := db.Query(query)
	if err != nil{
		fmt.Println(err.Error())
		resp := res{"status":"error"}
		json.NewEncoder(w).Encode(resp)
	}else{
		defer rows.Close()
		for rows.Next(){
			var email,password string
			err = rows.Scan(&email,&password)
			if person.Password == password{
				resp := res{"status":"success"}
				json.NewEncoder(w).Encode(resp)
			}else{
				resp := res{"status":"error","msg":"invalid password"}
				json.NewEncoder(w).Encode(resp)
			}
		}
	}
}

func Like(w http.ResponseWriter,r *http.Request){
	setupResponse(&w, r)
	if (*r).Method == "POST"{
		w.Header().Set("content-type","application/json")
	var cred Likes
	fmt.Println(r)
	_ = json.NewDecoder(r.Body).Decode(&cred)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
	" password=%s dbname=%s sslmode=disable",
	host,port,user,password,dbname)
	db, err := sql.Open("postgres",psqlInfo)
	if err != nil{
		panic(err)
	}
	defer db.Close()
	query := fmt.Sprintf("insert into likes(email,postid) values('%s','%s');",cred.Email,cred.postID)
	_, err1 := db.Exec(query)
	if err1 != nil{
		resp := res{"status":"error"}
		json.NewEncoder(w).Encode(resp)
	}else{
		result := res{"status":"success","msg":"proceed"}
		json.NewEncoder(w).Encode(result)
	}
	}
	
}

func Comment(w http.ResponseWriter, r *http.Request){
	setupResponse(&w, r)
	w.Header().Set("content-type","application/json")
	var comment Comments
	fmt.Println(r)
	_ = json.NewDecoder(r.Body).Decode(&comment)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
	" password=%s dbname=%s sslmode=disable",
	host,port,user,password,dbname)
	db, err := sql.Open("postgres",psqlInfo)
	if err != nil{
		panic(err)
	}
	defer db.Close()
	query := fmt.Sprintf("insert into comments(email,postid,comment) values('%s','%s','%s');",comment.Email,comment.postID,comment.Comment)
	_, err1 := db.Exec(query)
	if err1 != nil{
		resp := res{"status":"error"}
		json.NewEncoder(w).Encode(resp)
	}else{
		result := res{"status":"success","msg":"proceed"}
		json.NewEncoder(w).Encode(result)
	}
}
func setupResponse(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func main(){
	fmt.Println("Starting Application")
	fmt.Println("Application has started, running on http://127.0.0.1:5000")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
	" password=%s dbname=%s sslmode=disable",
	host,port,user,password,dbname)
	db, err := sql.Open("postgres",psqlInfo)
	if err != nil{
		panic(err)
	}
	defer db.Close()
	query := fmt.Sprintf("create table if not exists posts(id SERIAL,email varchar,header varchar,imgdata bytea,type varchar,body varchar,likes varchar,flag varchar);create table if not exists likes(id SERIAL, email varchar,postid varchar);create table if not exists comments(id SERIAL, email varchar,postid varchar,comment varchar);")
	_, err1 := db.Exec(query)
	if err1 != nil{
		panic(err1)
	}
	router := mux.NewRouter()
	router.HandleFunc("/v1/register",Register).Methods("POST")
	router.HandleFunc("/v1/login",Login).Methods("GET")
	router.HandleFunc("/v1/like",Like).Methods("POST")
	router.HandleFunc("/v1/comment",Comment).Methods("POST")
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(headers, methods, origins)(router)))
}