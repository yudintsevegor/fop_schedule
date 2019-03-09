package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)


func dbExplorer(db *sql.DB, group string) [][]Subject {
//func dbExplorer(db *sql.DB) [][]Subject {
	var tablesNames = make([]string, 0, 1)
	var tableName string

	//For debugging
	rowsTb, err := db.Query("SHOW TABLES")
	for rowsTb.Next() {
		err = rowsTb.Scan(&tableName)
		if err != nil {
			log.Fatal(err)
		}
		tablesNames = append(tablesNames, tableName)
	}
	rowsTb.Close()
	for _, key := range tablesNames {
		fmt.Println(key)
	}

	var allWeek = make([][]Subject, 0, 6)
//	for _, key := range tablesNames {
//		var allWeek = make([][]Subject, 0, 6)
		req := fmt.Sprintf("SELECT first, second, third, fourth, fifth FROM `%v`", group)
		rows, err := db.Query(req)
		for rows.Next() {
			var rawLes = make([]string, 5)
			err = rows.Scan(&rawLes[0], &rawLes[1], &rawLes[2], &rawLes[3], &rawLes[4])
			if err != nil {
				log.Fatal(err)
			}
			les := parsePercent(rawLes)
			allWeek = append(allWeek, les)
		}
//		fmt.Println("==================================" + group + "==============================")
//		for i, v := range allWeek {
//			fmt.Println("===========", i+1, "========")
//			for _, val := range v {
//				fmt.Println(val)
//			}
//		}
//	}

	return allWeek
}

//func main() {
//	db, err := sql.Open("mysql", DSN)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = db.Ping()
//	if err != nil {
//		log.Fatal(err)
//	}
//	group := "401"
//	allWeek := dbExplorer(db, group)
//	
//	fmt.Println("==================================" + group + "==============================")
//	for i, v := range allWeek {
//		fmt.Println("===========", i+1, "========")
//		for _, val := range v {
//			fmt.Println(val)
//		}
//	}
//}
//
var re = regexp.MustCompile("(.*)%(.*)%(.*)")

func parsePercent(arr []string) []Subject {
	nlr := Subject{}
	var result = make([]Subject, 0, 5)
	for _, val := range arr {
		res := re.FindStringSubmatch(val)
		nlr.Name = res[1]
		nlr.Lector = res[2]
		nlr.Room = res[3]
		result = append(result, nlr)
	}

	return result
}
