package main

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/fiatjaf/eventstore/sqlite3"
	"github.com/fiatjaf/khatru"
	"github.com/fiatjaf/khatru/blossom"
	"github.com/kehiy/blobstore/disk"
)

var config Config

func main() {
	log.SetPrefix("Blossom-Server ")
	log.Printf("Running %s\n", StringVersion())

	relay := khatru.NewRelay()

	LoadConfig()

	db := &sqlite3.SQLite3Backend{DatabaseURL: path.Join(config.WorkingDirectory, "database")}
	if err := db.Init(); err != nil {
		log.Fatalf("can't init database: %v\n", err)
	}

	bs := disk.New(path.Join(config.WorkingDirectory, "blobs"))

	bl := blossom.New(relay, fmt.Sprintf("http://localhost%s", config.Port))
	bl.Store = blossom.EventStoreBlobIndexWrapper{Store: db, ServiceURL: bl.ServiceURL}

	bl.StoreBlob = append(bl.StoreBlob, bs.Store)
	bl.DeleteBlob = append(bl.DeleteBlob, bs.Delete)
	bl.LoadBlob = append(bl.LoadBlob, bs.Load)

	if err := http.ListenAndServe(config.Port, relay); err != nil {
		log.Fatalf("Can't start the blossom server: %v", err)
	}
}
