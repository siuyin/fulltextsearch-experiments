package main

import (
	"fmt"
	"log"
	"os"

	"github.com/siuyin/fulltextsearch-experiments/doc"
	"github.com/siuyin/fulltextsearch-experiments/embnats"
	"github.com/siuyin/fulltextsearch-experiments/idx"
)

func main() {
	fmt.Println("full text search with bluge")
	checkUsage()

	em := embnats.New()
	createFullTextIndex(em)

	select {} // wait forever
}

func createFullTextIndex(em *embnats.Server) {
	em.KVBucketNew("mov")
	doc.Init(os.Args[1])
	idx.InitWriter()
	defer idx.WriterClose()
	for r := doc.Read(); r != nil; r = doc.Read() {
		fmt.Println(r[doc.ShowID])
		idx.Add(r)
		em.KVPut(r[doc.ShowID], r)
	}
}

func checkUsage() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ", os.Args[0], " <csv file>")
	}
}
