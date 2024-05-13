package syncer

import (
	"github.com/go-mysql-org/go-mysql/canal"

	"mysql-anywhere-syncer/global"
)

type Syncer interface {
	Connect(...any) error
	Ping() error
	Sync(*canal.RowsEvent, *global.Rule) error
	Close() error
}
