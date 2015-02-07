package utils

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func UTCTimeStamp() int32 {
	return int32(time.Now().UTC().Unix())
}

func ConcatenateString(src string, additionals ...string) string {
	var buffer bytes.Buffer
	buffer.WriteString(src)
	for _, element := range additionals {
		buffer.WriteString(element)
	}

	return buffer.String()
}

func ExecCommand(stripNewLines bool, cmd string, args ...string) string {

	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	var finalString string

	if stripNewLines {
		finalString = strings.Replace(string(out), "\n", "", -1)
	} else {
		finalString = string(out)
	}
	return finalString
}

// check if a file exists
func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		return true
	}
	return false
}

func ReadFile(fileName string) []byte {
	buf := bytes.NewBuffer(nil)

	f, err := os.Open(fileName)
	//defer f.Close()
	if err != nil {
		log.Fatal(err)
	} else {
		io.Copy(buf, f)
		f.Close()
	}
	//s := string(buf.Bytes())
	return buf.Bytes()
}

func ReadFileToString(fileName string) string {
	data := ReadFile(fileName)
	return string(data)
}

func WriteFile(fileName string, data []byte, permission os.FileMode) error {
	err := ioutil.WriteFile(fileName, data, permission)
	return err
}

func StringToFloat64(value string, onErrorOne bool) float64 {
	i, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatal(err)
		if onErrorOne {
			return float64(1)
		}
		return float64(0)
	}
	return i
}

func StringToInt(value string, onErrorOne bool) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(err)
		if onErrorOne {
			return 1
		}
		return 0
	}
	return i
}

func StringToUINT64(value string, onErrorOne bool) uint64 {
	i, err := strconv.ParseUint(value, 0, 64)
	if err != nil {
		log.Fatal(err)
		if onErrorOne {
			return 1
		}
		return 0
	}
	return i
}

func IntToString(i int64, base int) string {
	s := strconv.FormatInt(i, base)
	return s
}

func BuildURI(uris ...string) string {

	var buffer bytes.Buffer
	for _, rest := range uris {
		buffer.WriteString(rest)
		buffer.WriteString("/")
	}

	return buffer.String()

}

func BuildURL(URI, ip, version, protocol, port string) string {
	return ConcatenateString(protocol, "://", ip, ":", port, version, URI)
}

func Chomp(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}

func ParseOutputCommand(output string) [][]string {

	lines := strings.Split(output, "\n")
	totalLines := len(lines)

	data := make([][]string, totalLines)

	total := 0
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 {
			data[total] = fields
			total++
		}

	}

	return data[:total]

}

func ParseOutputCommandWithHeader(output string, totalHeaderLines int) [][]string {

	lines := strings.Split(output, "\n")
	totalLines := len(lines)

	data := make([][]string, totalLines)

	total := 0
	for index, line := range lines {
		if index < totalHeaderLines {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			data[total] = fields
			total++
		}

	}

	return data[:total]
}

func CommandOutputToMap(data [][]string, labelPosition, valuePosition int) (map[string]uint64, error) {

	mapData := make(map[string]uint64)
	totalEntries := len(data)

	if labelPosition > valuePosition {
		return mapData, errors.New("value position is smaller than label position...this is most probably wrong")
	}

	for _, line := range data {
		// the value position is > than labelPosition
		if valuePosition > totalEntries {
			return mapData, errors.New("Index out of range")
		}

		label := strings.Replace(line[labelPosition], ":", "", -1)

		mapData[label] = StringToUINT64(line[valuePosition], false)
	}

	return mapData, nil
}