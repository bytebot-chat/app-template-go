package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/bytebot-chat/gateway-irc/model"
	"github.com/go-redis/redis/v8"
	"github.com/satori/go.uuid"
)

// Flags and their default values.
var addr = flag.String("redis", "localhost:6379", "Redis server address")
var inbound = flag.String("inbound", "irc-inbound", "Pubsub queue to listen for new messages")
var outbound = flag.String("outbound", "irc", "Pubsub queue for sending messages outbound")

func main() {
	flag.Parse()
	ctx := context.Background()

	// We connect to the redis server
	rdb := redis.NewClient(&redis.Options{
		Addr: *addr,
		DB:   0,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		time.Sleep(3 * time.Second)
		err := rdb.Ping(ctx).Err()
		if err != nil {
			panic(err)
		}
	}

	// Reading the incoming messages into a channel
	topic := rdb.Subscribe(ctx, *inbound)
	channel := topic.Channel()

	// We iterate over each new message and reply to it appropriately.
	for msg := range channel {
		m := &model.Message{}
		err := m.Unmarshal([]byte(msg.Payload))
		if err != nil {
			fmt.Println(err)
		}

		// Below is the app logic, this template is just a parrot
		// that repeats every message
		m2 := model.Message{
			To:      m.From,
			Content: m.Content,
			Metadata: model.Metadata{
				Source: m.Metadata.Dest,
				Dest:   m.Metadata.Source,
				ID:     uuid.Must(uuid.NewV4(), *new(error)),
			},
		}

		stringMsg, _ := json.Marshal(m2)
		fmt.Println(string(stringMsg))
		rdb.Publish(ctx, *outbound, stringMsg)
	}
}