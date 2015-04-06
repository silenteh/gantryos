package state

import (
	"github.com/boltdb/bolt"
	log "github.com/golang/glog"
	"os"
	"time"
)

//var bucketName = []byte("slave")

type StateDB struct {
	db     *bolt.DB
	dbName string
}

func InitSlaveDB(dbName string) (StateDB, error) {
	var err error
	var state StateDB

	// remember the DB NAME
	state.dbName = dbName

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	state.db, err = bolt.Open(dbName, 0700, &bolt.Options{Timeout: 30 * time.Second})
	if err != nil {
		return state, err
	}

	return state, err
}

func (s StateDB) initBucket(bucket string) error {
	// init the slave bucket
	err := s.db.Update(func(tx *bolt.Tx) error {
		bucketName := []byte(bucket)
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
	return err
}

func (s StateDB) Close() {
	s.db.Close()
}

func (s StateDB) Set(bucket, key, value string) error {

	if err := s.initBucket(bucket); err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucket)).Put([]byte(key), []byte(value))
	})
}

func (s StateDB) Get(bucket, key string) (string, error) {

	var value string
	var err error

	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			value = string(b.Get([]byte(key)))
		}
		return nil
	})

	return value, err
}

func (s StateDB) GetAllTasks(bucket string) []string {

	var tasks []string
	total := 0

	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				//fmt.Printf("key=%s, value=%s\n", k, v)
				vString := string(v)
				tasks = append(tasks, []string{vString}...)
				total++
			}
		}
		return nil
	})

	return tasks

}

func (s StateDB) GetAllKeyValues(bucket string) map[string]*string {

	tasks := make(map[string]*string)

	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				//fmt.Printf("key=%s, value=%s\n", k, v)
				sKey := string(k)
				sValue := string(v)
				tasks[sKey] = &sValue
			}
		}
		return nil
	})

	return tasks

}

func (s StateDB) Exists(bucket, key string) bool {

	if data, err := s.Get(bucket, key); err != nil {
		return false
	} else {
		return data != ""
	}
}

func (s StateDB) Delete(bucket, key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			return b.Delete([]byte(key))
		}
		return nil
	})
}

// backs up the DB and returns the file name of the backup (because we add a time to it)
func (s StateDB) Backup() (string, error) {

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
