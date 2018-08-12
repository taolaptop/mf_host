package main

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/user/log"
	"github.com/user/middleware/security"
	"github.com/user/netty/udp"
	"github.com/user/service"
	"github.com/user/stringutil"
)

func main() {

	// type Serverslice2 struct {
	// 	Total int    `json:"total"`
	// 	Count int    `json:"count"`
	// 	Start int    `json:"start"`
	// 	V     string `json:"v"`
	// }

	// var s Serverslice2
	// str := `{"count": 10,"start": 0,"total": 1626,"v":"value"}`
	// err := json.Unmarshal([]byte(str), &s)
	// fmt.Println(s, err)
	// fmt.Printf("%+v", s)
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.Use(security.Filter)
	log.Info("hello")
	service.AddLocHandlerTo(router)
	service.AddCommonHandlerTo(router)
	service.AddLoginHandlerTo(router)
	fmt.Printf(stringutil.Reverse("!oG ,olleH"))
	go udp.Host()
	router.Run(":8080")
}
