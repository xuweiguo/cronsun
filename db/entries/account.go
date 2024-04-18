package entries

import (
	"context"
	"cronsun/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	Coll_Account = "cronsun_account"
)

type Account struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Role     Role               `bson:"role" json:"role"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
	Salt     string             `bson:"salt" json:"salt"`
	Status   UserStatus         `bson:"status" json:"status"`
	Session  string             `bson:"session" json:"-"`
	// If true, role and status are unchangeable, email and password can be change by it self only.
	Unchangeable bool      `bson:"unchangeable" json:"-"`
	CreateTime   time.Time `bson:"createTime" json:"createTime"`
}

type Role int

const (
	Administrator Role = 1
	Developer     Role = 2
	Reporter      Role = 3
)

func (r Role) Defined() bool {
	switch r {
	case Administrator, Developer, Reporter:
		return true
	}
	return false
}

func (r Role) String() string {
	switch r {
	case Administrator:
		return "Administrator"
	case Developer:
		return "Developer"
	case Reporter:
		return "Reporter"
	}
	return "Undefined"
}

type UserStatus int

const (
	UserBanned  UserStatus = -1
	UserActived UserStatus = 1
)

func (s UserStatus) Defined() bool {
	switch s {
	case UserBanned, UserActived:
		return true
	}
	return false
}

func GetAccounts(query bson.M) (list []Account, err error) {
	err = db.GetDb().WithC(Coll_Account, func(c *mongo.Collection) error {
		// 执行查询
		cursor, err := c.Find(context.Background(), query)
		if err != nil {
			return err
		}
		defer cursor.Close(context.Background())

		if err = cursor.All(context.Background(), &list); err != nil {
			return err
		}
		return nil
	})
	return
}

func GetAccountByEmail(email string) (u *Account, err error) {
	err = db.GetDb().FindOne(Coll_Account, bson.M{"email": email}, &u)
	return
}

func CreateAccount(u *Account) error {
	u.ID = primitive.NewObjectID()
	u.CreateTime = time.Now()
	return db.GetDb().Insert(Coll_Account, u)
}

func UpdateAccount(query bson.M, change bson.M) error {
	return db.GetDb().WithC(Coll_Account, func(c *mongo.Collection) error {
		_, err := c.UpdateMany(context.Background(), query, bson.M{"$set": change})
		return err
	})
}

func BanAccount(email string) error {
	return db.GetDb().WithC(Coll_Account, func(c *mongo.Collection) error {
		_, err := c.UpdateMany(context.Background(), bson.M{"email": email}, bson.M{"$set": bson.M{"status": UserBanned}})
		return err
	})
}

func EnsureAccountIndex() error {
	return db.GetDb().WithC(Coll_Account, func(c *mongo.Collection) error {
		// 定义唯一索引
		indexModel := mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		}
		// 创建唯一索引
		_, err := c.Indexes().CreateOne(context.Background(), indexModel)
		return err
	})
}
