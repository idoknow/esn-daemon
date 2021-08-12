package db

import (
	"database/sql"
	"esnd/src/util"

	_ "github.com/go-sql-driver/mysql"
)

var Cfg *util.Config
var DB *sql.DB

type InfoTables struct {
	name string
}

func Init(cfg *util.Config) error {
	Cfg = cfg
	db, err := sql.Open("mysql", cfg.GetAnyway("db.user", "esnd")+":"+cfg.GetAnyway("db.pass", "changeMe")+
		"@tcp("+cfg.GetAnyway("db.addr", "127.0.0.1:3306")+")/"+cfg.GetAnyway("db.database", "esnd"))
	if err != nil {
		util.SaySub("Database", "err:Cannot init database:"+err.Error())
		return err
	}
	DB = db

	//debug
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return err
	}
	defer rows.Close()

	user := false
	noti := false

	for rows.Next() {
		var table InfoTables
		err = rows.Scan(&table.name)
		if err != nil {
			util.DebugMsg("DB", "err:"+err.Error())
			continue
		}
		if table.name == "users" {
			user = true
		} else if table.name == "notis" {
			noti = true
		}
		util.DebugMsg("DB", "table:"+table.name)
	}
	//check
	if !user {
		err = CreateUserTable()
		if err != nil {
			return err
		}
	}
	if !noti {
		err = CreateNotiTable()
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateUserTable() error {
	util.SaySub("DB", "Creating users table")
	_, err := DB.Exec("CREATE TABLE users(id bigint not null primary key auto_increment,name varchar(255) not null,mask varchar(255) not null,priv varchar(255) not null);")
	return err
}
func CreateNotiTable() error {
	util.SaySub("DB", "Creating notis table")
	_, err := DB.Exec("CREATE TABLE notis (id bigint not null primary key auto_increment,target varchar(255) not null,time varchar(255) not null,title varchar(255) not null,content varchar(1023) not null,source varchar(255) not null,token varchar(255) not null)")
	return err
}

type count struct {
	count int
}

func Count(query string) int {
	util.DebugMsg("count", "Count:"+query)

	var c count
	row := DB.QueryRow(query)

	err := row.Scan(&c.count)

	if err != nil {
		return -1
	}
	return c.count

}
