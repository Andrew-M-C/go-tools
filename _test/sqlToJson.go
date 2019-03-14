package _test

import (
	"time"
	"github.com/Andrew-M-C/go-tools/log"
	amcjson "github.com/Andrew-M-C/go-tools/json"
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

func TestSqlToJson() {
	item := SqlExample{}
	item.Id = 10
	item.NName = sql.NullString{String:"Andrew", Valid:true}
	item.NTime = mysql.NullTime{Time:time.Now(), Valid:false}
	item.NNum = sql.NullFloat64{Float64:1.23, Valid:false}
	item.NBool = sql.NullBool{Bool:true, Valid:true}
	item.NInt = sql.NullInt64{Int64:9999, Valid:true}
	item.Int = 8888
	item.String = "Hello,\tworld!"
	item.NoTagString = "empty string"
	item.Ignore = "ignore string"
	item.Bool = true
	item.Float = 1.110
	item.Time = time.Now()

	var ret string

	// normal mode
	ret, _ = amcjson.SqlToJson(&item)
	log.Debug("json result: %s", ret)

	// test float digit count after decomal point
	item.Float = 0.0
	ret, _ = amcjson.SqlToJson(item)
	log.Debug("json result: %s", ret)

	// test ShowNull option
	ret, _ = amcjson.SqlToJson(item, amcjson.Option{ShowNull: true})
	log.Debug("json result: %s", ret)

	// test option of float digit count after decomal point
	ret, _ = amcjson.SqlToJson(item, amcjson.Option{FloatDigits: 6, TimeDigits: 3})
	log.Debug("json result: %s", ret)

	// test include mode
	ret, _ = amcjson.SqlToJson(item, amcjson.Option{
		FilterMode: amcjson.IncludeMode,
		FilterList: []string{"nul_name", "nul_time", "nul_float"} })
	log.Debug("json result: %s", ret)

	// test exclude mode
	ret, _ = amcjson.SqlToJson(item, amcjson.Option{
		FilterMode: amcjson.ExcludeMode,
		FilterList: []string{"nul_name", "nul_time"} })
	log.Debug("json result: %s", ret)
	return
}
