package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/patrickmn/go-cache"
)

var ConnectStrDataBase string = "user=adminl0 password=123 dbname=taskl0 sslmode=disable"

func addOrderDataDB(idData string, data []byte) {
	fmt.Println("add order data in data base")

	db, _ := sql.Open("postgres", ConnectStrDataBase)
	defer db.Close()

	var query string = "CALL public.\"addOrder\"('" + idData + "','" + string(data) + "');"

	db.Exec(query)

}

func getOrderDataDB(idData string) []byte { //

	db, err := sql.Open("postgres", ConnectStrDataBase)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var query string = "SELECT public.\"getDateId\"('" + idData + "');"

	var result []byte

	db.QueryRow(query).Scan(&result)

	return result

}

func addErrData(errDataString string) {
	fmt.Println("add order err in data base")

	db, err := sql.Open("postgres", ConnectStrDataBase)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var query string = "CALL public.\"addErrData\"('" + errDataString + "');"

	db.Exec(query)
}

func OrdersDBCache(cacheProgram *cache.Cache) *cache.Cache {

	fmt.Println("restart cache = db")
	db, err := sql.Open("postgres", ConnectStrDataBase)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var query string = "SELECT public.\"getOrders\"();"

	row, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer row.Close()

	for row.Next() {
		var jsonOrderDB []byte
		var order map[string]interface{}

		row.Scan(&jsonOrderDB)

		json.Unmarshal(jsonOrderDB, &order)

		cacheProgram.Add(fmt.Sprintf("%+v", order["order_uid"]), string(jsonOrderDB), cache.DefaultExpiration)

	}

	cacheProgram.SaveFile("cache")

	return cacheProgram

}
