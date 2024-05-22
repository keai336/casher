package grade

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

var MysqlOn bool = false
var Db *sqlx.DB

func init() {
	dsn := os.Getenv("DB_DSN")
	fmt.Println(dsn)

	if dsn == "" {
		fmt.Println("mysql option is not workable")
		return
	}

	database, err := sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	Db = database
	MysqlOn = true
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
func OneInsertChangeHistroy(groupname string, nowuse string, nowdelay int, newuse string, newdelay int) {
	currentTime := time.Now()
	t := currentTime.Format("2006-01-02 15:04:05")
	r, err := Db.Exec("insert into changehistory(groupname,nowuse,nowdelay,newuse,newdelay, datetime)values(?, ?, ?, ?,?,?)", groupname, nowuse, nowdelay, newuse, newdelay, t)

	if err != nil {
		fmt.Println("exec failed, ", err)
		return
	}
	_, err = r.LastInsertId()
	if err != nil {
		fmt.Println("exec failed, ", err)
		return
	}
}
