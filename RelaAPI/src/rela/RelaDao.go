package rela

import (
	"RelaAPI/src/commons"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// GetHandler
//  @Description: 关系Get方法
//  @param w
//  @param r
func GetHandler(w http.ResponseWriter,r *http.Request){
	//split分割参数
	args := strings.Split(fmt.Sprint(r.URL),"/")
	//判断格式
	if args[3] != "relationships" {
		return
	}

	//连接数据库
	db := commons.OpenConnection()

	//查询数据库
	rows ,err :=db.Query("SELECT targetid , state ,type FROM rela WHERE mainid=$1",args[2])
	if err != nil{
		log.Fatal(err)
	}

	//返回容器
	var relas []Rela

	for rows.Next() {
		var rela Rela
		//存入容器
		rows.Scan(&rela.TargetId,&rela.State,&rela.Type)
		relas = append(relas,rela)
	}

	//转换为输出流
	relasBytes,_ := json.MarshalIndent(relas,"","\t")

	w.Header().Set("Content-Type","application/json")
	w.Write(relasBytes)

	//关闭连接
	defer rows.Close()
	defer db.Close()
}

// PUTHandler
//  @Description: 关系Put操作
//  @param w
//  @param r
func PUTHandler(w http.ResponseWriter , r *http.Request){
	//分割参数
	args := strings.Split(fmt.Sprint(r.URL),"/")
	//判断格式
	if args[3] != "relationships" || args[4] == ""{
		return
	}

	//连接数据库
	db := commons.OpenConnection()

	var re Rela
	//暂存两个人的id
	mainId , targetId := args[2] ,args[4]

	//从前台读取输入数据 state
	err :=json.NewDecoder(r.Body).Decode(&re)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}

	//仅喜欢的情况需要处理
	if re.State =="liked"{
		//查询对方是否喜欢
		rows ,err :=db.Query("SELECT state FROM rela WHERE mainid= $1 AND targetid = $2", targetId , mainId)
		if err != nil{
			log.Fatal(err)
		}
		//喜欢分两种情况 单向和配对
		if rows.Next(){
			var rela Rela
			rows.Scan(&rela.State)
			//如果对方也喜欢,就配对
			if rela.State == "liked" {
				re.State ="matched"
				//同时更新对方
				rela.State ="matched"
				//执行UPDATE语句
				sql := `UPDATE rela SET state =$1 WHERE mainid = $2 AND targetid = $3`
				_, err = db.Exec(sql ,rela.State,targetId,mainId)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					panic(err)
				}
			}
		}
	}

	//拼接个人更新语句 验重
	sqlStatement := `INSERT INTO rela (mainid,targetid,state) VALUES ($1 ,$2 ,$3 )
					ON conflict(mainid,targetid)
					DO UPDATE SET state =$4`
	_, err = db.Exec(sqlStatement ,mainId,targetId,re.State,re.State)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	//关闭连接
	defer db.Close()
}