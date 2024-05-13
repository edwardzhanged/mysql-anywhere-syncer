package syncer

import (
	"context"
	"errors"
	"fmt"
	"mysql-anywhere-syncer/global"
	"mysql-anywhere-syncer/utils/logger"
	"strconv"
	"sync"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/schema"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type MongoConnectOptions struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Mongo struct {
	connectOptions *MongoConnectOptions

	client *mongo.Client
	lock   sync.Mutex
}

var MongoInstance *Mongo

func NewMongo(connectOptions *MongoConnectOptions) {
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
		logger.Logger.WithFields(logrus.Fields{
			"pk":     oldPk,
			"oldrow": oldRow,
		}).Info("Delted row")
	case canal.InsertAction:
		newRow := rowsEvent.Rows[0]
		newPk := getPrimaryKey(newRow, rowsEvent)

		doc := bson.M{"_id": newPk}
		buildUpsertDoc(doc, newRow, rowsEvent, rule)
		_, err := collection.InsertOne(context.Background(), doc)
		if err != nil {
			logger.Logger.WithError(err).Error("Failed to insert row")
		}
		logger.Logger.WithFields(logrus.Fields{
			"pk":     newPk,
			"newrow": newRow,
		}).Info("Inserted row")
	case canal.UpdateAction:
		oldRow, newRow := rowsEvent.Rows[0], rowsEvent.Rows[1]
		oldPK, newPK := getPrimaryKey(oldRow, rowsEvent), getPrimaryKey(newRow, rowsEvent)
		newDoc := bson.M{"_id": newPK}
		buildUpsertDoc(newDoc, newRow, rowsEvent, rule)
		_, err := collection.DeleteOne(context.Background(), bson.M{"_id": oldPK})
		if err != nil {
			logger.Logger.WithError(err).Error("Failed to delete row")
			return err
		}
		_, err = collection.InsertOne(context.Background(), newDoc)
		if err != nil {
			logger.Logger.WithError(err).Error("Failed to insert row")
		}
		logger.Logger.WithFields(logrus.Fields{
			"oldpk":  oldPK,
			"newpk":  newPK,
			"oldrow": oldRow,
			"newrow": newRow,
		}).Info("Updated row")
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
		case int64:
			pk += fmt.Sprintf("%d", v)
		case float64:
			pk += fmt.Sprintf("%f", v)
		case []byte:
			pk += string(v)
		case time.Time:
			pk += v.Format(time.RFC3339)
		case string:
			pk += v
		default:
			logger.Logger.WithField("pk", pk).Error("Unsupported type for primary key")
		}
	}
	return pk
}

func buildUpsertDoc(doc bson.M, newRow []interface{}, rowsEvent *canal.RowsEvent, rule *global.Rule) {
	for i, column := range rowsEvent.Table.Columns {
		if len(rule.IncludeColumnsConfig) > 0 {
			for _, incluedColumn := range rule.IncludeColumnsConfig {
				if column.Name == incluedColumn {
					addConvertedValue(doc, column, newRow[i])
				}
			}
		} else {
			addConvertedValue(doc, column, newRow[i])
		}
		if len(rule.ExcludeColumnsConfig) > 0 {
			for _, excludeColumn := range rule.ExcludeColumnsConfig {
				if column.Name == excludeColumn {
					delete(doc, column.Name)
				}
			}
		}
		if len(rule.ColumnMappingsConfig) > 0 {
			for _, columnMapping := range rule.ColumnMappingsConfig {
				if column.Name == columnMapping.Source {
					doc[columnMapping.Target] = newRow[i]
				}
			}
		}
	}
	if len(rule.NewColumnsConfig) > 0 {
		for _, newColumn := range rule.NewColumnsConfig {
			if newColumn.Type == "int" {
				intValue, err := strconv.Atoi(newColumn.Value)
				if err != nil {
					logger.Logger.WithError(err).Error("Failed to convert string to int")
					continue
				}
				doc[newColumn.Name] = intValue
			} else if newColumn.Type == "string" {
				doc[newColumn.Name] = newColumn.Value
			} else if newColumn.Type == "float" {
				floatValue, err := strconv.ParseFloat(newColumn.Value, 64)
				if err != nil {
					logger.Logger.WithError(err).Error("Failed to convert string to float")
					continue
				}
				doc[newColumn.Name] = floatValue
			} else if newColumn.Type == "bool" {
				boolValue, err := strconv.ParseBool(newColumn.Value)
				if err != nil {
					logger.Logger.WithError(err).Error("Failed to convert string to bool")
					continue
				}
				doc[newColumn.Name] = boolValue
			} else {
				logger.Logger.Error("Unsupported type for new column")
			}
		}
	}
}

func addConvertedValue(doc bson.M, column schema.TableColumn, newRow any) {

	switch column.Type {
	case schema.TYPE_TIMESTAMP, schema.TYPE_DATETIME, schema.TYPE_TIME:
		t, _ := time.Parse("2006-01-02 15:04:05", newRow.(string))
		doc[column.Name] = t
	case schema.TYPE_DATE:
		t, _ := time.Parse("2006-01-02", newRow.(string))
		doc[column.Name] = t
	case schema.TYPE_FLOAT:
		doc[column.Name] = newRow.(float64)
	default:
		doc[column.Name] = newRow
	}

}
