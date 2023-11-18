package topic

import (
	"fmt"
	"sync"
	"uno/service/subscriber"

	"github.com/google/uuid"
)

var bloc = sync.Map{}

type Entry struct {
	ID          string              `json:"id,omitempty"`
	Title       string              `json:"title,omitempty"`
	Subscribers []*subscriber.Entry `json:"subscribers,omitempty"`
	CreateAt    int64               `json:"create_at,omitempty"`
	Status      int                 `json:"status,omitempty"`
}

func Create(entry *Entry) error {
	if entry == nil {
		return fmt.Errorf("entry is nil")
	}

	entry.ID = uuid.New().String()

	bloc.Store(entry.ID, entry)

	go func() {
		// TODO prepare data for internal channel of subscribers
	}()

	return nil
}

func Get(id string) (*Entry, error) {
	if id == "" {
		return nil, fmt.Errorf("id is empty")
	}

	entry, ok := bloc.Load(id)
	if !ok {
		return nil, fmt.Errorf("topic not found")
	}

	return entry.(*Entry), nil
}
