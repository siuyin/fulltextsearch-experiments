package doc

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

var r *csv.Reader

type docIndex int

//go:generate stringer -type=docIndex

// These indices allows you to use r[Director] instead of r[3].
const (
	ShowID docIndex = iota
	Type
	Title
	Director
	Cast
	Country
	DateAdded
	ReleaseYear
	Rating
	Duration
	ListedIn
	Description
)

// Init sets up the system to read a csv file into document records.
func Init(fn string) {
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}

	r = csv.NewReader(f)
	skipHeader()
}

func skipHeader() {
	if _, err := r.Read(); err != nil {
		log.Fatal(err)
	}
}

// Read returns the next document or
// nil if there is no more data to read.
func Read() []string {
	recs, err := r.Read()
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	return recs
}
