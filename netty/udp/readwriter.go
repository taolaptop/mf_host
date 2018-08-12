package udp

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

import (
	"github.com/AlexStocks/getty"
	log "github.com/AlexStocks/log4go"
	"encoding/hex"
)

type EchoPackageHandler struct {
}

func NewEchoPackageHandler() *EchoPackageHandler {
	return &EchoPackageHandler{}
}

func (h *EchoPackageHandler) Read(ss getty.Session, data []byte) (interface{}, int, error) {
	var (
		err error
		len int
		pkg IPackage
		buf *bytes.Buffer
	)
	switch data[2] {
	case 10:
		pkg = &LatlngMessage{}
	case 22:
		pkg = &EchoPackage{}
	default:
		return nil, 0, errors.New("protocol err:"+  hex.EncodeToString(data))
	}
	buf = bytes.NewBuffer(data)
	len, err = pkg.Unmarshal(buf)
	if err != nil {
		if err == ErrNotEnoughStream {
			return nil, 0, nil
		}

		return nil, 0, err
	}

	return pkg, len, nil
}

func (h *EchoPackageHandler) Write(ss getty.Session, udpCtx interface{}) error {
	var (
		ok        bool
		err       error
		startTime time.Time
		echoPkg   *EchoPackage
		buf       *bytes.Buffer
		ctx       getty.UDPContext
	)

	ctx, ok = udpCtx.(getty.UDPContext)
	if !ok {
		log.Error("illegal UDPContext{%#v}", udpCtx)
		return fmt.Errorf("illegal @udpCtx{%#v}", udpCtx)
	}

	startTime = time.Now()
	if echoPkg, ok = ctx.Pkg.(*EchoPackage); !ok {
		log.Error("illegal pkg:%+v, addr:%s\n", ctx.Pkg, ctx.PeerAddr)
		return errors.New("invalid echo package!")
	}

	buf, err = echoPkg.Marshal()
	if err != nil {
		log.Warn("binary.Write(echoPkg{%#v}) = err{%#v}", echoPkg, err)
		return err
	}

	_, err = ss.Write(getty.UDPContext{Pkg: buf.Bytes(), PeerAddr: ctx.PeerAddr})
	log.Info("WriteEchoPkgTimeMs = %s", time.Since(startTime).String())

	return err
}
