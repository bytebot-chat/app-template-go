package main

import (
	"context"
	"encoding/json"
	"flag"
	"strings"
	"time"

	"github.com/bytebot-chat/gateway-irc/model"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/satori/go.uuid"
)

// This is the app's named, as seen in the message it sends to gateways.
var app_name = "app-template"

// This is the env template used in configuration variables.
var app_env_prefix = "TEMPLATE_"

var logger zerolog.Logger

// Flags and their default values.
var addr = flag.String("redis", "localhost:6379", "Redis server address")
var inbound = flag.String("inbound", "irc-inbound", "Pubsub queue to listen for new messages")
var outbound = flag.String("outbound", "irc", "Pubsub queue for sending messages outbound")

func main() {
	flag.Parse()
	parseEnv()

	logger = createLogger() // file log.go, along with configuration

	ctx := context.Background()

	logger.Info().
		Str("App name", app_name).
		Str("Redis DB address", *addr).
		Str("Inbound queue", *inbound).
		Str("Outbound queue", *outbound).
		Msg("Connecting to gateway...")

	// We connect to the redis server
	rdb := redis.NewClient(&redis.Options{
		Addr: *addr,
		DB:   0,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		time.Sleep(3 * time.Second)
		logger.Warn().Msg("Ping timeout, trying again...")
		err := rdb.Ping(ctx).Err()
		if err != nil {
			panic(err)
		}
	}

	logger.Info().Msg("Connected to the gateway.")

	// Reading the incoming messages into a channel
	topic := rdb.Subscribe(ctx, *inbound)
	channel := topic.Channel()

	// We iterate over each new message and reply to it appropriately.
	for msg := range channel {
		m := &model.Message{}
		err := m.Unmarshal([]byte(msg.Payload))
		if err != nil {
			logger.Error().
				Str("message payload", msg.Payload).
				Err(err)
		}

		logger.Debug().
			RawJSON("Received message", []byte(msg.Payload)).
			Msg("Received message")

		// Below is the app logic, this template is just a parrot
		// that repeats every message
		m2 := model.Message{
			To:      m.To,
			From:    app_name,
			Content: m.Content,
			Metadata: model.Metadata{
				Source: app_name,
				Dest:   m.Metadata.Source,
				ID:     uuid.Must(uuid.NewV4(), *new(error)),
			},
		}
		// DMs go back to source, channel goes back to channel
		if !strings.HasPrefix(m.To, "#") {
			m2.To = m.From
		}

		stringMsg, err := json.Marshal(m2)
		if err != nil {
			logger.Error().
				Err(err).
				Msg("Couldn't marshall message")
		}

		rdb.Publish(ctx, *outbound, stringMsg)
		logger.Info().
			RawJSON("Message", stringMsg).
			Msg("Message sent")
	}
}
