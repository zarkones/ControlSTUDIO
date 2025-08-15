package listeners

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"

	profiles "github.com/zarkones/ControlPROFILE"
)

type Listener struct {
	ProfileID string
	Address   string
	Port      string
}

func (l *Listener) GetID() (ID string) {
	return l.ProfileID + l.GetHost()
}

func (l *Listener) GetHost() (host string) {
	return net.JoinHostPort(l.Address, l.Port)
}

var (
	listeners    = map[string]Listener{}
	listenersMux = &sync.Mutex{}
	services     = map[string]ListenerService{}

	ErrAlreadyExists = errors.New("listener already exists")
)

func Insert(listener Listener, profile *profiles.Profile) (err error) {
	defer listenersMux.Unlock()
	listenersMux.Lock()

	id := listener.GetID()

	if _, ok := listeners[id]; ok {
		return ErrAlreadyExists
	}

	listeners[id] = listener

	handler, err := handlersFromProfile(profile)
	if err != nil {
		return err
	}

	services[id] = ListenerService{
		Server: &http.Server{
			Addr:    listener.GetHost(),
			Handler: handler,
		},
	}

	go func() {
		log.Println("starting listener service at:", listener.GetHost())
		if err := services[id].Server.ListenAndServe(); err != nil {
			log.Println("listener error:", id, err)
			Delete(listener)
		}
	}()

	return nil
}

func Delete(listener Listener) {
	defer listenersMux.Unlock()
	listenersMux.Lock()

	if _, ok := listeners[listener.GetID()]; !ok {
		return
	}

	services[listener.GetID()].Server.Shutdown(context.Background())

	delete(listeners, listener.GetID())
}

func AsMap() (ls map[string]Listener) {
	defer listenersMux.Unlock()
	listenersMux.Lock()

	return listeners
}

func AsSlice() (ls []Listener) {
	defer listenersMux.Unlock()
	listenersMux.Lock()

	ls = make([]Listener, len(listeners))
	i := 0
	for _, listener := range listeners {
		ls[i] = listener
		i++
	}

	return ls
}
