// Package metre is used to schedule end execute corn jobs in a simplified fashion
package metre

import (
    "github.com/garyburd/redigo/redis"
)

type Cache struct {
    ConnPool *redis.Pool
}

// Set sets the cache valur for the provided key
func (c Cache) Set(key string, val string) (interface{}, error) {
    conn := c.ConnPool.Get()
    defer conn.Close()

    return conn.Do("SET", key, val)
}

// Get gets the cache value for the provided key
func (c Cache) Get(key string) (string, error) {
    conn := c.ConnPool.Get()
    data, err := redis.String(conn.Do("GET", key))
    if err != nil {
        return "", err
    }

    return data, nil
}

// New acts as a queue constructor
func NewCache(uri string) (Cache, error) {
    u := "redis://" + uri
    p := &redis.Pool{
        MaxIdle: 80,
        MaxActive: 12000,
        Dial: func() (redis.Conn, error) {
            c, err := redis.DialURL(u)
            if err != nil {
                panic(err.Error())
            }
            return c, err
        },
    }

    return Cache{p}, nil
}
