package main

import (
	"elasticbulk/elastic"
	"elasticbulk/settings"
	"log"
	"os"
)

func terminate(err error) {
	log.Println(err)
	os.Exit(1)
}

func init() {
	if err := settings.InitConfig(); err != nil {
		terminate(err)
	}
	if err := elastic.InitElastic(); err != nil {
		terminate(err)
	}

}

func main() {
	elastic.BulkInsert()
}