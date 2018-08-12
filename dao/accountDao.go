package dao

import (
	"fmt"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"github.com/user/entity"

	"github.com/satori/go.uuid"
)

var db *xorm.Engine

type Account struct {
}

type IAccount interface {
	GetSource() *xorm.Engine
	Add(entity.Account) *entity.Account
	Get(account string) entity.Account
	Auth(account entity.Account) entity.Account
}

func init() {
	db = new(Account).GetSource()
}

func (s *Account) GetSource() *xorm.Engine {
	// masterDataSourceName := "postgres"
	// //driverName := "postgres://postgres:postgres@localhost:5432/huafengming?sslmode=disable"
	// driverName := "host=localhost port=5432 user=postgres password=postgres dbname=huafengming sslmode=disable"
	// dataSourceNameSlice := []string{masterDataSourceName}
	// engineGroup, err := xorm.NewEngineGroup(driverName, dataSourceNameSlice)
	//1.创建db引擎
	db, err := xorm.NewEngine("postgres", "postgres://postgres:postgres@localhost:5432/huafengming?sslmode=disable")
	if err != nil {
		defer fmt.Println(err)
	}

	//2.显示sql语句
	db.ShowSQL(true)

	//3.设置连接数
	db.SetMaxIdleConns(2000)
	db.SetMaxOpenConns(1000)

	cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 500) //缓存的条数
	db.SetDefaultCacher(cacher)
	db.SetSchema("mf")
	fmt.Println(db.Ping())
	// engineGroup.SetSchema("mf")
	return db
}

func (s *Account) Add(account entity.Account) *entity.Account {
	//results, err := engineGroup.Insert(&account)
	//results, err := engineGroup.Where("a = 1").Query()
	//db := s.GetSource()
	account.Id = uuid.Must(uuid.NewV4()).String()
	_, err := db.InsertOne(account)
	if err != nil {
		fmt.Println(err)
	}
	return &account
}

func (source *Account) Get(s string) entity.Account {
	account := entity.Account{}
	//db := source.GetSource()
	db.Where("account=?", s).Get(&account)
	return account
}

func (source *Account) Auth(account entity.Account) entity.Account {
	// accountdb := entity.Account{}
	//db := source.GetSource()
	has, _ := db.Where("account=?", account.Account).And("pwd=?", account.Pwd).Get(&account)
	if has {
		return account
	}
	return entity.Account{}
}
