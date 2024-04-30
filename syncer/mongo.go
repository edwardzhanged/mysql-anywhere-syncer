package syncer

import (
	"context"
	"errors"
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
	if len(rowsEvent.Table.PKColumns) == 0 {
		return errors.New("no primary key column")
	}
	collection := m.client.Database(rule.MongodbDatabase).Collection(rule.MongodbCollection)
	switch rowsEvent.Action {
	case canal.DeleteAction:
		oldRow := rowsEvent.Rows[0]
		oldPk := getPrimaryKey(oldRow, rowsEvent)
		_, err := collection.DeleteOne(context.Background(), bson.M{"_id": oldPk})
		if err != nil {
			logger.Logger.WithError(err).Error("Failed to delete row")
		}
	case canal.InsertAction:
		// TODO: 忽略的列
		newRow := rowsEvent.Rows[0]
		newPk := getPrimaryKey(newRow, rowsEvent)

		doc := bson.M{"_id": newPk}
		for i, column := range rowsEvent.Table.Columns {
			doc[column.Name] = newRow[i]
		}
		_, err := collection.InsertOne(context.Background(), doc)
		if err != nil {
			logger.Logger.WithError(err).Error("Failed to insert row")
		}
	case canal.UpdateAction:
		oldRow, newRow := rowsEvent.Rows[0], rowsEvent.Rows[1]
		oldPK, newPK := getPrimaryKey(oldRow, rowsEvent), getPrimaryKey(newRow, rowsEvent)
		newDoc := bson.M{"_id": newPK}
		for i, column := range rowsEvent.Table.Columns {
			newDoc[column.Name] = newRow[i]
		}
		_, err := collection.DeleteOne(context.Background(), bson.M{"_id": oldPK})
		if err != nil {
			logger.Logger.WithError(err).Error("Failed to delete row")
			return err
		}
		_, err = collection.InsertOne(context.Background(), newDoc)
		if err != nil {
			logger.Logger.WithError(err).Error("Failed to insert row")
		}
	}
	return nil
}

func (m *Mongo) Close() error {
	return m.client.Disconnect(context.Background())
}

func getPrimaryKey(row []interface{}, rowsEvent *canal.RowsEvent) (pk string) {

	for i, PKColumn := range rowsEvent.Table.PKColumns {
		if i != 0 {
			pk += "_"
		}
		switch v := row[PKColumn].(type) {
		case int32:
			pk += fmt.Sprintf("%d", v)
		case string:
			pk += v
		default:
			logger.Logger.WithField("pk", pk).Error("Unsupported type for primary key")
		}
		logger.Logger.WithField("pk", pk).Info("Primary key")
	}
	return pk
}
