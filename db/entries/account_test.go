package entries

import (
	"cronsun/db"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestGetAccounts(t *testing.T) {
	InitTestDb()
	if data, err := GetAccounts(bson.M{}); err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}
func TestCreateAccount(t *testing.T) {
	InitTestDb()
	if err := CreateAccount(&Account{
		Password:     "123456",
		Status:       UserActived,
		Unchangeable: true,
		Salt:         "123456",
		Email:        "weiguoxu@outlook.com",
		Role:         Administrator,
	}); err != nil {
		t.Error(err)
	}
}

func TestGetAccountByEmail(t *testing.T) {
	InitTestDb()
	if data, err := GetAccountByEmail("weiguoxu@outlook.com"); err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestUpdateAccount(t *testing.T) {
	InitTestDb()
	if err := UpdateAccount(bson.M{"email": "weiguoxu@outlook.com"}, bson.M{"role": Developer}); err != nil {
		t.Error(err)
	}
}

func TestBanAccount(t *testing.T) {
	InitTestDb()
	if err := BanAccount("weiguoxu@outlook.com"); err != nil {
		t.Error(err)
	}
}
func TestEnsureAccountIndex(t *testing.T) {
	InitTestDb()
	if err := EnsureAccountIndex(); err != nil {
		t.Error(err)
	}
}

func getConfig() *db.Config {
	return &db.Config{
		Database:   "cronsun",
		Hosts:      []string{"192.168.1.81:27018"},
		Timeout:    60,
		UserName:   "root",
		Password:   "123456",
		AuthSource: "admin",
	}
}

func InitTestDb() {
	mgoDB, err := db.NewMdb(getConfig())
	if err != nil {
		panic(err)
	} else {
		db.SetDb(mgoDB)
	}

}
