package main

import (
	"encoding/json"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/mroth/sseserver"
	"github.com/ryanlower/setting"
)

type Config struct {
	Port  string `env:"PORT" default:"3001"`
	Redis struct {
		Host     string `env:"REDIS_HOST"`
		Password string `env:"REDIS_PASSWORD"`
	}
}

type Server struct {
	conf Config
	sse  sseserver.Server
}

type Hit struct {
	Code string
}

func main() {
	config := new(Config)
	setting.Load(config)

	sse := sseserver.NewServer()

	server := &Server{
		conf: *config,
		sse:  *sse,
	}

	go server.subscribe()

	log.Printf("Reporting on port %v ...", config.Port)
	server.sse.Serve(":" + config.Port)
}

func (s *Server) subscribe() {
	log.Print("subscribe")

	c, err := redis.Dial("tcp", s.conf.Redis.Host)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer c.Close()

	if s.conf.Redis.Password != "" {
		c.Do("AUTH", s.conf.Redis.Password)
	}

	psc := redis.PubSubConn{c}
	psc.Subscribe("hits")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			var hit Hit
			json.Unmarshal(v.Data, &hit)

			s.broadcastHit(hit)
		}
	}
}

func (s *Server) broadcastHit(h Hit) {
	log.Print("hits: ", h.Code)
	s.sse.Broadcast <- sseserver.SSEMessage{"", []byte(h.Code), "/hits"}
}
