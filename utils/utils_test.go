package utils

import (
	"fmt"
	"testing"
)

func TestChomp(t *testing.T) {

	thisIsAString := "the quick brown fox jumps over the lazy dog\r\n"
	expected := "the quick brown fox jumps over the lazy dog"
	chomped := Chomp(thisIsAString, "\r\n")

	if chomped != expected {
		t.Fatalf("%s\n", "Error on string chomp!")
	}

	chomped = Chomp(thisIsAString, "\n\r")
	if chomped != expected {
		t.Fatalf("%s\n", "Error on string chomp!")
	}

	fmt.Println("Chomp: SUCCESS")
}

func TestParseOutputCommand(t *testing.T) {

	example := `MemTotal:        2049944 kB
	MemFree:         1795692 kB
	Buffers:           47180 kB
	Cached:           140220 kB
	SwapCached:            0 kB
	Active:            79052 kB
	Inactive:         132184 kB`

	data := ParseOutputCommand(example)

	if len(data) != 7 {
		t.Fatal("Error parsing the command output")
	}

	fmt.Println("ParseOutputCommand: SUCCESS")

}

func TestParseOutputCommandWithHeader(t *testing.T) {

	example := `Label		value unit

MemTotal:        2049944 kB
MemFree:         1795692 kB
Buffers:           47180 kB
Cached:           140220 kB
SwapCached:            0 kB
Active:            79052 kB
Inactive:         132184 kB
	`

	data := ParseOutputCommandWithHeader(example, 2)

	if len(data) != 7 {
		t.Fatal("Error parsing the command output")
	}

	fmt.Println("ParseOutputCommandWithHeader: SUCCESS")

}

func TestCommandOutputToMap(t *testing.T) {
	example := `Label		value unit

MemTotal:        2049944 kB
MemFree:         1795692 kB
Buffers:           47180 kB
Cached:           140220 kB
SwapCached:            0 kB
Active:            79052 kB
Inactive:         132184 kB
	`

	data := ParseOutputCommandWithHeader(example, 2)

	if len(data) != 7 {
		t.Fatal("Error parsing the command output")
	}

	if dataMap, err := CommandOutputToMap(data, 0, 1); err != nil {
		t.Fatal(err)
	} else {
		if dataMap["MemTotal"] != 2049944 {
			t.Fatal("Error converting the command output")
		}

		if dataMap["MemFree"] != 1795692 {
			t.Fatal("Error converting the command output")
		}
	}

	fmt.Println("CommandOutputToMap: SUCCESS")

}

func TestCommandOutputToMapArray(t *testing.T) {
	example := `Filesystem                          Size   Used  Avail Capacity  iused     ifree %iused  Mounted on
/dev/disk1                         465Gi   74Gi  391Gi    16% 19338597 102499001   16%   /
devfs                              183Ki  183Ki    0Bi   100%      634         0  100%   /dev
map -hosts                           0Bi    0Bi    0Bi   100%        0         0  100%   /net
map auto_home                        0Bi    0Bi    0Bi   100%        0         0  100%   /home
localhost:/JUBRlh2axvSNcf66Md3uWc  465Gi  465Gi    0Bi   100%        0         0  100%   /Volumes/MobileBackups
	`

	data := ParseOutputCommandWithHeader(example, 1)

	if dataMap, err := CommandOutputToMapArray(data, 8); err != nil {
		t.Fatal(err)
	} else {
		if dataMap["/"][0] != "/dev/disk1" {
			t.Fatal("Error converting the command output")
		}
	}

	fmt.Println("CommandOutputToMap: SUCCESS")

}
