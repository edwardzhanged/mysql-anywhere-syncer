package services

import (
	"log"
	"mysql-mongodb-syncer/global"

	"github.com/go-mysql-org/go-mysql/canal"
)

type Listener struct {
	canal   *canal.Canal
	handler *Handler
}

var ListenerService *Listener

func InitCanal() *canal.Canal {
	canalCfg := canal.NewDefaultConfig()
	canalCfg.Addr = global.GbConfig.Addr
	canalCfg.User = global.GbConfig.User
	canalCfg.Password = global.GbConfig.Password
	canalCfg.Charset = global.GbConfig.Charset
	canalCfg.ServerID = global.GbConfig.SlaveID

	c, err := canal.NewCanal(canalCfg)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func NewListener() {
	c := InitCanal()
	ListenerService = &Listener{}
	Handler := NewHandler()
	ListenerService.canal = c
	ListenerService.handler = Handler
}

func (l *Listener) Start() {
	l.canal.SetEventHandler(l.handler)
	l.handler.Start()
	l.canal.GetMasterPos()
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// err = l.canal.RunFrom(pos)
	// if err != nil {
	// 	fmt.Print(err)
	// }
}

func (l *Listener) Reload() {
	l.canal.Close()
	l.canal = InitCanal()
}
