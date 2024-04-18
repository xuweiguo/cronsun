package web

import (
	"cronsun/db/entries"
	"errors"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"net/http"
	"strings"
	"time"
)

type JobLog struct{}

func (jl *JobLog) GetDetail(ctx *Context) {
	vars := mux.Vars(ctx.R)
	id := strings.TrimSpace(vars["id"])
	if len(id) == 0 {
		outJSONWithCode(ctx.W, http.StatusBadRequest, "empty log id.")
		return
	}

	//objectId, err := primitive.ObjectIDFromHex(id)
	//if err != nil {
	//	outJSONWithCode(ctx.W, http.StatusBadRequest, "invalid ObjectId.")
	//	return
	//}

	logDetail, err := entries.GetJobLogById(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, mongo.ErrNoDocuments) {
			statusCode = http.StatusNotFound
		}
		outJSONWithCode(ctx.W, statusCode, err.Error())
		return
	}

	outJSON(ctx.W, logDetail)
}

func searchText(field string, keywords []string) (q []bson.M) {
	for _, k := range keywords {
		k = strings.TrimSpace(k)
		if len(k) == 0 {
			continue
		}
		q = append(q, bson.M{field: bson.M{"$regex": k, "$options": "i"}})
	}
	return q
}

func (jl *JobLog) GetList(ctx *Context) {
	hostnames := getStringArrayFromQuery("hostnames", ",", ctx.R)
	ips := getStringArrayFromQuery("ips", ",", ctx.R)
	names := getStringArrayFromQuery("names", ",", ctx.R)
	ids := getStringArrayFromQuery("ids", ",", ctx.R)
	begin := getTime(ctx.R.FormValue("begin"))
	end := getTime(ctx.R.FormValue("end"))
	page := getPage(ctx.R.FormValue("page"))
	failedOnly := ctx.R.FormValue("failedOnly") == "true"
	pageSize := getPageSize(ctx.R.FormValue("pageSize"))
	orderBy := bson.D{{Key: "beginTime", Value: -1}}

	query := bson.M{}
	var textSearch = make([]bson.M, 0, 2)
	textSearch = append(textSearch, searchText("hostname", hostnames)...)
	textSearch = append(textSearch, searchText("name", names)...)

	if len(ips) > 0 {
		query["ip"] = bson.M{"$in": ips}
	}

	if len(ids) > 0 {
		query["jobId"] = bson.M{"$in": ids}
	}

	if !begin.IsZero() {
		query["beginTime"] = bson.M{"$gte": begin}
	}
	if !end.IsZero() {
		query["endTime"] = bson.M{"$lt": end.Add(time.Hour * 24)}
	}

	if failedOnly {
		query["success"] = false
	}

	if len(textSearch) > 0 {
		query["$or"] = textSearch
	}

	var pager struct {
		Total int               `json:"total"`
		List  []*entries.JobLog `json:"list"`
	}
	var err error
	if ctx.R.FormValue("latest") == "true" {
		var latestLogList []*entries.JobLatestLog
		latestLogList, pager.Total, err = entries.GetJobLatestLogList(query, page, pageSize, orderBy)
		for i := range latestLogList {
			latestLogList[i].JobLog.Id, _ = primitive.ObjectIDFromHex(latestLogList[i].RefLogId)
			pager.List = append(pager.List, &latestLogList[i].JobLog)
		}
	} else {
		pager.List, pager.Total, err = entries.GetJobLogList(query, page, pageSize, orderBy)
	}
	if err != nil {
		outJSONWithCode(ctx.W, http.StatusInternalServerError, err.Error())
		return
	}

	pager.Total = int(math.Ceil(float64(pager.Total) / float64(pageSize)))
	outJSON(ctx.W, pager)
}
