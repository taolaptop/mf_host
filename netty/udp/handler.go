package udp

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

import (
	"github.com/AlexStocks/getty"
	log "github.com/AlexStocks/log4go"
	"reflect"
	"github.com/user/service"
)

const (
	WritePkgTimeout = 1e8
	WritePkgASAP    = 0e9
)

var (
	errTooManySessions = errors.New("Too many echo sessions!")
)

type PackageHandler interface {
	Handle(getty.Session, getty.UDPContext) error
}

////////////////////////////////////////////
// heartbeat handler
////////////////////////////////////////////

type HeartbeatHandler struct{}

func (h *HeartbeatHandler) Handle(session getty.Session, ctx getty.UDPContext) error {
	var (
		ok     bool
		pkg    *EchoPackage
		rspPkg EchoPackage
	)

	log.Debug("get echo heartbeat udp context{%#v}", ctx)
	if pkg, ok = ctx.Pkg.(*EchoPackage); !ok {
		return fmt.Errorf("illegal @ctx.Pkg:%#v", ctx.Pkg)
	}

	rspPkg.H = pkg.H
	rspPkg.B = echoHeartbeatResponseString
	rspPkg.H.Len = uint16(len(rspPkg.B) + 1)

	// return session.WritePkg(getty.UDPContext{Pkg: &rspPkg, PeerAddr: ctx.PeerAddr}, WritePkgTimeout)
	return session.WritePkg(getty.UDPContext{Pkg: &rspPkg, PeerAddr: ctx.PeerAddr}, WritePkgASAP)
}

////////////////////////////////////////////
// message handler
////////////////////////////////////////////

type MessageHandler struct{}

func (h *MessageHandler) Handle(session getty.Session, ctx getty.UDPContext) error {
	log.Debug("get echo ctx{%#v}", ctx)
	// write echo message handle logic here.
	// return session.WritePkg(ctx, WritePkgTimeout)
	return session.WritePkg(ctx, WritePkgASAP)
}

////////////////////////////////////////////
// EchoMessageHandler
////////////////////////////////////////////

type clientEchoSession struct {
	session getty.Session
	reqNum  int32
}

type EchoMessageHandler struct {
	handlers map[uint32]PackageHandler

	rwlock     sync.RWMutex
	sessionMap map[getty.Session]*clientEchoSession
}

func newEchoMessageHandler() *EchoMessageHandler {
	handlers := make(map[uint32]PackageHandler)
	handlers[heartbeatCmd] = &HeartbeatHandler{}
	handlers[echoCmd] = &MessageHandler{}

	return &EchoMessageHandler{sessionMap: make(map[getty.Session]*clientEchoSession), handlers: handlers}
}

func (h *EchoMessageHandler) OnOpen(session getty.Session) error {
	var (
		err error
	)

	h.rwlock.RLock()
	if conf.SessionNumber < len(h.sessionMap) {
		err = errTooManySessions
	}
	h.rwlock.RUnlock()
	if err != nil {
		return err
	}

	log.Info("got session:%s", session.Stat())
	h.rwlock.Lock()
	h.sessionMap[session] = &clientEchoSession{session: session}
	h.rwlock.Unlock()
	return nil
}

func (h *EchoMessageHandler) OnError(session getty.Session, err error) {
	log.Info("session{%s} got error{%v}, will be closed.", session.Stat(), err)
	h.rwlock.Lock()
	delete(h.sessionMap, session)
	h.rwlock.Unlock()
}

func (h *EchoMessageHandler) OnClose(session getty.Session) {
	log.Info("session{%s} is closing......", session.Stat())
	h.rwlock.Lock()
	delete(h.sessionMap, session)
	h.rwlock.Unlock()
}

func (h *EchoMessageHandler) OnMessage(session getty.Session, udpCtx interface{}) {
	ctx, ok := udpCtx.(getty.UDPContext)
	if !ok {
		log.Error("illegal UDPContext{%#v}", udpCtx)
		return
	}

	tp := reflect.TypeOf(ctx.Pkg)
	fmt.Println(tp.Kind().String())
	//switch tp.Kind() {
	//case reflect.Ptr:
	//vv:=	reflect.ValueOf(ctx.Pkg)
	//vvk:=vv.Kind()
	elm := tp.Elem()
	//
	//fmt.Println(vv)
	//	fmt.Println(vvk)
	fmt.Println(elm)
	////log id:1, sequence:1, command:heartbeat, echo string:you are welcome he
	////ptr
	////udp.EchoPackage
	//
	//}
	var
	(
		handler PackageHandler
		//ok1      bool
	)
	if elm == reflect.TypeOf(&LatlngMessage{}).Elem() {

		p, ok := ctx.Pkg.(*LatlngMessage)
		if !ok {
			log.Error("illegal pkg{%#v}", ctx.Pkg)
			return
		}
		service.StrQ <- p.String()
		fmt.Println(p.String())
		//handler, ok = h.handlers[p.H.Command]
		//if !ok {
		//	log.Error("illegal command{%d}", p.H.Command)
		//	return
		//}
	} else if elm == reflect.TypeOf(&EchoPackage{}).Elem() {

		p, ok := ctx.Pkg.(*EchoPackage)
		if !ok {
			log.Error("illegal pkg{%#v}", ctx.Pkg)
			return
		}

		handler, ok = h.handlers[p.H.Command]
		if !ok {
			log.Error("illegal command{%d}", p.H.Command)
			return
		}
		err := handler.Handle(session, ctx)
		if err != nil {
			h.rwlock.Lock()
			if _, ok := h.sessionMap[session]; ok {
				h.sessionMap[session].reqNum++
			}
			h.rwlock.Unlock()
		}
	}

}

func (h *EchoMessageHandler) OnCron(session getty.Session) {
	var (
		//flag   bool
		active time.Time
	)
	h.rwlock.RLock()
	if _, ok := h.sessionMap[session]; ok {
		active = session.GetActive()
		if conf.sessionTimeout.Nanoseconds() < time.Since(active).Nanoseconds() {
			//flag = true
			log.Error("session{%s} timeout{%s}, reqNum{%d}",
				session.Stat(), time.Since(active).String(), h.sessionMap[session].reqNum)
		}
	}
	h.rwlock.RUnlock()
	// udp session是根据本地udp socket fd生成的，如果关闭则连同socket也一同关闭了
	//if flag {
	//	h.rwlock.Lock()
	//	delete(h.sessionMap, session)
	//	h.rwlock.Unlock()
	//	session.Close()
	//}
}
