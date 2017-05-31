package socketio

import (
	"sync"
	"fmt"
)

// BroadcastAdaptor is the adaptor to handle broadcasts.
type BroadcastAdaptor interface {

	// Join causes the conn to join a room.
	Join(room string,conn Conn) error

	// Leave causes the conn to leave a room.
	Leave(room string, conn Conn) error

	// Send will send an event with args to the room. If "ignore" is not nil, the event will be excluded from being sent to "ignore".
	BroadcastTo(ignore Conn, room, event string, args ...interface{}) error

	// Find target belong to which room
	Belong(target Conn)(string,error)

	// List
	List()error
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
	fmt.Println(conn.ID()," leaves room "+room)
	return nil
}

//谁发就ignore谁
func (b *broadcast) BroadcastTo(ignore Conn, room, event string, args ...interface{}) error {
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

func (b *broadcast) Belong(target Conn)(string,error){
	b.RLock()
	for room, conns:=range b.m{
		for id, _:=range conns{
			if target.ID() == id{
				return room,nil
			}
		}
	}
	b.RUnlock()
	return "",nil
}

func (b *broadcast) List()error{
	b.RLock()
	for room, conns:=range b.m{
		fmt.Printf("room name=%s: ",room)
		for id, _:=range conns{
			fmt.Printf("%s,",id)
		}
		fmt.Println()
	}
	b.RUnlock()
	return nil
}