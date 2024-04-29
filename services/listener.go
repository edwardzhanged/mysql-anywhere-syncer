package services

import (
	"mysql-mongodb-syncer/global"
	"mysql-mongodb-syncer/utils/logger"

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
		logger.Logger.WithError(err).Fatal("Failed to create canal")
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
	pos, err := l.canal.GetMasterPos()
	if err != nil {
		logger.Logger.WithError(err).Fatal("Failed to get master pos")
	}
	go func() { l.canal.RunFrom(pos) }()

}

func (l *Listener) Reload() {
	l.canal.Close()
	l.canal = InitCanal()
	l.Start()
}
