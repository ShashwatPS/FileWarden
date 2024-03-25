package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	currentTime := time.Now()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Println("Date only:", currentTime.Format("2006-01-02"))
				fmt.Println("Time only:", currentTime.Format("15:04:05"))
				fmt.Println("event:", event.Op)
				fmt.Println("event:", event.Name)
				//if event.Has(fsnotify.Write) {
				//	log.Println("modified file:", event.Name)
				//}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add("./")
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
