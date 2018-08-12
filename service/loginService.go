package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/user/dao"
	"github.com/user/entity"
	"github.com/user/log"
	"github.com/user/util"
)

func AddLoginHandlerTo(engine *gin.Engine) {

	engine.POST("/login", LoginHandler)
	engine.POST("/register", RegisterHandler)

}

func RegisterHandler(ctx *gin.Context) {
	var a, _ = ctx.GetRawData()
	fmt.Println(string(a))
	var account entity.Account
	//var authJson = string(a)
	var err = json.Unmarshal(a, &account)
	if err != nil {
		fmt.Println("error:", err)
	}

	account.Pwd = util.Md5(account.Pwd)
	account.Created = time.Now().Unix()
	var idb dao.IAccount
	db := new(dao.Account)
	idb = db
	idb.Add(account)
	ctx.JSON(http.StatusOK, gin.H{"status": "200"})

}

func LoginHandler(ctx *gin.Context) {
	var a, e = ctx.GetRawData()
	fmt.Println(string(a))

	log.Info("hello")
	var account entity.Account
	//var authJson = string(a)
	var err = json.Unmarshal(a, &account)
	if err != nil {
		fmt.Println("error:", err)
	}
	//panic(e)
	//log.Info(account.account)
	fmt.Printf("%+v", account)

	account.Pwd = util.Md5(account.Pwd)

	var idb dao.IAccount
	db := new(dao.Account)
	idb = db
	u := idb.Auth(account)
	fmt.Printf("%+v", u)
	if e != nil {
		defer fmt.Print(e.Error())
	}
	if u.Created > 0 {
		session := sessions.Default(ctx)
		session.Set("account", u)
		session.Save()

		ctx.JSON(http.StatusOK, gin.H{"status": 200})

	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 404})

	}

}
