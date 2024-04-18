package entries

import (
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func TestGetJobLogById(t *testing.T) {
	InitTestDb()
	l, err := GetJobLogById("664c47d8921911193cea155e")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(l)
}

func TestCreateJobLog(t *testing.T) {

	InitTestDb()

	jl := JobLog{
		JobId: "j.ID",

		JobGroup: "j.Group",
		Name:     "j.Name",
		User:     "j.User",

		Node:     "j.runOn",
		Hostname: "j.hostname",
		IP:       "j.ip",

		Command: "j.Command",
		Output:  "rs",
		Success: true,

		BeginTime: time.Now(),
		EndTime:   time.Now().Add(time.Minute),
		Cleanup:   time.Now().Add(time.Hour * 24),
	}
	CreateJobLog(jl, nil)
}

func TestGetJobLatestLogList(t *testing.T) {
	InitTestDb()
	query := bson.M{}
	List, Total, err := GetJobLogList(query, 1, 10, bson.D{{Key: "beginTime", Value: -1}})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(List, Total)

}

func TestGetJobLatestLogListByJobIds(t *testing.T) {
	InitTestDb()
	m, err := GetJobLatestLogListByJobIds([]string{"j.ID"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(m)

}

func TestJobLogDailyStat(t *testing.T) {
	InitTestDb()
	ls, err := JobLogDailyStat(time.Now().Add(-time.Hour*24*7), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ls)
}
