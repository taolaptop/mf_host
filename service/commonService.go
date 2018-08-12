package service

import (
	"github.com/gin-gonic/gin"
	"encoding/hex"
	"github.com/AlexStocks/log4go"
	"net/http"
	"bytes"
	"encoding/binary"
	"github.com/user/util"
)

//整形转换成字节
func IntToBytes(n uint16) []byte {
	tmp := uint16(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int(tmp)
}

func AddCommonHandlerTo(engine *gin.Engine) {

	engine.GET("/crc", CrcHandler)

}

func CrcHandler(ctx *gin.Context)  {
	str,_:=ctx.GetQuery("hex")
	bts,err :=hex.DecodeString(str)

	if err != nil {
		log4go.Info("err:"+err.Error())
	}
	crc:=util.CheckSum(bts)


	ctx.JSON(http.StatusOK,gin.H{"crc16":hex.EncodeToString(IntToBytes(crc))})
}