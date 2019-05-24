package main

import (
	"./controllers"
	"./controllers/index"
	"./global"
	"./model"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)
const (
	certificateData = "./crtdate"
	defaultCrt      = "server.crt"
	defaultKey      = "server.key"
)




func route(w http.ResponseWriter, r *http.Request) {

	var Var global.Var
	Var.W=w
	Var.R=r

	if(strings.Contains(r.UserAgent(),"spider")){fmt.Println(time.Now(),r.UserAgent())}//统计蜘蛛UA


	dir := strings.Split(r.URL.Path,"/")
	//根据一级目录返回数据
	switch dir[1]{
	case "static":
		//返回静态文件 css js
		http.ServeFile(w, r, "."+r.URL.Path)
	case "page":
		//文章内容页面
		controllers.GetHtml_page(Var,dir[2])
	case "do":
		model.Do(Var,dir[2])
	case "tag":
		controllers.GetHtml_tag(Var,dir[2])
	case "author":
		controllers.GetHtml_Author(Var,dir[2])
	default:
		index_route(Var,dir[1])
	}

}
func index_route(Var global.Var,dir string){
	switch dir{
	case "":
		index.GetHtml_index(Var)
	case "sitemap":
		index.GetHtml_sitemap(Var)
	case "robots.txt":
		fmt.Fprint(Var.W,"Sitemap: https://"+Var.R.Host+"/sitemap")
	default:
		http.Redirect(Var.W, Var.R, "/", http.StatusFound)
	}
}

func main() {

	Init_config("./config.json")
	Init_server("./utf8.sql")

}

/*初始化配置*/
func Init_config(path string) {

	file1, err := os.OpenFile(path, os.O_RDWR, os.ModeType)
	if err != nil { log.Fatalf("%+v\n",errors.WithStack(err)) }
	defer file1.Close()
	marshal_data, err := ioutil.ReadAll(file1)
	//fmt.Println(string(marshal_data))
	if err != nil { log.Fatalf("%+v\n",errors.WithStack(err)) }
	err = json.Unmarshal(marshal_data, &global.Config)
	if err != nil { log.Fatalf("%+v\n",errors.WithStack(err)) }

	log.Println("载入配置文件"+path)
	log.Println(global.Config.Sql_str)



}
/*初始化服务端*/
func Init_server(path string){
	/*连接数据库*/
	log.Println("连接数据库"+global.Config.Sql_str)
	var err error
	global.Db, err = sql.Open("mysql", global.Config.Sql_str)
	err = global.Db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	global.Db.SetMaxOpenConns(global.Config.Max_connections)
	global.Db.SetMaxIdleConns(global.Config.Max_connections)
	log.Println("连接数据库成功")

	/*设置访问的路由*/
	http.HandleFunc("/", route)



}
