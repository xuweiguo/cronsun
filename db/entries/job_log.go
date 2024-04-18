package entries

import (
	"context"
	"cronsun/db"
	"cronsun/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	Coll_JobLog       = "cronsun_job_log"
	Coll_JobLatestLog = "cronsun_job_latest_log"
	Coll_Stat         = "cronsun_stat"
)

// 任务执行记录
type JobLog struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	JobId     string             `bson:"jobId" json:"jobId"`               // 任务 Id，索引
	JobGroup  string             `bson:"jobGroup" json:"jobGroup"`         // 任务分组，配合 Id 跳转用
	User      string             `bson:"user" json:"user"`                 // 执行此次任务的用户
	Name      string             `bson:"name" json:"name"`                 // 任务名称
	Node      string             `bson:"node" json:"node"`                 // 运行此次任务的节点 id，索引
	Hostname  string             `bson:"hostname" json:"hostname"`         // 运行此次任务的节点主机名称，索引
	IP        string             `bson:"ip" json:"ip"`                     // 运行此次任务的节点主机IP，索引
	Command   string             `bson:"command" json:"command,omitempty"` // 执行的命令，包括参数
	Output    string             `bson:"output" json:"output,omitempty"`   // 任务输出的所有内容
	Success   bool               `bson:"success" json:"success"`           // 是否执行成功
	BeginTime time.Time          `bson:"beginTime" json:"beginTime"`       // 任务开始执行时间，精确到毫秒，索引
	EndTime   time.Time          `bson:"endTime" json:"endTime"`           // 任务执行完毕时间，精确到毫秒
	Cleanup   time.Time          `bson:"cleanup,omitempty" json:"-"`       // 日志清除时间标志
}

type JobLatestLog struct {
	JobLog   `bson:",inline"`
	RefLogId string `bson:"refLogId,omitempty" json:"refLogId"`
}

type StatExecuted struct {
	Total     int64  `bson:"total" json:"total"`
	Successed int64  `bson:"successed" json:"successed"`
	Failed    int64  `bson:"failed" json:"failed"`
	Date      string `bson:"date" json:"date"`
}

func GetJobLogById(id string) (l *JobLog, err error) {
	err = db.GetDb().FindId(Coll_JobLog, id, &l)
	return
}

func CreateJobLog(jl JobLog, logger log.Logger) {
	if jl.Id == primitive.NilObjectID {
		jl.Id = primitive.NewObjectID()
	}
	if err := db.GetDb().Insert(Coll_JobLog, jl); err != nil {
		if logger != nil {
			logger.Errorf(err.Error())
		}
	}

	latestLog := &JobLatestLog{
		RefLogId: jl.Id.Hex(),
		JobLog:   jl,
	}
	latestLog.Id = primitive.NilObjectID
	if err := db.GetDb().Upsert(Coll_JobLatestLog, bson.M{"node": jl.Node, "hostname": jl.Hostname, "ip": jl.IP, "jobId": jl.JobId, "jobGroup": jl.JobGroup}, bson.M{"$set": latestLog}); err != nil {
		if logger != nil {
			logger.Errorf(err.Error())
		}
	}

	var inc = bson.M{"total": 1}
	if jl.Success {
		inc["successed"] = 1
	} else {
		inc["failed"] = 1
	}

	err := db.GetDb().Upsert(Coll_Stat, bson.M{"name": "job-day", "date": time.Now().Format("2006-01-02")}, bson.M{"$inc": inc})
	if err != nil {
		if logger != nil {
			logger.Errorf("increase stat.job-day %s", err.Error())
		}
	}
	err = db.GetDb().Upsert(Coll_Stat, bson.M{"name": "job"}, bson.M{"$inc": inc})
	if err != nil {
		if logger != nil {
			logger.Errorf("increase stat.job %s", err.Error())
		}
	}

}

var selectForJobLogList = bson.M{"command": 0, "output": 0}

func GetJobLogList(query bson.M, page, size int, sort bson.D) (list []*JobLog, total int, err error) {
	err = db.GetDb().WithC(Coll_JobLog, func(c *mongo.Collection) error {
		totalTmp, err := c.CountDocuments(context.Background(), query)
		if err != nil {
			return err
		}
		total = int(totalTmp)
		findOptions := options.Find()
		findOptions.SetLimit(int64(size))
		findOptions.SetSkip(int64((page - 1) * size))
		findOptions.SetSort(sort)

		cursor, err := c.Find(context.Background(), query, findOptions)
		if err != nil {
			return err
		}
		defer cursor.Close(context.Background())

		err = cursor.All(context.Background(), &list)
		return err
	})
	return
}

