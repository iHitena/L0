package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/patrickmn/go-cache"
)

var cacheProgram *cache.Cache

func main() {

	cacheProgram = restartCache(cacheProgram)

	go subChan()

	go serverHtmlStart(cacheProgram)

	fmt.Scanln()
}

func subChan() {
	fmt.Println("Connect channel")
	nc, err := nats.Connect("mytoken@localhost")
	if err != nil {
		println("not connetc channel")
		time.Sleep(5 * time.Second)
		subChan()
	}

	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	recvCh := make(chan map[string]interface{})

	ec.BindRecvChan("order", recvCh)

	order := <-recvCh

	orderJsonData, err := json.Marshal(order)
	if err != nil {
		fmt.Println(err)
		go addErrData(fmt.Sprintf("%+v", order))

	} else {

		idOrder := fmt.Sprintf("%+v", order["order_uid"])

		go addOrderDataDB(idOrder, orderJsonData)

		addCache := func() {
			var mutex sync.Mutex
			cacheProgram = addDataCache(cacheProgram, idOrder, orderJsonData, &mutex)
		}
		go addCache()
	}

	nc.Close()
	ec.Close()
	subChan()

}
