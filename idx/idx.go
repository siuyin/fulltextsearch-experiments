package idx

import (
	"context"
	"fmt"
	"log"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/search"
	querystr "github.com/blugelabs/query_string"
	"github.com/siuyin/dflt"
	"github.com/siuyin/fulltextsearch-experiments/doc"
)

var (
	config bluge.Config
	w      *bluge.Writer
	r      *bluge.Reader
	err    error
)

// InitWriter sets up a search index writer.
func InitWriter() {
	config = bluge.DefaultConfig(dflt.EnvString("BLUGE_PATH", "./blugeidx"))
	w, err = bluge.OpenWriter(config)
	if err != nil {
		log.Fatal(err)
	}
}

// WriterClose closes the search index writer.
func WriterClose() {
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// InitReader sets up a search index reader.
func InitReader() {
	config = bluge.DefaultConfig(dflt.EnvString("BLUGE_PATH", "./blugeidx"))
	config.DefaultSearchField = doc.Title.String()
	//config.DefaultSearchField = doc.Description.String()
	r, err = bluge.OpenReader(config)
	if err != nil {
		log.Fatal(err)
	}
}

// ReaderClose closes the search index reader.
func ReaderClose() {
	err = r.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// Add adds a document to be indexed.
func Add(rec []string) {
	d := newDoc(rec)

	err = w.Update(d.ID(), d)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("   Added: ", rec[doc.Title])
}

// AddBatch add n documents to the search index and returns the number of entries added.
func AddBatch(n int) int {
	b := bluge.NewBatch()

	i := 0
	for rec := doc.Read(); rec != nil; rec = doc.Read() {
		d := newDoc(rec)
		b.Update(d.ID(), d)

		i++
		if i == n {
			break
		}
	}
	w.Batch(b)

	return i
}

func newDoc(rec []string) *bluge.Document {
	d := bluge.NewDocument(rec[doc.ShowID])
	d.AddField(bluge.NewTextField(doc.Type.String(), rec[doc.Type]))
	d.AddField(bluge.NewTextField(doc.Title.String(), rec[doc.Title]).StoreValue()) // store this value, IDs are automatically stored
	d.AddField(bluge.NewTextField(doc.Director.String(), rec[doc.Director]))
	d.AddField(bluge.NewTextField(doc.Cast.String(), rec[doc.Cast]))
	d.AddField(bluge.NewTextField(doc.Country.String(), rec[doc.Country]))
	d.AddField(bluge.NewTextField(doc.DateAdded.String(), rec[doc.DateAdded]))
	d.AddField(bluge.NewTextField(doc.ReleaseYear.String(), rec[doc.ReleaseYear]))
	d.AddField(bluge.NewTextField(doc.Rating.String(), rec[doc.Rating]))
	d.AddField(bluge.NewTextField(doc.Duration.String(), rec[doc.Duration]))
	d.AddField(bluge.NewTextField(doc.ListedIn.String(), rec[doc.ListedIn]))
	d.AddField(bluge.NewTextField(doc.Description.String(), rec[doc.Description]).StoreValue())

	return d
}

// TopNSearch searches for the top n matches given search term s.
func TopNSearch(n int, s string) {
	// TODO: use this analyser as a query string option:
	// https://go.dev/play/p/G_KNuNm0c9p
	// https://pkg.go.dev/github.com/blugelabs/bluge@v0.1.9/analysis/analyzer#NewKeywordAnalyzer
	req := getRequest(n, s)

	dmi, err := r.Search(context.Background(), req) // dmi: document match iterator
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("total hits: %d (%v) %v\n\n", dmi.Aggregations().Count(), dmi.Aggregations().Duration(), dmi.Aggregations().Metric("max_score"))

	iterateAndShow(dmi)
}

func getRequest(n int, s string) *bluge.TopNSearch {
	userQuery, err := querystr.ParseQueryString(s, querystr.DefaultOptions())
	if err != nil {
		log.Fatalf("errror parsing query string '%s': %v", s, err)
	}

	q := bluge.NewBooleanQuery().AddMust(userQuery)
	req := bluge.NewTopNSearch(n, q).WithStandardAggregations()
	log.Println("searching top ", n, s)

	return req
}

func iterateAndShow(dmi search.DocumentMatchIterator) {
	match, err := dmi.Next()
	for err == nil && match != nil {
		showMatches(match)

		match, err = dmi.Next()
	}

	if err != nil {
		log.Fatalf("error iterating results: %v", err)
	}
}
func showMatches(match *search.DocumentMatch) {
	entry := map[string]string{}
	err = match.VisitStoredFields(func(field string, value []byte) bool {
		entry[field] = string(value)
		return true
	})

	if err != nil {
		log.Fatalf("error accessing stored fields: %v", err)
	}
	fmt.Printf("%s %v %s\n%s\n----\n", entry["_id"], match.Score, entry["Title"], entry["Description"])
}
