package udp

import (
	"fmt"
	"github.com/user/util"
	"bytes"
	"github.com/pkg/errors"
	"encoding/binary"
)

import (
	log "github.com/AlexStocks/log4go"
)

type IPackage interface {
	Marshal() (*bytes.Buffer, error)
	Unmarshal(buf *bytes.Buffer) (int, error)
	//Crc16()
}

var (
	latLngPkgHeaderLen int
	magic uint16=0x4d46
)

func init() {
	latLngPkgHeaderLen = 21//(int)((uint)(unsafe.Sizeof(LatlngMessage{})))
}

type LatlngMessage struct {
	Magic  uint16
	T      uint8 //10
	Mobile uint64
	Lat    uint32
	Lng    uint32
	Crc    uint16
}

func (m LatlngMessage) String() string {
	return fmt.Sprintf("Mobile:%d,Lat:%d,Lng:%d",
		m.Mobile, m.Lat, m.Lng)
}

func  Crc16(bt []byte) uint16 {

	checksum := util.CheckSum(bt) //调用计算CRC函数 CheckSum

	fmt.Printf("check sum:%X \n", checksum)
	return checksum
}

func (m *LatlngMessage) Marshal() (*bytes.Buffer, error) {
	return nil, nil
}

func (m *LatlngMessage)  Unmarshal(buf *bytes.Buffer) (int, error) {
	var (
		err error
		crc uint16
	)
	bts:=buf.Bytes()
	if buf.Len() < latLngPkgHeaderLen {
		return 0, errors.New("msg buffer len was too small")
	}
	err = binary.Read(buf, binary.BigEndian, m)
	if err != nil {
		return 0, err
	}
	if m.Magic != magic {
		log.Error("@Magic{%x}, right magic{%x}", m.Magic, magic)
		return 0, errors.New("magic is not right")
	}

	crc=Crc16(bts[:len(bts)-2])
	if m.Crc != crc {
		log.Error("@Crc{%x}, computed crc{%x}", m.Crc, crc)
		return 0 ,errors.New("crc validate fail")
	}

	return latLngPkgHeaderLen, nil
}
