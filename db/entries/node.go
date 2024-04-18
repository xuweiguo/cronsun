package entries

import (
	"context"
	"cronsun/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	Coll_Node = "cronsun_node"
)

type Node struct {
	ID       string `bson:"_id" json:"id"`  // machine id
	PID      string `bson:"pid" json:"pid"` // 进程 pid
	PIDFile  string `bson:"-" json:"-"`
	IP       string `bson:"ip" json:"ip"` // node ip
	Hostname string `bson:"hostname" json:"hostname"`

	Version  string    `bson:"version" json:"version"`
	UpTime   time.Time `bson:"up" json:"up"`     // 启动时间
	DownTime time.Time `bson:"down" json:"down"` // 上次关闭时间

	Alived    bool `bson:"alived" json:"alived"` // 是否可用
	Connected bool `bson:"-" json:"connected"`   // 当 Alived 为 true 时有效，表示心跳是否正常
}

func GetNodes() (nodes []*Node, err error) {
	return GetNodesBy(bson.M{})
}

func GetNodesBy(query interface{}) (nodes []*Node, err error) {
	err = db.GetDb().WithC(Coll_Node, func(c *mongo.Collection) error {
		find, err := c.Find(context.Background(), query)
		if err != nil {
			return err
		}
		return find.All(context.Background(), &nodes)
	})
	return
}

func GetNodesByID(id string) (node *Node, err error) {
	err = db.GetDb().FindOne(Coll_Node, bson.M{"_id": id}, &node)
	return
}

func RemoveNodeById(id string) error {
	query := bson.M{"_id": id}
	return db.GetDb().WithC(Coll_Node, func(c *mongo.Collection) error {
		_, err := c.DeleteMany(context.Background(), query)
		return err
	})
}
func RemoveNode(query interface{}) error {
	return db.GetDb().WithC(Coll_Node, func(c *mongo.Collection) error {
		_, err := c.DeleteMany(context.Background(), query)
		return err
	})
}

func SyncNodeToMgo(node *Node) error {
	return db.GetDb().Upsert(Coll_Node, bson.M{"_id": node.ID}, bson.M{"$set": node})
}
