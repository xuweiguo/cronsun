package entries

import "testing"

func TestGetNodes(t *testing.T) {
	InitTestDb()
	nodes, err := GetNodes()
	if err != nil {
		t.Error(err)
	}
	t.Log(nodes)
}

func TestSyncToMgo(t *testing.T) {
	InitTestDb()
	node := &Node{
		ID:       "1",
		IP:       "127.0.0.1",
		Version:  "1.0.0",
		Hostname: "localhost",
	}
	err := SyncNodeToMgo(node)
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveNodeById(t *testing.T) {
	InitTestDb()
	err := RemoveNodeById("1")
	if err != nil {
		t.Error(err)
	}

}

func TestGetNodesByID(t *testing.T) {
	InitTestDb()
	node, err := GetNodesByID("1")
	if err != nil {
		t.Error(err)
	}
	t.Log(node)
}
