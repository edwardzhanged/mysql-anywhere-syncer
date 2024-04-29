package syncer

import (
	"github.com/go-mysql-org/go-mysql/canal"

	"mysql-mongodb-syncer/global"
)

type Syncer interface {
	Connect(...any) error
	Ping() error
	Sync(*canal.RowsEvent, *global.Rule) error
	Close() error
}
