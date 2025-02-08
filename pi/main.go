package main

import (
	"github.com/stianeikeland/go-rpio/v4"
	"log"
	"time"
	"sync"
	"net/http"
	"timum-viw/supercoop/tcr"
	"timum-viw/supercoop/hash"
)

var (
	pin = rpio.Pin(17)
	mu sync.Mutex
	timer *time.Timer
	wg sync.WaitGroup
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}
	defer rpio.Close()

	pin.Output()
	pin.Low()

	wg.Add(3)
	go runServer(&wg)
	go pollButtons(&wg)
	go pollTcr(&wg)
	wg.Wait()
}

func pollTcr(wg *sync.WaitGroup) {
	defer wg.Done()
	tcr, err := tcr.Open(1)
	if err != nil {
		log.Fatal(err)
	}
	defer tcr.Close()
	for {
		str, err := tcr.Read()
		if err != nil {
			log.Println(err)
			continue
		}
		id, err := hash.Validate(str)
		if id > 0 {
			switchOn()
			log.Println(id)
		}
		time.Sleep(time.Second/2)
	}
}

func pollButtons(wg *sync.WaitGroup) {
	defer wg.Done()

	buttonPin := rpio.Pin(6)
	buttonPin.Input()
	buttonPin.PullUp()
	buttonPin.Detect(rpio.FallEdge) 
	defer buttonPin.Detect(rpio.NoEdge)

	log.Println("press a button")

	for {
		if buttonPin.EdgeDetected() {
			switchOn()
		}
		time.Sleep(time.Second / 2)
	}
}

func runServer(wg *sync.WaitGroup) error {
	defer wg.Done()
	return http.ListenAndServe(":8080", http.HandlerFunc(requestHandler))
}

func switchOn() {
	mu.Lock()
	defer mu.Unlock()
	pin.High()
	if timer != nil {
		timer.Stop()
	}

	// Set a new debounce timer
	timer = time.AfterFunc(5*time.Second, func() {
		pin.Low()
	})
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	switchOn()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request received"))
}
