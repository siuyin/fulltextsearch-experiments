package embnats

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	svr "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/siuyin/dflt"
)

type Server struct {
	s  *svr.Server
	nc *nats.Conn
	js nats.JetStreamContext
	kv nats.KeyValue
}

func New() *Server {
	s, err := svr.NewServer(&svr.Options{Port: 4222, JetStream: true})
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

func (s *Server) connectSvr() {
	waitForSvrRdy(s.s)

	url := dflt.EnvString("NATS_SVR", "nats://localhost:4222")
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

func (s *Server) KVBucketNew(b string) {
	kv, err := s.js.CreateKeyValue(&nats.KeyValueConfig{Bucket: b})
	if err != nil {
		log.Fatal(err)
	}

	s.kv = kv
}

func (s *Server) KVPut(k string, v []string) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		log.Fatal(err)
	}

	s.kv.Put(k, buf.Bytes())
}

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
