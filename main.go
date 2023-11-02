/*
Views (c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

main.go

Main Views generator file.
*/

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	hostname    = flag.String("host", "", "MySQL server hostname")
	port        = flag.Int("port", 3306, "MySQL port")
	username    = flag.String("username", "", "MySQL server username")
	password    = flag.String("password", "", "MySQL user password")
	database    = flag.String("database", "", "Database name")
	packageName = flag.String("package", "main", "package name")
	output      = flag.String("output", "", "output filename (default <snake_case_database_name>_views.go)")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "views: generate Go data structures and iterators for MySQL views.\nAvailable options:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("views: ")
	flag.Usage = Usage

	flag.Parse()
	if len(*hostname)+len(*username)+len(*database) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", *username, *password, *hostname, *port, *database)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("sql.Open: %v", err)
	}
	defer db.Close()
	data := NewTemplateData(*packageName)
	if err := data.processDatabase(db); err != nil {
		log.Fatal(err)
	}
	outputFileName := *output
	if outputFileName == "" {
		outputFileName = fmt.Sprintf("%s_views.go", ToSnakeCase(*database))
	}
	if err := data.Save(outputFileName); err != nil {
		log.Fatal(err)
	}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
