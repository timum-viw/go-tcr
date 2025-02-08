package main

import (
	"github.com/stianeikeland/go-rpio/v4"
	"os"
	"context"
	"log"
	"time"
	"sync"
	"net/http"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

var (
	pin = rpio.Pin(17)
	mu sync.Mutex
	timer *time.Timer
	wg sync.WaitGroup
)

func main() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	pin.Output()
	pin.Low()

	wg.Add(2)
	go runServer(context.Background(), &wg)
	go pollButtons(&wg)
	wg.Wait()
}

func pollButtons(wg *sync.WaitGroup) {
	defer wg.Done()

	buttonPin := rpio.Pin(26)
	buttonPin.Input()
	buttonPin.PullUp()
	buttonPin.Detect(rpio.FallEdge) // enable falling edge event detection
	defer buttonPin.Detect(rpio.NoEdge)

	log.Println("press a button")

	for {
		if pin.EdgeDetected() {
			log.Println("button pressed")
		}
		time.Sleep(time.Second / 2)
	}
}

func runServer(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()
	listener, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return err
	}

	log.Println("App URL is ", listener.URL())
	return http.Serve(listener, http.HandlerFunc(requestHandler))
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
