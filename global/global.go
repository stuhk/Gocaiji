package global

import (
	"database/sql"
	"log"
	"net/http"
)

type ConfiG struct {
	Sql_str			 string		`json:"sql_str"`
	Max_connections	 int		`json:"Max_connections"`
	Static_domain	 string		`json:"static_domain"`
	SiteName	     string
	Port		     string
}
type Var struct {
	W		   http.ResponseWriter
	R 		   *http.Request
}
var (
	Config     *ConfiG
	Db         *sql.DB

	LastPage	int64
)
func CheckErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
