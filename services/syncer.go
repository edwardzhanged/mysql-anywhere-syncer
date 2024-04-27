package services

type Syncer interface {
	Client() error
	Sync() error
}
