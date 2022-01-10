package main

import (
	"fmt"
	"sync"

	"github.com/patrickmn/go-cache"
)

type Cache struct {
	ordersDataJson string
}

func restartCache(cacheProgram *cache.Cache) *cache.Cache {

	fmt.Println("restart cache")
	cacheProgram = cache.New(cache.DefaultExpiration, cache.DefaultExpiration)

	err := cacheProgram.LoadFile("cache") //не добавлена процедура синхронизации кэша с бд

	if err != nil {
		fmt.Println("new cache")
		cacheProgram = OrdersDBCache(cacheProgram)
	}

	return cacheProgram

}
func addDataCache(cacheProgram *cache.Cache, orderId string, orderDataJson []byte, mutex *sync.Mutex) *cache.Cache {

	mutex.Lock()
	stringOrderDataJson := string(orderDataJson)
	err := cacheProgram.Add(orderId, stringOrderDataJson, cache.DefaultExpiration)
	if err != nil {
		println(err)
	}

	mutex.Unlock()

	return cacheProgram

}
