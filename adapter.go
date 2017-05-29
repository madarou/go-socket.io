package socketio

import (
	"sync"
)

// BroadcastAdaptor is the adaptor to handle broadcasts.
type BroadcastAdaptor interface {

	// Join causes the conn to join a room.
	Join(room string,conn Conn) error

	// Leave causes the conn to leave a room.
	Leave(room string, conn Conn) error

	// Send will send an event with args to the room. If "ignore" is not nil, the event will be excluded from being sent to "ignore".
	Send(ignore Conn, room, event string, args ...interface{}) error
}

type broadcast struct {
	m map[string]map[string]Conn
	sync.RWMutex
}

func newBroadcastDefault() BroadcastAdaptor {
	return &broadcast{
		m: make(map[string]map[string]Conn),
	}
}

//加入房间
func (b *broadcast) Join(room string, conn Conn) error {
	b.Lock()
	conns, ok := b.m[room]
	if !ok {
		conns = make(map[string]Conn)
	}
	conns[conn.ID()] = conn
	b.m[room] = conns
	b.Unlock()
	return nil
}

func (b *broadcast) Leave(room string, conn Conn) error{
	b.Lock()
	defer b.Unlock()
	conns, ok := b.m[room]
	if !ok {
		return nil
	}
	delete(conns, conn.ID())
	if len(conns) == 0 {
		delete(b.m, room)
		return nil
	}
	b.m[room] = conns
	return nil
}

//谁发就ignore谁
func (b *broadcast) Send(ignore Conn, room, event string, args ...interface{}) error {
	b.RLock()
	conns := b.m[room]
	for id, s := range conns {
		if ignore != nil && ignore.ID() == id {
			continue
		}
		s.Emit(event, args...)
	}
	b.RUnlock()
	return nil
}
