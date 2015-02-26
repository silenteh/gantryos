package state

import (
	"github.com/boltdb/bolt"
	log "github.com/golang/glog"
	"os"
	"time"
)

var bucketName = []byte("slave")

type stateDB struct {
	db     *bolt.DB
	dbName string
}

func InitSlaveDB(dbName string) (stateDB, error) {
	var err error
	var state stateDB

	// remember the DB NAME
	state.dbName = dbName

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	state.db, err = bolt.Open(dbName, 0600, &bolt.Options{Timeout: 30 * time.Second})
	if err != nil {
		return state, err
	}

	// init the slave bucket
	err = state.db.Update(func(tx *bolt.Tx) error {
		if tx.Bucket(bucketName) == nil {
			_, err := tx.CreateBucket(bucketName)
			if err != nil {
				log.Fatalf("create bucket: %s - %s", bucketName, err)
				return err
			}
			return nil
		}
		return nil
	})

	return state, err
}

func (s stateDB) Close() {
	s.db.Close()
}

func (s stateDB) Set(key, value string) error {

	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketName).Put([]byte(key), []byte(value))
	})
}

func (s stateDB) Get(key string) (string, error) {

	var value string
	var err error

	s.db.View(func(tx *bolt.Tx) error {
		value = string(tx.Bucket(bucketName).Get([]byte(key)))
		return nil
	})

	return value, err
}

func (s stateDB) GetAllTasks() []string {

	var tasks []string
	total := 0

	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			//fmt.Printf("key=%s, value=%s\n", k, v)
			vString := string(v)
			tasks = append(tasks, []string{vString}...)
			total++
		}

		return nil
	})

	return tasks

}

func (s stateDB) Exists(key string) bool {
	if data, err := s.Get(key); err != nil {
		return false
	} else {
		return data != ""
	}
}

func (s stateDB) Delete(key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketName).Delete([]byte(key))
	})
}

// backs up the DB and returns the file name of the backup (because we add a time to it)
func (s stateDB) Backup() (string, error) {

	now := time.Now()
	nowFormat := now.Format("2006_01_02__15_04_05")
	backupFileName := s.dbName + "_" + nowFormat + ".db"

	f, err := os.Create(backupFileName)
	if err != nil {
		return backupFileName, err
	}

	err = s.db.View(func(tx *bolt.Tx) error {
		return tx.Copy(f)
	})

	return backupFileName, err
}
