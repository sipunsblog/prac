package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"satrtp/worker"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type Client struct {
	name   string
	events chan *DashBoard
}
type DashBoard struct {
	User uint
}

func main() {
	e := echo.New()
	e.GET("/", runWorker)
	e.Logger.Fatal(e.Start(":1323"))
}

func runBuff(c chan int, wg *sync.WaitGroup) {
	c <- 7
	defer wg.Done()
}

func runBuff2(c chan int, wg *sync.WaitGroup) {
	vv := <-c
	fmt.Println(vv)
	for i := 0; i < 3; i++ {
		vv += i
	}
	c <- vv
	defer wg.Done()
	close(c)
}

func helloWorld(c echo.Context) error {
	myc := make(chan int, 3)
	var wg sync.WaitGroup
	wg.Add(1)
	go runBuff(myc, &wg)
	wg.Add(1)
	go runBuff2(myc, &wg)
	wg.Wait()
	// for {
	// 	val, ok := <-myc
	// 	if ok {
	// 		fmt.Println(val)
	// 	} else {
	// 		fmt.Println("closed")
	// 		break
	// 	}
	// }
	fmt.Println("Length of the channel is: ", len(myc))
	for val := range myc {

		fmt.Println(val)

	}
	return echo.ErrRendererNotRegistered

}

func dashboardHandler(c echo.Context) error {
	r := c.Request()
	w := c.Response()
	client := &Client{name: r.RemoteAddr, events: make(chan *DashBoard, 10)}
	go updateDashboard(client)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	timeout := time.After(1 * time.Second)
	select {
	case ev := <-client.events:
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.Encode(ev)
		fmt.Fprintf(w, "data: %v\n\n", buf.String())
		fmt.Printf("data: %v\n", buf.String())
	case <-timeout:
		fmt.Fprintf(w, ": nothing to sent\n\n")
	}

	return echo.ErrBadGateway
}

func updateDashboard(client *Client) {
	for {
		db := &DashBoard{
			User: uint(rand.Uint32()),
		}
		client.events <- db
	}
}

func runWorker(c echo.Context) error {
	worker.CreateClient()
	worker.WorkerInIt()
	return nil
}
