package main

import (
	"github.com/stianeikeland/go-rpio/v4"
	"log/slog"
	"log"
	"time"
	"sync"
	"net/http"
	"timum-viw/go-tcr/tcr"
	"timum-viw/go-tcr/hash"
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
		slog.Error(err.Error())
		return
	}
	defer tcr.Close()
	slog.Info("polling tcr")
	for {
		time.Sleep(time.Second/2)
		str, err := tcr.Read()
		if err != nil {
			slog.Warn(err.Error())
			continue
		} else if str == "" {
			continue
		}
		id, err := hash.Validate(str)
		if err != nil {
			slog.Warn(err.Error())
		}
		if id > 0 {
			switchOn()
			slog.Info("qr code scanned", "id", id)
		}
	}
}

func pollButtons(wg *sync.WaitGroup) {
	defer wg.Done()

	buttonPin := rpio.Pin(6)
	buttonPin.Input()
	buttonPin.PullUp()
	buttonPin.Detect(rpio.FallEdge) 
	defer buttonPin.Detect(rpio.NoEdge)

	slog.Info("polling buttons")

	for {
		slog.Info("polling button")
		if buttonPin.EdgeDetected() {
			slog.Info("Button pressed")
			go switchOn()
		}
		time.Sleep(time.Second / 2)
	}
}

func runServer(wg *sync.WaitGroup) {
	defer wg.Done()
	slog.Info("server listening on :8080")
	err := http.ListenAndServe(":8080", http.HandlerFunc(requestHandler))
	if err != nil {
		slog.Error(err.Error())
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	switchOn()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request received"))
}

func switchOn() {
	mu.Lock()
	defer mu.Unlock()
	pin.High()
	if timer != nil {
		timer.Stop()
	}

	timer = time.AfterFunc(5*time.Second, func() {
		pin.Low()
	})
}


