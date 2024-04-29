package syncer

import (
	"context"
	"fmt"
	"mysql-mongodb-syncer/global"
	"mysql-mongodb-syncer/utils/logger"
	"sync"

	"github.com/go-mysql-org/go-mysql/canal"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
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

	defer func() {
		if err != nil {
			if closeErr := client.Disconnect(context.TODO()); closeErr != nil {
				logger.Logger.WithError(closeErr).Error("Failed to close mongo client")
			}
		}
	}()

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	m.client = client
	return nil
}

func (m *Mongo) Ping() error {
	return m.client.Ping(context.Background(), nil)
}

type Document struct {
	Data []interface{} `bson:"data"`
}

func (m *Mongo) Sync(rowsEvent *canal.RowsEvent, rule *global.Rule) error {
	collection := m.client.Database(rule.MongodbDatabase).Collection(rule.MongodbCollection)
	for _, row := range rowsEvent.Rows {
		switch rowsEvent.Action {
		case canal.DeleteAction:
			id := rowsEvent.Table.PKColumns[0]
			_, err := collection.DeleteOne(context.Background(), bson.M{id: id})
			if err != nil {
				logger.Logger.WithError(err).Error("Failed to delete row")
			}
		case canal.InsertAction:
			doc := &Document{
				Data: row,
			}
			_, err := collection.InsertOne(context.Background(), doc)
			if err != nil {
				logger.Logger.WithError(err).Error("Failed to insert row")

			}
		}

	}
	return nil
}

func (m *Mongo) Close() error {
	return m.client.Disconnect(context.Background())
}
