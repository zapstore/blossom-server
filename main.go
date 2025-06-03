package main

import (
	"context"
	"log"
	"net/http"
	"path"

	"github.com/fiatjaf/eventstore/sqlite3"
	"github.com/fiatjaf/khatru"
	"github.com/fiatjaf/khatru/blossom"
	"github.com/kehiy/blobstore/disk"
	"github.com/nbd-wtf/go-nostr"
)

const limitInMB int = 600 * 1024 * 1024

var (
	config Config
	db     *sqlite3.SQLite3Backend
)

func main() {
	log.SetPrefix("Blossom-Server ")
	log.Printf("Running %s\n", StringVersion())

	relay := khatru.NewRelay()

	relay.RejectEvent = append(relay.RejectEvent, func(context.Context, *nostr.Event) (reject bool, msg string) {
		return true, "blocked: not a relay"
	})

	relay.RejectFilter = append(relay.RejectFilter, func(context.Context, nostr.Filter) (reject bool, msg string) {
		return true, "blocked: not a relay"
	})

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

	bl := blossom.New(relay, config.ServerURL)
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
	if size > limitInMB {
		return true, "blocked: max upload limit is 600MB", 400
	}

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
