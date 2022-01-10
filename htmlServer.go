package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/patrickmn/go-cache"
)

type ViewData struct {
	OrderId   string
	OrderInfo string
}

func serverHtmlStart(cacheProgram *cache.Cache) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Start html /")

		http.ServeFile(w, r, "tmpl/index.html")
	})

	http.HandleFunc("/postform", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Start html /postform")

		orderId := r.FormValue("orderId")

		iOrderJson, flag := cacheProgram.Get(orderId)

		var data ViewData

		if flag {

			var strOrderJson string = fmt.Sprintf("%+v", iOrderJson)

			var orderMapJson map[string]interface{}
			json.Unmarshal([]byte(strOrderJson), &orderMapJson)

			var orderId string = fmt.Sprintf("%+v", orderMapJson["order_uid"])

			data = ViewData{
				OrderId:   orderId,
				OrderInfo: strOrderJson,
			}

		} else {
			//Стоит добавить проверку есть ли в бд данные или нет,,т.к. нет синхронизации данных с бд и кэшем
			data = ViewData{
				OrderId:   "нет информации",
				OrderInfo: "нет информации",
			}
		}

		tmpl, _ := template.ParseFiles("tmpl/orderID.html")

		tmpl.Execute(w, data)

	})

	http.ListenAndServe(":8181", nil)
}
