package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var processing = false

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	dbConnect()
	defer dbClose()
	initialise()

	http.HandleFunc("GET /", getHome)
	http.HandleFunc("GET /ip/{ip}", getIp)
	http.HandleFunc("GET /random/{ipVersion}", getRandomIp)
	http.HandleFunc("GET /benchmark/{ipVersion}/{times}", getBenchmark)

	fmt.Printf("starting server on %s:%s\n", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")), nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func initialise() {
	loadDbStructure()

	initialised, missing := loadCheckInitialised()

	if !initialised {
		fmt.Println("initialising data source(s)...")
		go upgrade(missing)
	}

	if len(os.Getenv("UPDATE_TIME")) > 0 {
		normaliser	:= time.NewTicker(time.Second)
		checker		:= time.NewTicker(2 * time.Minute)
		quit		:= make(chan struct{})

		go func() {
			for {
				select {
					case <- normaliser.C:
						checker.Reset(durationUntil(os.Getenv("UPDATE_TIME")))
						normaliser.Stop()
					case <- checker.C:
						update(checker)
					case <- quit:
						checker.Stop()
						return
				}
			}
		}()
	}
}

func upgrade(missing []string) {
	processing = true;

	dataToLoad := downloadDataToLoad(missing)
	loadData(dataToLoad)

	processing = false;
}

func update(checker *time.Ticker) {
	if !processing {
		fmt.Println("checking for updates...")
		go upgrade([]string{})

		checker.Reset(24 * time.Hour)
	}
}