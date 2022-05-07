package idx

import (
	"os"
	"testing"

	"github.com/siuyin/fulltextsearch-experiments/doc"
)

func TestAddBatch(t *testing.T) {
	size := 3
	n := 0
	for {
		n = AddBatch(size)
		if n != size {
			break
		}
	}

	exp := 8 % size
	if n != exp {
		t.Errorf("expected %d, got %d", exp, n)
	}
}

func TestMain(m *testing.M) {
	InitWriter()
	doc.Init("testdata/dat.csv")
	os.Exit(m.Run())
}
