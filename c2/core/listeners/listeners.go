package listeners

import (
	"errors"
	"net"
	"sync"
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

	ErrAlreadyExists = errors.New("listener already exists")
)

func Insert(listener Listener) (err error) {
	defer listenersMux.Unlock()
	listenersMux.Lock()

	id := listener.GetID()

	if _, ok := listeners[id]; ok {
		return ErrAlreadyExists
	}

	listeners[id] = listener

	return nil
}

func Delete(listener Listener) {
	defer listenersMux.Unlock()
	listenersMux.Lock()

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
