package security

import (
	"net/http"

	"github.com/user/entity"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"regexp"
)

func Filter(ctx *gin.Context) {

	b,_:=regexp.MatchString("^/crc[\\?|/].*", ctx.Request.RequestURI)
	if path := ctx.Request.RequestURI;
		path == "/login" ||
			path == "/register" ||
			path == "/location" ||
			b {
		ctx.Next()

	} else {
		sessions := sessions.Default(ctx)
		if sessions.Get("account")==nil {
			ctx.JSON(http.StatusMethodNotAllowed, gin.H{"status": 403})
			ctx.Abort()
		}
		session := sessions.Get("account").(*entity.Account)
		if session.Created == 0 {

			//if session.GetAccount().Created == 0 {
			ctx.JSON(http.StatusMethodNotAllowed, gin.H{"status": 403})
			ctx.Abort()
			//}
		}
		ctx.Next()
	}
}
