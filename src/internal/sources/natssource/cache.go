package natssource

import (
	"context"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sandrolain/event-runner/src/config"
	"github.com/sandrolain/event-runner/src/internal/utils"
)

type NatsEventCache struct {
	slog   *slog.Logger
	config config.Cache
	nats   *nats.Conn
	js     *jetstream.JetStream
	ctx    context.Context
	ctxC   context.CancelFunc
	kv     *jetstream.KeyValue
}

func (s *NatsEventCache) init() (err error) {
	var ttl time.Duration
	if s.config.Ttl == "" {
		ttl, err = time.ParseDuration(s.config.Ttl)
		if err != nil {
			return
		}
	}

	s.ctx, s.ctxC = context.WithCancel(context.TODO())
	kv, err := (*s.js).CreateOrUpdateKeyValue(s.ctx, jetstream.KeyValueConfig{
		Bucket: s.config.Bucket,
		TTL:    ttl,
		//TODO other options?
	})
	if err != nil {
		return
	}
	s.kv = &kv
	return
}

func (s *NatsEventCache) Close() (err error) {
	s.ctxC()
	return
}

func (s *NatsEventCache) Get(key string) (res any, err error) {
	val, err := (*s.kv).Get(s.ctx, key)
	if err != nil {
		if err == jetstream.ErrKeyNotFound {
			return nil, nil
		}
		return
	}
	byt := val.Value()
	err = utils.Unmarshal(s.config.Marshal, byt, &res)
	return
}

func (s *NatsEventCache) Set(key string, value any) (err error) {
	byt, err := utils.Marshal(s.config.Marshal, value)
	if err != nil {
		return
	}
	_, err = (*s.kv).Put(s.ctx, key, byt)
	return
}

func (s *NatsEventCache) Del(key string) (err error) {
	err = (*s.kv).Delete(s.ctx, key)
	return
}
