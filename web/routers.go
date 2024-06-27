package web

import (
	"cronsun/db/entries"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"cronsun"
)

func GetVersion(ctx *Context) {
	outJSON(ctx.W, cronsun.Version)
}

func initRouters() (s *http.Server, err error) {
	jobHandler := &Job{}
	nodeHandler := &Node{}
	jobLogHandler := &JobLog{}
	infoHandler := &Info{}
	configHandler := &Configuration{}
	authHandler := &Authentication{}
	adminHandler := &Administrator{}

	r := mux.NewRouter()
	subrouter := r.PathPrefix("/v1").Subrouter()
	subrouter.Handle("/version", NewBaseHandler(GetVersion)).Methods("GET")

	h := NewBaseHandler(authHandler.GetAuthSession)
	subrouter.Handle("/session", h).Methods("GET")
	h = NewBaseHandler(authHandler.DeleteAuthSession)
	subrouter.Handle("/session", h).Methods("DELETE")

	h = NewBaseHandler(authHandler.SetPassword)
	subrouter.Handle("/user/setpwd", h).Methods("POST")

	h = NewAdminAuthHandler(adminHandler.GetAccount)
	subrouter.Handle("/admin/account/{email}", h).Methods("GET")
	h = NewAdminAuthHandler(adminHandler.GetAccountList)
	subrouter.Handle("/admin/accounts", h).Methods("GET")
	h = NewAdminAuthHandler(adminHandler.AddAccount)
	subrouter.Handle("/admin/account", h).Methods("PUT")
	h = NewAdminAuthHandler(adminHandler.UpdateAccount)
	subrouter.Handle("/admin/account", h).Methods("POSt")

	// get job list
	h = NewAuthHandler(jobHandler.GetList, entries.Reporter)
	subrouter.Handle("/jobs", h).Methods("GET")
	// get a job group list
	h = NewAuthHandler(jobHandler.GetGroups, entries.Reporter)
	subrouter.Handle("/job/groups", h).Methods("GET")
	// create/update a job
	h = NewAuthHandler(jobHandler.UpdateJob, entries.Developer)
	subrouter.Handle("/job", h).Methods("PUT")
	// pause/start
	h = NewAuthHandler(jobHandler.ChangeJobStatus, entries.Developer)
	subrouter.Handle("/job/{group}-{id}", h).Methods("POST")
	// batch pause/start
	h = NewAuthHandler(jobHandler.BatchChangeJobStatus, entries.Developer)
	subrouter.Handle("/jobs/{op}", h).Methods("POST")
	// get a job
	h = NewAuthHandler(jobHandler.GetJob, entries.Reporter)
	subrouter.Handle("/job/{group}-{id}", h).Methods("GET")
	// remove a job
	h = NewAuthHandler(jobHandler.DeleteJob, entries.Developer)
	subrouter.Handle("/job/{group}-{id}", h).Methods("DELETE")

	h = NewAuthHandler(jobHandler.GetJobNodes, entries.Reporter)
	subrouter.Handle("/job/{group}-{id}/nodes", h).Methods("GET")

	h = NewAuthHandler(jobHandler.JobExecute, entries.Developer)
	subrouter.Handle("/job/{group}-{id}/execute", h).Methods("PUT")

	// query executing job
	h = NewAuthHandler(jobHandler.GetExecutingJob, entries.Reporter)
	subrouter.Handle("/job/executing", h).Methods("GET")

	// kill an executing job
	h = NewAuthHandler(jobHandler.KillExecutingJob, entries.Developer)
	subrouter.Handle("/job/executing", h).Methods("DELETE")

	// get job log list
	h = NewAuthHandler(jobLogHandler.GetList, entries.Reporter)
	subrouter.Handle("/logs", h).Methods("GET")
	// get job log
	h = NewAuthHandler(jobLogHandler.GetDetail, entries.Developer)
	subrouter.Handle("/log/{id}", h).Methods("GET")

	h = NewAuthHandler(nodeHandler.GetNodes, entries.Reporter)
	subrouter.Handle("/nodes", h).Methods("GET")
	h = NewAuthHandler(nodeHandler.DeleteNode, entries.Developer)
	subrouter.Handle("/node/{ip}", h).Methods("DELETE")
	// get node group list
	h = NewAuthHandler(nodeHandler.GetGroups, entries.Reporter)
	subrouter.Handle("/node/groups", h).Methods("GET")
	// get a node group by group id
	h = NewAuthHandler(nodeHandler.GetGroupByGroupId, entries.Reporter)
	subrouter.Handle("/node/group/{id}", h).Methods("GET")
	// create/update a node group
	h = NewAuthHandler(nodeHandler.UpdateGroup, entries.Developer)
	subrouter.Handle("/node/group", h).Methods("PUT")
	// delete a node group
	h = NewAuthHandler(nodeHandler.DeleteGroup, entries.Developer)
	subrouter.Handle("/node/group/{id}", h).Methods("DELETE")

	h = NewAuthHandler(infoHandler.Overview, entries.Reporter)
	subrouter.Handle("/info/overview", h).Methods("GET")

	h = NewAuthHandler(configHandler.Configuratios, entries.Reporter)
	subrouter.Handle("/configurations", h).Methods("GET")

	r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", disableDirectoryListing(http.FS(webUi))))
	r.NotFoundHandler = NewBaseHandler(notFoundHandler)

	s = &http.Server{
		Handler: r,
	}
	return s, nil
}

func staticFileHandler() http.HandlerFunc {
	fs := http.FS(webUi)
	return func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(fs).ServeHTTP(w, r)
	}
}

func disableDirectoryListing(fs http.FileSystem) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 获取请求的路径
		upath := r.URL.Path
		sFilePath := "/ui/dist/" + upath
		if upath == "" {
			sFilePath += "index.html" // 重定向根目录到index.html，根据实际情况调整
		}
		// 检查路径是否指向一个目录
		info, err := fs.Open(sFilePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer info.Close()

		stat, err := info.Stat()
		if err != nil {
			http.Error(w, "Not a file", http.StatusForbidden)
			return
		}

		if stat.IsDir() {
			http.NotFound(w, r) // 或者重定向到某个默认页面
			return
		}
		http.ServeContent(w, r, stat.Name(), stat.ModTime(), info)
	})
}

func notFoundHandler(c *Context) {
	_notFoundHandler(c.W, c.R)
}

func _notFoundHandler(w http.ResponseWriter, r *http.Request) {
	const html = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>404 page not found</title>
</head>
<body>
    The page you are looking for is not found. Redirect to <a href="/ui/">Dashboard</a> after <span id="s">5</span> seconds.
</body>
<script type="text/javascript">
var s = 5;
setInterval(function(){
    s--;
    document.getElementById('s').innerText = s;
    if (s === 0) location.href = '/ui/';
}, 1000);
</script>
</html>`

	if strings.HasPrefix(strings.TrimLeft(r.URL.Path, "/"), "v1") {
		outJSONWithCode(w, http.StatusNotFound, "Api not found.")
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(html))
	}
}
