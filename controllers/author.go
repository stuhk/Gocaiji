package controllers

import (
	"../global"
	"fmt"
	"html/template"
	"log"
)
type author struct {
	Domain			string
	Tags 			[]string
	Author			string
	SiteName		string
	AuthorNum			int
	Post_title		[]string
	Author_page		[]author_page
}
type author_page struct {
	ID				string
	Post_title		string
}
func GetHtml_Author(Var global.Var,Author string){
	var name string
	var tags []string

	/*随机获取十个标签*/
	rows,err := global.Db.Query("SELECT name FROM `caiji_terms` AS t1 JOIN (SELECT ROUND(RAND() * ((SELECT MAX(term_id) FROM `caiji_terms`)-(SELECT MIN(term_id) FROM `caiji_terms`))+(SELECT MIN(term_id) FROM `caiji_terms`)) AS term_id) AS t2 WHERE t1.term_id >= t2.term_id ORDER BY t1.term_id LIMIT 10")
	if err != nil {
		log.Println(err)
		fmt.Fprintf(Var.W, "出错,请刷新")
		return
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&name)
		tags=append(tags,name)
	}

	/*根据作者名字获取最新10帖子*/
	rows,err = global.Db.Query("select post_title,ID from caiji_posts where post_author = ? limit 18", Author)
	defer rows.Close()
	var Author_page []author_page
	var authorPage author_page
	for rows.Next() {
		rows.Scan(&authorPage.Post_title,&authorPage.ID)
		fmt.Println(authorPage)
		Author_page =append(Author_page,authorPage)
	}

	/*模板解析*/
	Author_ := author{
		Domain:Var.R.Host,
		SiteName:global.Config.SiteName,
		Tags:tags,
		Author:Author,
		AuthorNum:len(Author_page),
		Author_page:Author_page,
	}

	t := template.Must(template.ParseGlob("././views/*"))
	err = t.ExecuteTemplate(Var.W, "author.html", Author_)
	global.CheckErr(err)
	/*模板解析*/


}



