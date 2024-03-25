package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ShashwatPS/FileWarden/db"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	location := "./"
	runWatcher(location)
}

func runWatcher(location string) {
	watcher, err := fsnotify.NewWatcher()
	currentTime := time.Now()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				var typeOperation string
				if event.Has(fsnotify.Write) {
					typeOperation = "WRITE"
				} else if event.Has(fsnotify.Remove) {
					typeOperation = "REMOVE"
				} else if event.Has(fsnotify.Create) {
					typeOperation = "CREATE"
				} else if event.Has(fsnotify.Rename) {
					typeOperation = "RENAME"
				}

				fmt.Println("Date only:", currentTime.Format("2006-01-02"))
				fmt.Println("Time only:", currentTime.Format("15:04:05"))
				fmt.Println("event:", event.Op)
				fmt.Println("event:", event.Name)
				saveToDataBase(currentTime, location, event.Name, typeOperation)

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(location) // Can add multiple paths to track multiple locations
	if err != nil {
		log.Fatal(err)
	}
	<-make(chan struct{})
}

func saveToDataBase(currentTime time.Time, path string, name string, typeop string) {
	if err := run(currentTime, path, name, typeop); err != nil {
		panic(err)
	}
}

func run(currentTime time.Time, path string, name string, typeop string) error {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return err
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	SetTime := currentTime.Format("15:04:05")
	SetDate := currentTime.Format("2006-01-02")

	createdPost, err := client.Post.CreateOne(
		db.Post.CreatedAt.Set(SetDate),
		db.Post.PublishTime.Set(SetTime),
		db.Post.FilePath.Set(path),
		db.Post.FileName.Set(name),
		db.Post.Type.Set(typeop),
	).Exec(ctx)
	if err != nil {
		return err
	}

	result, _ := json.MarshalIndent(createdPost, "", "  ")
	fmt.Printf("created post: %s\n", result)

	return nil
}
