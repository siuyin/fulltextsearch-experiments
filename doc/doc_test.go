package doc

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Init("../netflix_titles.csv")
	os.Exit(m.Run())
}

func TestRead(t *testing.T) {
	recs, err := r.Read()
	if err != nil {
		t.Error(err)
	}
	if n := len(recs); n != 12 {
		t.Error("incorrect length", n)
	}
}