func GetJobLatestLogList(query bson.M, page, size int, sort bson.D) (list []*JobLatestLog, total int, err error) {
	err = db.GetDb().WithC(Coll_JobLatestLog, func(c *mongo.Collection) error {
		totalTmp, err := c.CountDocuments(context.Background(), query)
		if err != nil {
			return err
		}
		total = int(totalTmp)
		findOptions := options.Find()
		findOptions.SetLimit(int64(size))
		findOptions.SetSkip(int64((page - 1) * size))
		findOptions.SetSort(sort)
		findOptions.SetProjection(selectForJobLogList)

		cursor, err := c.Find(context.Background(), query, findOptions)
		if err != nil {
			return err
		}
		defer cursor.Close(context.Background())

		err = cursor.All(context.Background(), &list)
		return err
	})
	return
}

func GetJobLatestLogListByJobIds(jobIds []string) (m map[string]*JobLatestLog, err error) {
	var list []*JobLatestLog

	err = db.GetDb().WithC(Coll_JobLatestLog, func(c *mongo.Collection) error {
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{Key: "beginTime", Value: 1}})
		findOptions.SetProjection(selectForJobLogList)
		if len(jobIds) == 0 {
			cursor, err := c.Find(context.Background(), bson.M{}, findOptions)
			if err != nil {
				return err
			}
			defer cursor.Close(context.Background())
			err = cursor.All(context.Background(), &list)
		} else {
			cursor, err := c.Find(context.Background(), bson.M{"jobId": bson.M{"$in": jobIds}}, findOptions)
			if err != nil {
				return err
			}
			defer cursor.Close(context.Background())
			err = cursor.All(context.Background(), &list)
		}
		return err
	})
	if err != nil {
		return
	}

	m = make(map[string]*JobLatestLog, len(list))
	for i := range list {
		m[list[i].JobId] = list[i]
	}
	return
}

func JobLogStat() (s *StatExecuted, err error) {
	err = db.GetDb().FindOne(Coll_Stat, bson.M{"name": "job"}, &s)
	return
}

func JobLogDailyStat(begin, end time.Time) (ls []*StatExecuted, err error) {
	const oneDay = time.Hour * 24
	err = db.GetDb().WithC(Coll_Stat, func(c *mongo.Collection) error {
		dateList := make([]string, 0, 8)

		cur := begin
		for {
			dateList = append(dateList, cur.Format("2006-01-02"))
			cur = cur.Add(oneDay)
			if cur.After(end) {
				break
			}
		}
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{Key: "date", Value: 1}})
		cursor, err := c.Find(context.Background(), bson.M{"name": "job-day", "date": bson.M{"$in": dateList}}, findOptions)
		if err != nil {
			return err
		}
		defer cursor.Close(context.Background())

		err = cursor.All(context.Background(), &ls)
		return err
	})

	return
}

func ClearJobLogs(expiration time.Duration) error {
	err := db.GetDb().WithC(Coll_JobLog, func(c *mongo.Collection) error {
		_, err := c.DeleteMany(context.Background(), bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"cleanup": bson.M{"$exists": true}},
				{"cleanup": bson.M{"$lte": time.Now()}},
			}},
			{"$and": []bson.M{
				{"cleanup": bson.M{"$exists": false}},
				{"endTime": bson.M{"$lte": time.Now().Add(-expiration)}},
			}},
		}})
		return err
	})
	return err
}

func EnsureJobLogIndex() error {

	return db.GetDb().WithC(Coll_JobLog, func(c *mongo.Collection) error {
		// 获取集合的索引视图
		indexView := c.Indexes()

		// 创建索引
		_, err := indexView.CreateMany(context.Background(), []mongo.IndexModel{
			{
				Keys: bson.D{{"beginTime", 1}},
			},
			{
				Keys: bson.D{{"hostname", 1}},
			},
			{
				Keys: bson.D{{"ip", 1}},
			},
		})
		return err
	})
}
