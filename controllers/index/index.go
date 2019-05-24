package index

import (
	"../../global"
	"fmt"
	"html/template"
	"log"
)
type index struct {
	Domain			string
	Tags 			[]string
	SiteName		string
}

func GetHtml_index(Var global.Var){

	var tags []string
	var name string

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

	err = rows.Err()
	global.CheckErr(err)


	/*模板解析*/
	index := index{
		SiteName:global.Config.SiteName,
		Domain:Var.R.Host,
		Tags:tags,
	}
	t := template.Must(template.ParseGlob("././views/*"))
	err = t.ExecuteTemplate(Var.W, "index.html", index)
	global.CheckErr(err)
	/*模板解析*/


}



