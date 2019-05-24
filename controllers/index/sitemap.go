package index

import (
	"../../global"
	"fmt"
)

func GetHtml_sitemap(Var global.Var){

	/*获取最新200帖子*/
	rows,err := global.Db.Query("SELECT post_title,ID,post_date from caiji_posts ORDER BY ID DESC limit 200")
	global.CheckErr(err)
	defer rows.Close()
	var xml string
	for rows.Next() {
		var post_title,ID,post_date string

		rows.Scan(&post_title,&ID,&post_date)
		xml=xml+"<url><loc>https://"+Var.R.Host+"/page/"+ID+"</loc><lastmod>"+post_date+"</lastmod><changefreq>monthly</changefreq><priority>0.6</priority></url>"
	}
	xml="<?xml version=\"1.0\" encoding=\"UTF-8\"?><!-- baidu-sitemap-generator-version=\"1.6.5\" --><urlset xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd\" xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">"+xml+"</urlset>"

	Var.W.Header().Set("Content-Type", "application/xml;charset=UTF-8")
	fmt.Fprint(Var.W,xml)

}



