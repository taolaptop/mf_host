package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/AlexStocks/goext/time"
	"fmt"
	"strings"
	"reflect"
	"time"
	"github.com/user/dao"
	"github.com/user/entity"
)

var (
	StrQ  = make(chan string)
	wheel = gxtime.NewWheel(gxtime.TimeMillisecondDuration(100), 1200) // wheel longest span is 2 minute

)

type Location struct {
	Mobile  string
	Lat     string
	Lng     string
	Created int64
}

func AddLocHandlerTo(engine *gin.Engine) {
	engine.GET("/location", LocHandler)
	go observer()
}

func LocHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"loc": <-StrQ})
}

func (l *Location) parse(s string) *Location {
	arr := strings.Split(s, ",")

	//ss := reflect.ValueOf(l).Elem()
	//typeOfT := ss.Type()
	//for i := 0; i < ss.NumField(); i++ {
	//	f := ss.Field(i)
	//	fmt.Printf("%d:%v %s %s = %v\n", i, ss.Kind(), typeOfT.Field(i).Name, f.Type(), f.Interface())
	//}

	//fmt.Println(reflect.ValueOf(ss.Field(0)).NumField())

	fmt.Println(reflect.ValueOf(l).Elem().NumField())
	fmt.Println(reflect.TypeOf(l))
	for _, item := range arr {
		kv := strings.Split(item, ":")
		k := kv[0]
		v := kv[1]
		//fmt.Println(reflect.ValueOf(l).Elem())
		fmt.Println(reflect.ValueOf(l).Elem().FieldByName(k))
		reflect.ValueOf(l).Elem().FieldByName(k).SetString(v)

	}

	fmt.Printf("%#v", l)
	return l
}

func observer() {
	var str string
	var idb dao.IAccount
	db := new(dao.Account)
	idb = db


	var ilocdb dao.ILocation
	locdb:= dao.Location{}
	ilocdb = locdb

	for {

		str = <-StrQ

		//fmt.Println(str)
		l := Location{}

		l.parse(str)
		l.Created = time.Now().Unix()

		account := idb.Get(l.Mobile)

		location := entity.LocationSnapshot{Id: account.Id, Lng: l.Lng, Lat: l.Lat, Updated: l.Created}

		ilocdb.Save(location)
		//select {
		//case <-Str:
		//	break
		//case <-wheel.After(60 * 1e9):
		//	fmt.Println("quit")
		//	break
		//}

	}

}
