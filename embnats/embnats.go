package embnats

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	svr "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/siuyin/dflt"
)

// Server represents an embedded NATS server.
type Server struct {
	s  *svr.Server
	nc *nats.Conn
	js nats.JetStreamContext
	kv nats.KeyValue
}

// New returns a reference to an embedded NATS server.
func New() *Server {
	s, err := svr.NewServer(&svr.Options{Port: natsPort(), JetStream: true,
		StoreDir: dflt.EnvString("NATS_STOREDIR", ".")})
	if err != nil {
		log.Fatal(err)
	}

	em := Server{}
	em.s = s

	s.Start()
	em.connectSvr()
	em.initJetStream()
	log.Println("NATS server started")
	return &em
}

func natsPort() int {
	port, err := dflt.EnvInt("NATS_PORT", 4222)
	if err != nil {
		log.Fatal(err)
	}
	return port
}

func (s *Server) connectSvr() {
	waitForSvrRdy(s.s)

	url := fmt.Sprintf("nats://localhost:%d", natsPort())
	var err error
	s.nc, err = nats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to NATS server: ", url)
}

func (s *Server) initJetStream() {
	var err error
	s.js, err = s.nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("JetStream event streaming persistence initialised")
}

func waitForSvrRdy(s *svr.Server) {
	for rdy := s.ReadyForConnections(time.Second); !rdy; rdy = s.ReadyForConnections(time.Second) {
	}
}

// KVBucketNew optionally creates and returns a key-value bucket to hold key-value pairs.
func (s *Server) KVBucketNew(b string) {
	kv, err := s.js.CreateKeyValue(&nats.KeyValueConfig{Bucket: b})
	if err != nil {
		log.Fatal(err)
	}

	s.kv = kv
}

// KVPut stores value v associated with key k.
func (s *Server) KVPut(k string, v []string) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := s.kv.Put(k, buf.Bytes()); err != nil {
		log.Fatal(err)
	}
}

// KVGet retrieve the value associated with key k.
func (s *Server) KVGet(k string) []string {
	entry, err := s.kv.Get(k)
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(entry.Value())
	dec := gob.NewDecoder(buf)
	recs := []string{}
	err = dec.Decode(&recs)
	if err != nil {
		log.Fatal(err)
	}

	return recs
}
