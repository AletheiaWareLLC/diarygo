package main

import (
	"aletheiaware.com/bcgo"
	"aletheiaware.com/diarygo"
	"aletheiaware.com/spaceclientgo"
	"encoding/base64"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		usage()
		return
	}

	// Create Space Client
	c := spaceclientgo.NewSpaceClient()

	// Create Diary
	d := diarygo.NewDiary(c)

	// Get Node
	n, err := c.Node()
	if err != nil {
		log.Fatal(err)
	}

	// Refresh Diary Entries
	if err := d.Refresh(n); err != nil {
		log.Fatal(err)
	}

	switch args[0] {
	case "add":
		// Read data from system in
		reader := os.Stdin
		if len(args) > 1 {
			// Read data from file
			file, err := os.Open(args[1])
			if err != nil {
				log.Fatal(err)
			}
			reader = file
		} else {
			log.Println("Reading from stdin, use CTRL-D to terminate")
		}

		// Add New Diary Entry
		id, err := d.Add(n, reader)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Added:", id)
	case "list":
		list(d)
	case "show":
		var id string
		if len(args) > 1 {
			id = d.FindID(args[1])
		} else if l := d.Length(); l > 0 {
			// Choose latest
			id = d.ID(l - 1)
		}
		if id == "" {
			log.Println("Usage: diary show <id>")
			return
		}
		if err := show(c, n, d, id); err != nil {
			log.Fatal(err)
		}
	default:
		usage()
	}
}

func list(d diarygo.Diary) {
	for i := 0; i < d.Length(); i++ {
		id := d.ID(i)
		timestamp := d.Timestamp(id)
		meta := d.Meta(id)
		log.Println(id, bcgo.TimestampToString(timestamp), meta)
	}
}

func show(c spaceclientgo.SpaceClient, n bcgo.Node, d diarygo.Diary, id string) error {
	hash, err := base64.RawURLEncoding.DecodeString(id)
	if err != nil {
		return err
	}

	reader, err := c.ReadFile(n, hash)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	log.Println("ID:", id)
	log.Println("Timestamp:", bcgo.TimestampToString(d.Timestamp(id)))
	log.Println("Meta:", d.Meta(id))
	log.Println("Content:", string(bytes))
	return nil
}

func usage() {
	log.Println("Diary Usage:")
	log.Println("\tdiary - display usage")
	log.Println("\tdiary add - add a new diary entry from stdin")
	log.Println("\tdiary add [file] - add a new diary entry from file")
	log.Println("\tdiary list - display all entries")
	log.Println("\tdiary show - display latest entry")
	log.Println("\tdiary show [id] - display entry with given id")
}
