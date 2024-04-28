package syncer

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConnectOptions struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Mongo struct {
	connectOptions *ConnectOptions

	client *mongo.Client
	lock   sync.Mutex
}

var MongoInstance *Mongo

func NewMongo(connectOptions *ConnectOptions) {
	MongoInstance = &Mongo{
		connectOptions: connectOptions,
	}
}

func (m *Mongo) Connect() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%d", m.connectOptions.Username,
		m.connectOptions.Password, m.connectOptions.Host, m.connectOptions.Port)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	if err != nil {
		return err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	m.client = client
	return nil
}

func (m *Mongo) Ping() error {
	return m.client.Ping(context.Background(), nil)
}

func (m *Mongo) Sync() error {
	return nil
}

func (m *Mongo) Close() error {
	return m.client.Disconnect(context.Background())
}
