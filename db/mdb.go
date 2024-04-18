package db

var (
	mgoDB *Mdb
)

func GetDb() *Mdb {
	return mgoDB
}

func SetDb(db *Mdb) {
	mgoDB = db
}
