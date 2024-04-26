package services

import (
	"log"
	"mysql-mongodb-syncer/global"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
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
	err := l.canal.RunFrom(mysql.Position{})
	if err != nil {
		log.Fatal(err)
	}
}

func (l *Listener) Reload() {
	l.canal.Close()
	l.canal = InitCanal()
}
