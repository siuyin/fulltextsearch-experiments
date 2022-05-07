package main

import (
	"fmt"
	"log"
	"os"

	"github.com/siuyin/fulltextsearch-experiments/idx"
)

func main() {
	fmt.Println("full text search with bluge")
	checkUsage()

	topNSearch(10, os.Args[1])

}

func topNSearch(n int, s string) {
	idx.InitReader()
	defer idx.ReaderClose()

	idx.TopNSearch(n, s)
}

func checkUsage() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ", os.Args[0], "<search term>")
	}
}
