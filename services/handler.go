package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type updateQueue struct {
	mu   sync.Mutex
	data []string
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
	fmt.Printf("+++++++OnRow+++++++++")
	fmt.Printf("%+v\n", rowsEvent)
	// 根据配置构造一个规则表，当event命中，执行插入MongoDB的操作
	h.updateQueue.data = append(h.updateQueue.data, "row")
	h.updateQueue.mu.Unlock()

	return nil
}
func (h *Handler) OnXID(*replication.EventHeader, mysql.Position) error { return nil }
func (h *Handler) OnGTID(*replication.EventHeader, mysql.GTIDSet) error { return nil }
func (h *Handler) OnPosSynced(*replication.EventHeader, mysql.Position, mysql.GTIDSet, bool) error {
	return nil
}

func (h *Handler) String() string { return "mysql-mongodb-handler" }

func (h *Handler) Start() error {
	go func() {
		for {
			h.updateQueue.mu.Lock()
			if len(h.updateQueue.data) > 0 {
				fmt.Println("Received row:", h.updateQueue.data[0])
				h.updateQueue.data = h.updateQueue.data[1:]
			}
			h.updateQueue.mu.Unlock()
			// 如何控制，如果lock时候会不会丢失更新
			time.Sleep(1 * time.Second)
		}

	}()
	return nil
}
