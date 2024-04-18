package cronsun

import (
	"cronsun/db/entries"
	"cronsun/log"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"cronsun/conf"
	client "github.com/coreos/etcd/clientv3"
)

// 执行 cron cmd 的进程
// 注册到 /cronsun/node/<id>
type Node struct {
	Data *entries.Node
}

func (n *Node) String() string {
	return "node[" + n.Data.ID + "] pid[" + n.Data.PID + "]"
}

func (n *Node) Put(opts ...client.OpOption) (*client.PutResponse, error) {
	return DefalutClient.Put(conf.Config.Node+n.Data.ID, n.Data.PID, opts...)
}

func (n *Node) Del() (*client.DeleteResponse, error) {
	return DefalutClient.Delete(conf.Config.Node + n.Data.ID)
}

// 判断 node 是否已注册到 etcd
// 存在则返回进行 pid，不存在返回 -1
func (n *Node) Exist() (pid int, err error) {
	resp, err := DefalutClient.Get(conf.Config.Node + n.Data.ID)
	if err != nil {
		return
	}

	if len(resp.Kvs) == 0 {
		return -1, nil
	}

	if pid, err = strconv.Atoi(string(resp.Kvs[0].Value)); err != nil {
		if _, err = DefalutClient.Delete(conf.Config.Node + n.Data.ID); err != nil {
			return
		}
		return -1, nil
	}

	p, err := os.FindProcess(pid)
	if err != nil {
		return -1, nil
	}

	// TODO: 暂时不考虑 linux/unix 以外的系统
	if p != nil && p.Signal(syscall.Signal(0)) == nil {
		return
	}

	return -1, nil
}

func GetNodeGroups() (list []*Group, err error) {
	resp, err := DefalutClient.Get(conf.Config.Group, client.WithPrefix(), client.WithSort(client.SortByKey, client.SortAscend))
	if err != nil {
		return
	}

	list = make([]*Group, 0, resp.Count)
	for i := range resp.Kvs {
		g := Group{}
		err = json.Unmarshal(resp.Kvs[i].Value, &g)
		if err != nil {
			err = fmt.Errorf("node.GetGroups(key: %s) error: %s", string(resp.Kvs[i].Key), err.Error())
			return
		}
		list = append(list, &g)
	}

	return
}

func WatchNode() client.WatchChan {
	return DefalutClient.Watch(conf.Config.Node, client.WithPrefix())
}

// On 结点实例启动后，在 mongoDB 中记录存活信息
func (n *Node) On() {
	n.Data.Alived, n.Data.Version, n.Data.UpTime = true, Version, time.Now()
	n.SyncToMgo()
}

// On 结点实例停用后，在 mongoDB 中去掉存活信息
func (n *Node) Down() {
	n.Data.Alived, n.Data.DownTime = false, time.Now()
	n.SyncToMgo()
}

func (n *Node) SyncToMgo() {
	if err := entries.SyncNodeToMgo(n.Data); err != nil {
		log.Errorf(err.Error())
	}
}

// RmOldInfo remove old version(< 0.3.0) node info
func (n *Node) RmOldInfo() {
	entries.RemoveNodeById(n.Data.IP)
	DefalutClient.Delete(conf.Config.Node + n.Data.IP)
}
