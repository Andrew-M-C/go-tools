package _test

import (
	"time"
	"github.com/Andrew-M-C/go-tools/log"
	"github.com/Andrew-M-C/go-tools/jsonconv"
	"github.com/Andrew-M-C/go-tools/sqlconv"
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


func TestReadSqlKVs() {
	log.Info("====== now test GetValidKVsFromStruct")
	item := SqlExample{}
	item.Id = 10
	item.NName = sql.NullString{String:"中文", Valid:true}
	item.NTime = mysql.NullTime{Time:time.Now(), Valid:false}
	item.NNum = sql.NullFloat64{Float64:1.23, Valid:false}
	item.NBool = sql.NullBool{Bool:true, Valid:true}
	item.NInt = sql.NullInt64{Int64:9999, Valid:true}
	item.Int = 8888
	item.String = "Hello,\tworld!"
	item.NoTagString = "\"\"'''```"
	item.Ignore = "ignore string"
	item.Bool = true
	item.Float = 1.110
	item.Time = time.Now()

	k, v, err := sqlconv.GetValidKVsFromStruct(&item, "'")
	if err != nil {
		log.Error("Getting failed: %s", err.Error())
		return
	}

	for i := 0; i < len(k); i++ {
		log.Info("%s -- %s", k[i], v[i])
	}

	return
}


func TestSqlToJson() {
	log.Info("====== now test SqlToJson")
	item := SqlExample{}
	item.Id = 10
	item.NName = sql.NullString{String:"中文", Valid:true}
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
	ret, _ = jsonconv.SqlToJson(&item)
	log.Debug("json result: %s", ret)

	// test float digit count after decomal point
	item.Float = 0.0
	ret, _ = jsonconv.SqlToJson(item)
	log.Debug("json result: %s", ret)

	// test ShowNull option
	ret, _ = jsonconv.SqlToJson(item, jsonconv.Option{ShowNull: true})
	log.Debug("json result: %s", ret)

	// test option of float digit count after decomal point
	ret, _ = jsonconv.SqlToJson(item, jsonconv.Option{FloatDigits: 6, TimeDigits: 3})
	log.Debug("json result: %s", ret)

	// test include mode
	ret, _ = jsonconv.SqlToJson(item, jsonconv.Option{
		FilterMode: jsonconv.IncludeMode,
		FilterList: []string{"nul_name", "nul_time", "nul_float"} })
	log.Debug("json result: %s", ret)

	// test exclude mode
	ret, _ = jsonconv.SqlToJson(item, jsonconv.Option{
		FilterMode: jsonconv.ExcludeMode,
		FilterList: []string{"nul_name", "nul_time"} })
	log.Debug("json result: %s", ret)

	// ensure ascii
	ret, _ = jsonconv.SqlToJson(&item, jsonconv.Option{EnsureAscii: true})
	log.Debug("json result: %s", ret)
	return
}
