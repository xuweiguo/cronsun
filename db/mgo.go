package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	Hosts []string
	// AuthSource Specify the database name associated with the user’s credentials.
	// authSource defaults to the database specified in the connection string.
	AuthSource string
	UserName   string
	Password   string
	Database   string
	Timeout    time.Duration // second
}

type Mdb struct {
	*Config
	*mongo.Client
}

func NewMdb(c *Config) (*Mdb, error) {
	m := &Mdb{
		Config: c,
	}
	return m, m.connect()
}

func (m *Mdb) connect() error {
	// connectionString: [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	// via: https://docs.mongodb.com/manual/reference/connection-string/
	connectionString := strings.Join(m.Config.Hosts, ",")
	if len(m.Config.UserName) > 0 && len(m.Config.Password) > 0 {
		connectionString = m.Config.UserName + ":" + url.QueryEscape(m.Config.Password) + "@" + connectionString
	}

	if len(m.Config.Database) > 0 {
		connectionString += "/" + m.Config.Database
	}

	if len(m.Config.AuthSource) > 0 {
		connectionString += "?authSource=" + m.Config.AuthSource
	}

	ctx, _ := context.WithTimeout(context.Background(), m.Config.Timeout)
	// 创建MongoDB客户端
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+connectionString))
	if err != nil {
		return err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}
	m.Client = client
	return nil
}

func (m *Mdb) WithC(collection string, job func(*mongo.Collection) error) error {
	database := m.Client.Database(m.Config.Database)
	collectionObj := database.Collection(collection)

	// 创建会话选项
	sessionOptions := options.Session()
	sessionOptions.SetCausalConsistency(true) // 设置原因一致性

	// 开始会话
	session, err := m.Client.StartSession(sessionOptions)
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	// 在会话中执行操作
	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		err := job(collectionObj)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (self *Mdb) Upsert(collection string, selector interface{}, change interface{}) error {
	return self.WithC(collection, func(c *mongo.Collection) error {
		// 执行upsert操作
		opts := options.Update().SetUpsert(true)
		_, err := c.UpdateOne(context.Background(), selector, change, opts)
		return err
	})
}

func (self *Mdb) Insert(collection string, data ...interface{}) error {
	return self.WithC(collection, func(c *mongo.Collection) error {
		_, err := c.InsertMany(context.Background(), data)
		return err
	})
}

func (self *Mdb) FindId(collection string, id string, result interface{}) error {
	return self.WithC(collection, func(c *mongo.Collection) error {
		idHex, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		// 构造查询
		filter := bson.M{"_id": idHex}

		err = c.FindOne(context.Background(), filter).Decode(result)
		return err
	})
}

func (self *Mdb) FindOne(collection string, query interface{}, result interface{}) error {
	return self.WithC(collection, func(c *mongo.Collection) error {
		return c.FindOne(context.Background(), query).Decode(result)
	})
}

func (self *Mdb) RemoveId(collection string, id string) error {
	return self.WithC(collection, func(c *mongo.Collection) error {
		idHex, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		// 构造查询
		filter := bson.M{"_id": idHex}

		_, err = c.DeleteOne(context.Background(), filter)
		return err
	})
}
