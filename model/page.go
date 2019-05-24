package model

import (
	"../global"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"time"
)


type NewPage struct {
	Post_content	string
	Post_author		string
	Post_title		string
	Post_source		string
	Comment_status	string
	Post_password	string
	Tags		    string
	CssFile			string//css文件路径 ,分割
}
type NewComment struct {
	Comment_author			string
	Comment_content			string
	Comment_postID			int
	Comment_parent			int
}
func Do(Var global.Var,function string){
	/*初始化一下要反射的函数列表*/
	funcs := map[string]interface{}{
		"Page_add":Page_add,
		"Comment_add":Comment_add,
		"Page_isExist":Page_isExist,//根据原文url检测文章是否存在
		"Page_Number":Page_Number,//返回文章数量
	}

	Call(funcs, function,Var)
}

func Call(m map[string]interface{}, name string, params ... interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		fmt.Println("The number of params is not adapted.")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}

func Page_add(Var global.Var){
	/*接受数据*/
	body_, _ := ioutil.ReadAll(Var.R.Body) //获取post的数据
	var body NewPage
	err := json.Unmarshal(body_, &body)
	tags:=strings.Split(body.Tags,",")
	//fmt.Println(tags)
	global.CheckErr(err)
	/*接受数据解析到json*/

	/*文章插入到数据库*/
	stmt, err := global.Db.Prepare("insert into caiji_posts(CssFile,post_content,post_author,post_title,post_source,comment_status,post_password,post_date,post_status)values(?,?,?,?,?,?,?,?,?)")
	global.CheckErr(err)
	//fmt.Println(body.CssFile,body.Post_author,body.Post_title,body.Post_source,body.Comment_status,body.Post_password)
	h, _ := time.ParseDuration("-1h")
	rs,err :=stmt.Exec(body.CssFile,body.Post_content,body.Post_author,body.Post_title,body.Post_source,body.Comment_status,body.Post_password,time.Now().Add(9 * h),"已发表")
	if(err != nil){
		return
	}
	LastPage,_:=rs.LastInsertId()
	/*标签插入到数据库*/


	for i, v := range tags {
		var term_id int64
		err=global.Db.QueryRow("SELECT term_id from caiji_terms where name=? limit 1", tags[i]).Scan(&term_id)
		fmt.Println(term_id,tags[i])
		if(term_id==0){
			stmt, err = global.Db.Prepare("insert into caiji_terms SET name=?")
			rs,_ :=stmt.Exec(v)
			fmt.Println(v)
			term_id,_=rs.LastInsertId()

		}
		stmt, err = global.Db.Prepare("insert into caiji_term_relationships SET `object_id`=?,`term_taxonomy_id`=?")
		stmt.Exec(LastPage,term_id)
	}
	global.LastPage=LastPage
	fmt.Println("发表文章: "+body.Post_title)

}
func Comment_add(Var global.Var){
	/*接受数据*/
	body_, _ := ioutil.ReadAll(Var.R.Body) //获取post的数据
	var body NewComment
	err := json.Unmarshal(body_, &body)
	global.CheckErr(err)
	/*接受数据解析到json*/

	/*数据更新*/
	stmt, err := global.Db.Prepare("insert into caiji_comments(comment_postID,comment_author,comment_author_IP,comment_status,comment_content,comment_agent,comment_date,comment_parent)values(?,?,inet_aton(?),?,?,?,?,?)")
	global.CheckErr(err)
	stmt.Exec(body.Comment_postID,body.Comment_author,strings.SplitN(Var.R.RemoteAddr,":",2)[0] ,"已发表",body.Comment_content,Var.R.UserAgent(),time.Now(),body.Comment_parent)
	/*数据更新*/

	fmt.Println("更新评论: "+body.Comment_content)

}

func Page_isExist(Var global.Var)  {
	/*接受数据*/
	PageSource:=Var.R.PostFormValue("PageSource") //获取Get的数据

	/*查询是否存在*/
	var post_title string
	global.Db.QueryRow("SELECT post_title from caiji_posts where post_source=? limit 1", PageSource).Scan(&post_title)

	if(post_title==""){
		fmt.Fprintf(Var.W,"Flase")
	}else {
		fmt.Fprintf(Var.W,"True")
	}
}

func Page_Number(Var global.Var)  {

	fmt.Fprint(Var.W,global.LastPage)
}