package user

import (
	"RelaAPI/src/commons"
	"encoding/json"
	"log"
	"net/http"
)

// GetHandler
//  @Description: 用户get方法
//  @param w
//  @param r
func GetHandler(w http.ResponseWriter,r *http.Request){
	//连接数据库
	db := commons.OpenConnection()

	//查询全部
	rows ,err :=db.Query("SELECT * FROM users")
	if err != nil{
		log.Fatal(err)
	}

	//返回容器
	var users []User

	for rows.Next() {
		var user User
		//存入容器
		rows.Scan(&user.Id,&user.Name,&user.Type)
		users = append(users,user)
	}
	//转为输出流
	peopleBytes,_ := json.MarshalIndent(users,"","\t")
	w.Header().Set("Content-Type","application/json")
	w.Write(peopleBytes)

	//关闭连接
	defer rows.Close()
	defer db.Close()
}

// POSTHandler
//  @Description: 用户Post方法
//  @param w
//  @param r
func POSTHandler(w http.ResponseWriter , r *http.Request){
	//连接数据库
	db := commons.OpenConnection()

	//从前台获取数据
	var u User
	err :=json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}

	//执行INSERT操作，人名不验重
	sqlStatement := `INSERT INTO users (id,name) VALUES (newid() ,$1 )`
	_, err = db.Exec(sqlStatement ,u.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	//关闭连接
	defer db.Close()
}
