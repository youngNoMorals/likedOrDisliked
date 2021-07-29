package main

import (
	"RelaAPI/src/rela"
	"RelaAPI/src/user"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main(){
	//对/users 和 /users/分别处理
	http.HandleFunc("/users",UserHandler)
	http.HandleFunc("/users/",RelaHandler)

	log.Fatal(http.ListenAndServe(":80",nil))
}

// UserHandler
//  @Description: 用户处理
//  @param w
//  @param r
func UserHandler(w http.ResponseWriter, r *http.Request){
	if r.Method =="GET"{
		user.GetHandler(w,r)
	}else if r.Method == "POST"{
		user.POSTHandler(w,r)
	}else{
		return
	}
}

// RelaHandler
//  @Description: 关系处理
//  @param w
//  @param r
func RelaHandler(w http.ResponseWriter, r *http.Request){
	if r.Method =="GET"{
		rela.GetHandler(w,r)
	}else if r.Method == "PUT"{
		rela.PUTHandler(w,r)
	}else{
		return
	}
}

