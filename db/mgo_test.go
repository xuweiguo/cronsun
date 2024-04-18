package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestNewMdb(t *testing.T) {

	_, err := NewMdb(getConfig())
	if err != nil {
		t.Error(err)
	}
}

func TestMdb_Insert(t *testing.T) {
	m, err := NewMdb(getConfig())
	if err != nil {
		t.Error(err)
	}
	err = m.Insert("test", bson.M{"name": "Alice", "age": 30}, bson.M{"name": "Bob", "age": 35}, bson.M{"name": "Charlie", "age": 40})
	if err != nil {
		t.Error(err)
	}
}
func TestMdb_FindId(t *testing.T) {
	m, err := NewMdb(getConfig())
	if err != nil {
		t.Error(err)
	}
	var result = &bson.M{}
	err = m.FindId("test", "664c0a2ca232190f883acebb", result)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}
func TestMdb_FindOne(t *testing.T) {
	m, err := NewMdb(getConfig())
	if err != nil {
		t.Error(err)
	}
	var result = &bson.M{}
	err = m.FindOne("test", bson.M{"name": "Alice"}, result)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}
func TestMdb_RemoveId(t *testing.T) {
	m, err := NewMdb(getConfig())
	if err != nil {
		t.Error(err)
	}
	err = m.RemoveId("test", "664c0a2ca232190f883acebb")
	if err != nil {
		t.Error(err)
	}
}

func getConfig() *Config {
	return &Config{
		Database:   "cronsun",
		Hosts:      []string{"192.168.1.81:27018"},
		Timeout:    60,
		UserName:   "root",
		Password:   "123456",
		AuthSource: "admin",
	}
}
