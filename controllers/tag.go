package controllers

import (
	"../global"
	"fmt"
	"html/template"
	"log"
)
type Tag struct {
	Domain			string
	Tags 			[]string
	Tag				string
	SiteName		string
	TagNum			int
	Object_id		[]string
	Post_title		[]string
	RowMaps			[]map[string]interface{}
}

func GetHtml_tag(Var global.Var,tag string){
	var object_id,post_title string
	var tags []string
	var term_taxonomy_id int
	var name string

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

	/*根据标签名字获取最新帖子*/
	global.Db.QueryRow("select term_id from caiji_terms where name = ?", tag).Scan(&term_taxonomy_id)
	rows,err = global.Db.Query("SELECT object_id from caiji_term_relationships where term_taxonomy_id=? ORDER BY object_id DESC limit 10",term_taxonomy_id)
	defer rows.Close()
	var rowMaps []map[string]interface{}
	if(err==nil) {
		for rows.Next() {
			rows.Scan(&object_id)
			global.Db.QueryRow("SELECT post_title from caiji_posts where ID=? limit 1", object_id).Scan(&post_title)

			each := map[string]interface{}{}
			each["post_title"] = post_title
			each["object_id"] = object_id
			//fmt.Println(object_id,post_title)
			rowMaps =append(rowMaps,each)
		}
	}
	//rows,err = global.Db.Query("SELECT post_title from caiji_posts where ID IN (?) limit 10",strings.Replace(strings.Trim(fmt.Sprint(object_id), "[]"), " ", ",", -1))

	//fmt.Println(rowMaps)
	/*模板解析*/
	tag_ := Tag{
		Domain:Var.R.Host,
		SiteName:global.Config.SiteName,
		Tags:tags,
		Tag:tag,
		TagNum:len(rowMaps),
		RowMaps:rowMaps,
	}
	//fmt.Println(tag_.TagNum)
	t := template.Must(template.ParseGlob("././views/*"))
	err = t.ExecuteTemplate(Var.W, "tag.html", tag_)
	global.CheckErr(err)
	/*模板解析*/


}



