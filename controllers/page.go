package controllers

import (
	"../global"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)
type page struct {
	Domain			string
	Post_title 		string
	Post_content 	template.HTML
	Post_author		string
	Post_date		string
	Look_count		string
	Tags			[]string
	TagNum			int
	RandPage1		RandPage
	RandPage2		RandPage
	Comments		[]Comments
	CssFile			template.CSS
}
type RandPage struct {
	ID				string
	Post_title 		string
}
type Comments struct {
	Comment_author		string
	Comment_date 		string
	Comment_content		string
}

func GetHtml_page(Var global.Var,page_id string){

	if ok, _ := regexp.MatchString("^[0-9]+$", page_id); !ok {
		http.Redirect(Var.W,Var.R, "/", http.StatusFound)
		return
	}
	/*获取文章基本信息*/
	var CssFile,post_content,post_title,post_author,post_date,look_count string
	err := global.Db.QueryRow("select CssFile,post_content,post_title,post_author,post_date,look_count from caiji_posts where ID = ?", page_id).Scan(&CssFile,&post_content,&post_title,&post_author,&post_date,&look_count)
	if err != nil {
		fmt.Println(err)
		http.Redirect(Var.W,Var.R, "/", http.StatusFound)
		return
	}
	CssFile_:=strings.Split(CssFile,",")
	CssFile=""
	for _,v :=range CssFile_{
		b, err := ioutil.ReadFile("././"+v)
		if err != nil {
			fmt.Print(err)
		}
		CssFile =CssFile+ string(b)
	}

	/*随机取一篇文章*/
	var RandPage [2]RandPage
	global.Db.QueryRow("SELECT post_title,ID FROM caiji_posts AS t1 JOIN (SELECT ROUND(RAND() * ((SELECT MAX(ID) FROM `caiji_posts`)-(SELECT MIN(ID) FROM `caiji_posts`))+(SELECT MIN(ID) FROM `caiji_posts`)) AS ID_2) AS t2 WHERE t1.ID >= t2.ID_2 ORDER BY t1.ID LIMIT 1").Scan(&RandPage[0].Post_title,&RandPage[0].ID)
	global.Db.QueryRow("SELECT post_title,ID FROM caiji_posts AS t1 JOIN (SELECT ROUND(RAND() * ((SELECT MAX(ID) FROM `caiji_posts`)-(SELECT MIN(ID) FROM `caiji_posts`))+(SELECT MIN(ID) FROM `caiji_posts`)) AS ID_2) AS t2 WHERE t1.ID >= t2.ID_2 AND post_author=? ORDER BY t1.ID LIMIT 1",post_author).Scan(&RandPage[1].Post_title,&RandPage[1].ID)
	/*获取最新评论*/
	var comments []Comments
	var comments_ Comments
	rows,err := global.Db.Query("SELECT comment_author,comment_date,comment_content from caiji_comments where `comment_postID`=? limit 10",page_id)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&comments_.Comment_author,&comments_.Comment_date,&comments_.Comment_content)
		comments=append(comments,comments_)
	}
	/*获取标签*/
	var tags []string
	var name,term_taxonomy_id string
	rows,err = global.Db.Query("SELECT term_taxonomy_id from caiji_term_relationships where `object_id`=?",page_id)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&term_taxonomy_id)
		global.Db.QueryRow("SELECT `name` from caiji_terms where `term_id`=? limit 1", term_taxonomy_id).Scan(&name)
		tags=append(tags,name)
	}
	/*模板解析*/

	page := page{
		CssFile:template.CSS( CssFile),
		Domain:Var.R.Host,
		Post_title:post_title,
		Post_content:template.HTML(strings.Replace( post_content,"{{static}}",global.Config.Static_domain,-1)),
		Post_author:post_author,
		Post_date:post_date,
		Look_count:look_count,
		Tags:tags,
		TagNum:len(tags),
		RandPage1:RandPage[0],
		RandPage2:RandPage[1],
		Comments:comments,
	}
	t := template.Must(template.ParseGlob("././views/*"))
	err = t.ExecuteTemplate(Var.W, "page.html", page)
	global.CheckErr(err)
	/*模板解析*/

	/*数据更新*/
	stmt, err := global.Db.Prepare("UPDATE caiji_posts SET look_count=look_count+1 WHERE ID = ?")
	global.CheckErr(err)
	stmt.Exec(page_id)
	/*数据更新*/
}



