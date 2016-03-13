package main

import "sync"

var hubs *hubMap

type hubMap struct {
	m  map[string](*hub)
	mu sync.RWMutex
}

func HubsInitialize() {
	hubs = &hubMap{m: map[string]*hub{}}
}

func GetHub(id string) *hub {
	hubs.mu.RLock()

	//Если хаб уже создан, то просто его вернем
	if hubs.m[id] != nil {
		defer hubs.mu.RUnlock()
		return hubs.m[id]
	}
	hubs.mu.RUnlock()

	//Если хаб еще не создан, создадим его и запустим
	hubs.mu.Lock()
	defer hubs.mu.Unlock()

	h := NewHub(id)

	hubs.m[id] = h
	go h.run()

	return h
}
