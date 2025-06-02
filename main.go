package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/fiatjaf/eventstore/sqlite3"
	"github.com/fiatjaf/khatru"
	"github.com/fiatjaf/khatru/blossom"
	"github.com/kehiy/blobstore/disk"
	"github.com/nbd-wtf/go-nostr"
)

var (
	config Config
	db     *sqlite3.SQLite3Backend
)

func main() {
	log.SetPrefix("Blossom-Server ")
	log.Printf("Running %s\n", StringVersion())

	relay := khatru.NewRelay()

	LoadConfig()

	db = &sqlite3.SQLite3Backend{DatabaseURL: path.Join(config.WorkingDirectory, "database")}
	if err := db.Init(); err != nil {
		log.Fatalf("can't init database: %v\n", err)
	}

	for _, ddl := range ddls {
		_, err := db.DB.Exec(ddl)
		if err != nil {
			log.Fatalf("can't init database: %v\n", err)
		}
	}

	bs := disk.New(path.Join(config.WorkingDirectory, "blobs"))

	bl := blossom.New(relay, fmt.Sprintf("http://localhost%s", config.Port))
	bl.Store = blossom.EventStoreBlobIndexWrapper{Store: db, ServiceURL: bl.ServiceURL}

	bl.StoreBlob = append(bl.StoreBlob, bs.Store)
	bl.DeleteBlob = append(bl.DeleteBlob, bs.Delete)
	bl.LoadBlob = append(bl.LoadBlob, bs.Load)
	bl.RejectUpload = append(bl.RejectUpload, rejectUpload)

	if err := http.ListenAndServe(config.Port, relay); err != nil {
		log.Fatalf("Can't start the blossom server: %v", err)
	}
}

func rejectUpload(ctx context.Context, auth *nostr.Event, size int, ext string) (bool, string, int) {
	level, err := getWhitelistLevel(context.Background(), auth.PubKey)
	if err != nil {
		log.Printf("Can't read the whitelist from db: %v\n", err)
		return true, "internal: reading from database", 500
	}

	if level != 2 && level != 3 {
		return true, "blocked: you are not whitelisted", 403
	}

	return false, "", 200
}
