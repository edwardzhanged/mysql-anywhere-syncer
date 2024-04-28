package services

import (
	"fmt"
	"mysql-mongodb-syncer/global"
	"mysql-mongodb-syncer/syncer"
	"sync"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type updateQueue struct {
	mu         sync.Mutex
	rowsEvents []*canal.RowsEvent
}

type Handler struct {
	updateQueue updateQueue
	stop        chan any
}

func NewHandler() *Handler {
	return &Handler{}
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

// table命中config中的表，才会触发OnRow
func (h *Handler) OnRow(rowsEvent *canal.RowsEvent) error {
	h.updateQueue.mu.Lock()
	schemaTable := fmt.Sprintf("%s.%s", rowsEvent.Table.Schema, rowsEvent.Table.Name)
	if _, ok := global.RulesMap[schemaTable]; ok {
		h.updateQueue.rowsEvents = append(h.updateQueue.rowsEvents, rowsEvent)
		h.updateQueue.mu.Unlock()
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
	// TODO: 要解决并发问题
	go func() {
		for {
			h.updateQueue.mu.Lock()
			if len(h.updateQueue.rowsEvents) > 0 {
				fmt.Println("Received row:", h.updateQueue.rowsEvents[0])
				rowEvent := h.updateQueue.rowsEvents[0]
				for _, rule := range global.RulesMap[rowEvent.Table.Schema+"."+rowEvent.Table.Name] {
					switch rule.Target {
					case global.TargetMongoDB:
						syncer.MongoInstance.Sync()
						fmt.Println("Sync to MongoDBcc")
					default:
						fmt.Println("Unknown target")
					}
				}
				h.updateQueue.rowsEvents = h.updateQueue.rowsEvents[1:]
			}
			h.updateQueue.mu.Unlock()
			// 如何控制，如果lock时候会不会丢失更新
			time.Sleep(1 * time.Second)
		}

	}()
}
