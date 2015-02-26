package state

import (
	"fmt"
	"github.com/silenteh/gantryos/utils"
	"testing"
)

func TestInitSlaveDB(t *testing.T) {

	key := "TEST"
	value := "OK"
	bucket := "test"
	dbt, err := InitSlaveDB("./test_db.db")

	if err != nil {
		t.Fatal(err)
	}

	defer dbt.Close()

	for i := 0; i < 100; i++ {
		dbt.Set(bucket, key, value)
	}

	// GET
	if data, err := dbt.Get(bucket, key); err != nil {
		t.Fatal(err)
	} else {
		if data != value {
			t.Fatal("Retrived key from DB mismatch!")
		}
	}

	// EXISTS
	if !dbt.Exists(bucket, key) {
		t.Fatal("Key do not exist, but it supposed to !")
	}

	// GET ALL
	values := dbt.GetAllTasks(bucket)
	if len(values) != 1 {
		t.Fatal("Retrieving all values failed !")
	}

	if values[0] != value {
		t.Fatal("Retrieving all values failed: the retrived key mismatches !")
	}

	// GET ALL KEY VALUES
	keyValues := dbt.GetAllKeyValues(bucket)
	if len(keyValues) != 1 {
		t.Fatal("Retrieving all key values failed !")
	}

	if *keyValues[key] != value {
		t.Fatal("Retrieving all key values failed: the retrived key mismatches !")
	}

	// DELETE
	if err := dbt.Delete(bucket, key); err != nil {
		t.Fatal(err)
	}

	// EXISTS
	if dbt.Exists(bucket, key) {
		t.Fatal("Key exists, but we deleted it previously !")
	}

	if data, err := dbt.Get(bucket, key); err != nil {
		t.Fatal(err)
	} else {
		if data != "" {
			t.Fatal("There should be no keys !")
		}
	}

	// Backup
	if filename, err := dbt.Backup(); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(filename)
		utils.RemoveDir("./" + filename)
		utils.RemoveDir("./test_db.db")
	}

	fmt.Println("- State: OK")

}
