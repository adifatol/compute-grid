package response

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/op/go-logging"
)

type DbHandler struct {
	dbname string
	dbHnd  *bolt.DB
}

type ResultRow struct {
	id  string
	avg float64
	sum float64
}

func ConnectDb(dbname string) *DbHandler {
	p := new(DbHandler)
	p.dbname = dbname

	db, err := bolt.Open(dbname, 0600, nil)
	if err != nil {
		Log := logging.MustGetLogger("database")
		Log.Critical(err)
	}
	p.dbHnd = db

	return p
}

func Save(results []interface{}, db *DbHandler) {
	for k := range results {
		item := results[k].(map[string]interface{})
		row := &ResultRow{
			id:  item["id"].(string),
			avg: item["avg"].(float64),
			sum: item["sum"].(float64),
		}

		db.dbHnd.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("results"))
			if err != nil {
				return err
			}
			encoded, err := json.Marshal(row)
			if err != nil {
				return err
			}
			return b.Put([]byte(row.id), []byte(encoded))
		})

	}
}

func Disconnect(p *DbHandler) {
	p.dbHnd.Close()
}
