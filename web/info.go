package web

import (
	"cronsun/db/entries"
	"time"

	v3 "github.com/coreos/etcd/clientv3"

	"cronsun"
	"cronsun/conf"
)

type Info struct{}

func (inf *Info) Overview(ctx *Context) {
	var overview = struct {
		TotalJobs        int64                   `json:"totalJobs"`
		JobExecuted      *entries.StatExecuted   `json:"jobExecuted"`
		JobExecutedDaily []*entries.StatExecuted `json:"jobExecutedDaily"`
	}{}

	const day = 24 * time.Hour
	days := 7

	overview.JobExecuted, _ = entries.JobLogStat()
	end := time.Now()
	begin := end.Add(time.Duration(1-days) * day)
	statList, _ := entries.JobLogDailyStat(begin, end)
	list := make([]*entries.StatExecuted, days)
	cur := begin

	for i := 0; i < days; i++ {
		date := cur.Format("2006-01-02")
		var se *entries.StatExecuted

		for j := range statList {
			if statList[j].Date == date {
				se = statList[j]
				statList = statList[1:]
				break
			}
		}

		if se != nil {
			list[i] = se
		} else {
			list[i] = &entries.StatExecuted{Date: date}
		}

		cur = cur.Add(day)
	}

	overview.JobExecutedDaily = list
	gresp, err := cronsun.DefalutClient.Get(conf.Config.Cmd, v3.WithPrefix(), v3.WithCountOnly())
	if err == nil {
		overview.TotalJobs = gresp.Count
	}

	outJSON(ctx.W, overview)
}
