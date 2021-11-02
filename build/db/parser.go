package db

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

type XOQuery struct {
	TypeName string
	Query    string
}

func init() {
	os.Setenv("MAGEFILE_VERBOSE", "true")
	// os.Setenv("CGO_ENABLED", "0")
}

func FindXOQueries(databaseSvc string) []XOQuery {
	files, err := ioutil.ReadDir("./db/" + databaseSvc + "/")
	if err != nil {
		log.Fatalln(err)
	}

	rxp, err := regexp.Compile(`xo-([a-zA-Z]{1,25})\.sql`)

	if err != nil {
		log.Fatalln("Hey you, check the regex in line 21 of build/db/parser")
	}

	queries := make([]XOQuery, 0)

	for _, v := range files {
		if match := rxp.FindSubmatch([]byte(v.Name())); len(match) > 0 {
			fmt.Println()
			body, err := ioutil.ReadFile("./db/" + databaseSvc + "/" + v.Name())

			if err != nil {
				log.Fatalln("Can't open a mfing file you dumbo")
			}
			xoQ := &XOQuery{
				TypeName: string(match[1]),
				Query:    string(body),
			}
			queries = append(queries, *xoQ)
		}
	}

	return queries
}
