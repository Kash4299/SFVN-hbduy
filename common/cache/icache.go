package cache

import (
	"time"
)

type IMemCache interface {
	Set(key string, value interface{}) error
	SetTTL(key string, value interface{}, t time.Duration) error
	Get(key string) (interface{}, error)
	Del(key string) error
	Close()
}

var MCache IMemCache
