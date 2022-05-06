package cache

import (
	"log"
	"sync"
)

type Cache interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error
	GetStatus() Status
}

type inMemoryCache struct {
	c     map[string][]byte
	mutex sync.RWMutex
	Status
}

func (i *inMemoryCache) Set(k string, v []byte) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.c[k] = v
	i.add(k, v)
	return nil
}

func (i *inMemoryCache) Get(k string) ([]byte, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	return i.c[k], nil
}

func (i *inMemoryCache) Del(k string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	v, exist := i.c[k]
	if exist {
		delete(i.c, k)
		i.del(k, v)
	}
	return nil
}

func (i *inMemoryCache) GetStatus() Status {
	return i.Status
}

type Status struct {
	Count     int64
	KeySize   int64
	ValueSize int64
}

func (s *Status) add(k string, v []byte) {
	s.Count += 1
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
}

func (s *Status) del(k string, v []byte) {
	s.Count -= 1
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
}

func New(typ string) Cache {
	var c Cache
	if typ == "inmemory" {
		c = newInMemoryCache()
	}
	if c == nil {
		panic("unknow cache type " + typ)
	}
	log.Println(typ, "ready to serve")
	return c
}

func newInMemoryCache() Cache {
	c := &inMemoryCache{}
	c.c = make(map[string][]byte)
	c.mutex = sync.RWMutex{}
	c.Status = Status{}
	return c
}
