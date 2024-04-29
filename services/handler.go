package services

import (
	"fmt"
	"mysql-mongodb-syncer/global"
	"mysql-mongodb-syncer/syncer"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type Handler struct {
	rowsEvents chan *canal.RowsEvent
	stop       chan any
}

func NewHandler() *Handler {
	return &Handler{
		rowsEvents: make(chan *canal.RowsEvent, 5000),
	}
}

func (h *Handler) OnRotate(*replication.EventHeader, *replication.RotateEvent) error {
	return nil
}
func (h *Handler) OnTableChanged(*replication.EventHeader, string, string) error {
	return nil
}
func (h *Handler) OnDDL(*replication.EventHeader, mysql.Position, *replication.QueryEvent) error {
	return nil
}

func (h *Handler) OnRow(rowsEvent *canal.RowsEvent) error {
	if _, ok := global.RulesMap[getSchemaTable(rowsEvent)]; ok {
		h.rowsEvents <- rowsEvent
	}
	return nil
}

func (h *Handler) OnXID(*replication.EventHeader, mysql.Position) error { return nil }
func (h *Handler) OnGTID(*replication.EventHeader, mysql.GTIDSet) error { return nil }
func (h *Handler) OnPosSynced(*replication.EventHeader, mysql.Position, mysql.GTIDSet, bool) error {
	return nil
}

func (h *Handler) String() string { return "mysql-mongodb-handler" }

func (h *Handler) Start() {
	go func() {
		for rowsEvent := range h.rowsEvents {

			for _, rule := range global.RulesMap[getSchemaTable(rowsEvent)] {
				switch rule.Target {
				case global.TargetMongoDB:
					syncer.MongoInstance.Sync(rowsEvent, rule)
				}
			}
		}

	}()

}

func getSchemaTable(rowsEvent *canal.RowsEvent) string {
	return fmt.Sprintf("%s.%s", rowsEvent.Table.Schema, rowsEvent.Table.Name)
}
