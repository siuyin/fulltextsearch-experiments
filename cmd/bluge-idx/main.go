package main

import (
	"fmt"
	"log"
	"os"

	"github.com/siuyin/dflt"
	"github.com/siuyin/fulltextsearch-experiments/doc"
	"github.com/siuyin/fulltextsearch-experiments/embnats"
	"github.com/siuyin/fulltextsearch-experiments/idx"
)

func main() {
	fmt.Println("full text search with bluge")
	checkUsage()

	createFullTextIndex()
}

func createFullTextIndex() {
	doc.Init(os.Args[1])
	idx.InitWriter()
	defer idx.WriterClose()

	batchAdd(batchSize())
	fmt.Println("\nindexed and ready for queries")
}

func batchSize() int {
	n, err := dflt.EnvInt("BLUGE_BATCHSIZE", 500)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

func singleAdd(em *embnats.Server) {
	for r := doc.Read(); r != nil; r = doc.Read() {
		fmt.Println(r[doc.ShowID])
		idx.Add(r)
		em.KVPut(r[doc.ShowID], r)
	}
}

func batchAdd(size int) {
	n := 0
	for {
		n = idx.AddBatch(size)
		fmt.Print(".")
		if n != size {
			break
		}
	}
}

func checkUsage() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ", os.Args[0], " <csv file>")
	}
}
