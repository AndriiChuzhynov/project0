package main

import (
	"fmt"
	"github.com/labstack/echo"
	"sync"
	"time"
)

type mockDatabase struct {
	variable int
	sync.Mutex
}

type counter struct {
	variable int
	sync.Mutex
}

var db mockDatabase
var globalCounter counter

func (db *mockDatabase) changeCounter(val int) {

	db.Lock()
	defer db.Unlock()
	time.Sleep(time.Second * 5)
	db.variable = db.variable + val
	fmt.Println(db.variable)

}

func globalIncrement(e echo.Context) error {

	globalCounter.Lock()
	globalCounter.variable = globalCounter.variable + 1
	globalCounter.Unlock()

	return nil

}

func globalDecrement(e echo.Context) error {

	globalCounter.Lock()
	globalCounter.variable = globalCounter.variable - 1
	globalCounter.Unlock()

	return nil

}

func (c *counter) sendToDB() {

	for {
		if c.variable == 0 {
			// Nothing to send
			time.Sleep(1 * time.Second)
			continue
		}

		c.Lock()
		valToSend := c.variable
		c.variable = 0
		c.Unlock()

		db.changeCounter(valToSend)
		time.Sleep(1 * time.Second)
	}

}

func main() {

	e := echo.New()
	go globalCounter.sendToDB()
	e.POST("/globalIncrement", globalIncrement)
	e.POST("/globalDecrement", globalDecrement)

	e.Logger.Fatal(e.Start(":8080"))

}
