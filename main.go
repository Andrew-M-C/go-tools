package main

import (
	"time"
	"./log"
	amcjson "./json"
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

type SqlExample struct {
	Id		int32				`db:"id"`
	NName	sql.NullString		`db:"nul_name"`
	NTime	mysql.NullTime		`db:"nul_time"`
	NNum	sql.NullFloat64		`db:"nul_float"`
	NBool	sql.NullBool		`db:"nul_bool"`
	NInt	sql.NullInt64		`db:"nul_int"`
	Int		int					`db:"int" json:"int_num"`
	String	string				`db:"string" json:"str,omitempty"`
	EmptyString	string			`json:"empty_str,omitempty"`
	NoTagString	string
	Ignore	string				`json:"-"`
	Bool	bool				`json:"bool"`
	Float	float32				`json:"float32"`
	Time	time.Time			`json:"time"`
}

func main() {
	item := SqlExample{}
	item.Id = 10
	item.NName = sql.NullString{String:"Andrew", Valid:true}
	item.NTime = mysql.NullTime{Time:time.Now(), Valid:false}
	item.NNum = sql.NullFloat64{Float64:1.23, Valid:false}
	item.NBool = sql.NullBool{Bool:true, Valid:true}
	item.NInt = sql.NullInt64{Int64:9999, Valid:true}
	item.Int = 8888
	item.String = "Hello, world!"
	item.NoTagString = "empty string"
	item.Ignore = "ignore string"
	item.Bool = true
	item.Float = 0.001
	item.Time = time.Now()

	var ret string
	// ret, _ = amcjson.SqlToJson(&item)
	// log.Debug("json result: %s", ret)
	ret, _ = amcjson.SqlToJson(item)
	log.Debug("json result: %s", ret)
	return
}
