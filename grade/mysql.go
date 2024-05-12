package grade

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

var Db *sqlx.DB

func init() {
	database, err := sqlx.Open("mysql", "root:lk021104@tcp(127.0.0.1:3306)/delay_proxy")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	Db = database
}

func OneInsertHistory(proxy *GradeProxy) {
	currentTime := time.Now()
	t := currentTime.Format("2006-01-02 15:04:05")
	r, err := Db.Exec("insert into delay_test_history(proxy_name, delay, datetime,provider)values(?, ?, ?,?)", proxy.Name, proxy.DelayNow, t, proxy.Provider.Name)
	if err != nil {
		fmt.Println("exec failed, ", err)
		return
	}
	_, err = r.LastInsertId()
	if err != nil {
		fmt.Println("exec failed, ", err)
		return
	}

	//fmt.Println("insert succ:", id)
}
