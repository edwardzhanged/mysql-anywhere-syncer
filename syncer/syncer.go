package syncer

import "github.com/go-mysql-org/go-mysql/canal"

type Syncer interface {
	Connect(...any) error
	Ping() error
	Sync(*canal.RowsEvent) error
	Close() error
}
