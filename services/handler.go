package services

import (
	"fmt"
	"mysql-mongodb-syncer/global"
	"mysql-mongodb-syncer/syncer"
	"mysql-mongodb-syncer/utils/logger"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type Handler struct {
	rowsEvents chan *canal.RowsEvent
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
		logger.Logger.Infof("Received rows event: %s", rowsEvent)
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
	// Create a buffered channel to hold the jobs
	jobs := make(chan *canal.RowsEvent, 50)

	// Start 50 workers
	for w := 1; w <= 50; w++ {
		go func(w int, jobs <-chan *canal.RowsEvent) {
			for rowsEvent := range jobs {
				for _, rule := range global.RulesMap[getSchemaTable(rowsEvent)] {
					switch rule.Target {
					case global.TargetMongoDB:
						err := syncer.MongoInstance.Sync(rowsEvent, rule)
						if err != nil {
							logger.Logger.WithError(err).Error("Failed to sync to mongodb")
						}
					}
				}
			}
		}(w, jobs)
	}

	// Send jobs to the workers
	go func() {
		for rowsEvent := range h.rowsEvents {
			jobs <- rowsEvent
		}
		close(jobs)
	}()

}

func getSchemaTable(rowsEvent *canal.RowsEvent) string {
	return fmt.Sprintf("%s.%s", rowsEvent.Table.Schema, rowsEvent.Table.Name)
}
