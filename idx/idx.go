package idx

import (
	"context"
	"fmt"
	"log"

	"github.com/blugelabs/bluge"
	"github.com/siuyin/dflt"
	"github.com/siuyin/fulltextsearch-experiments/doc"
	"github.com/siuyin/fulltextsearch-experiments/embnats"
)

var (
	config bluge.Config
	w      *bluge.Writer
	r      *bluge.Reader
	err    error
)

func InitWriter() {
	config = bluge.DefaultConfig(dflt.EnvString("BLUGE_PATH", "./blugeidx"))
	w, err = bluge.OpenWriter(config)
	if err != nil {
		log.Fatal(err)
	}
}

func WriterClose() {
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func InitReader() {
	config = bluge.DefaultConfig(dflt.EnvString("BLUGE_PATH", "./blugeidx"))
	config.DefaultSearchField = doc.Title.String()
	//config.DefaultSearchField = doc.Description.String()
	r, err = bluge.OpenReader(config)
	if err != nil {
		log.Fatal(err)
	}
}

func ReaderClose() {
	err = r.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func Add(rec []string) {
	d := bluge.NewDocument(rec[doc.ShowID])
	d.AddField(bluge.NewTextField(doc.Type.String(), rec[doc.Type]))
	d.AddField(bluge.NewTextField(doc.Title.String(), rec[doc.Title]))
	d.AddField(bluge.NewTextField(doc.Director.String(), rec[doc.Director]))
	d.AddField(bluge.NewTextField(doc.Cast.String(), rec[doc.Cast]))
	d.AddField(bluge.NewTextField(doc.Country.String(), rec[doc.Country]))
	d.AddField(bluge.NewTextField(doc.DateAdded.String(), rec[doc.DateAdded]))
	d.AddField(bluge.NewTextField(doc.ReleaseYear.String(), rec[doc.ReleaseYear]))
	d.AddField(bluge.NewTextField(doc.Rating.String(), rec[doc.Rating]))
	d.AddField(bluge.NewTextField(doc.Duration.String(), rec[doc.Duration]))
	d.AddField(bluge.NewTextField(doc.ListedIn.String(), rec[doc.ListedIn]))
	d.AddField(bluge.NewTextField(doc.Description.String(), rec[doc.Description]))

	err = w.Update(d.ID(), d)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("   Added: ", rec[doc.Title])
}

func TopNSearch(n int, s string) {
	em := embnats.New()
	em.KVBucketNew(dflt.EnvString("NATS_BUCKET", "mov"))
	q := bluge.NewMatchQuery(s)
	req := bluge.NewTopNSearch(n, q)
	log.Println("searching top ", n, s)

	dmi, err := r.Search(context.Background(), req) // dmi: document match iterator
	if err != nil {
		log.Fatal(err)
	}

	next, err := dmi.Next()
	for err == nil && next != nil {
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			if field == "_id" {
				rec := em.KVGet(string(value))
				fmt.Println(string(value), next.Score, rec[doc.Type], rec[doc.Title])
				fmt.Println("   ", rec[doc.Description])
			}
			return true
		})

		if err != nil {
			log.Fatalf("error accessing stored fields: %v", err)
		}
		next, err = dmi.Next()
	}

	if err != nil {
		log.Fatalf("error iterating results: %v", err)
	}
}
