package common

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"
)

func JsonStringify(obj interface{}) (v string) {
	b, err := json.Marshal(obj)
	if err != nil {
		log.Fatal("Stringify error ", err)
		return
	}
	v = string(b)
	return
}

func JsonParse(obj interface{}, s string) (err error) {
	err = json.Unmarshal([]byte(s), obj)
	if err != nil {
		log.Fatal("ParseJson:", s)
		log.Fatal(err)
	}
	return
}

func IsExist( f string) bool {
	_, err := os.Stat(f);
	return err == nil || os.IsExist(err)
}

func GetMinutesDiffer (start, end string) int {
	var minutes int
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix() //
		minutes = int(diff / 60)
		return minutes
	} else {
		return minutes
	}
}

func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

