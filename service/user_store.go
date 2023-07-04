package service

import (
	"sync"
)

type UserStore interface {
	Save(user *User) error
	Find(username string) (*User, error)
}

type InMemoryUserStore struct {
	mutex sync.Mutex
	data  map[string]*User
}

func NewInMemoryUserStore() UserStore {
	return &InMemoryUserStore{
		data: make(map[string]*User),
	}
}

func (store *InMemoryUserStore) Save(user *User) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if store.data[user.Username] != nil {
		return ErrAlreadyExists
	}
	other := user.Clone()
	store.data[other.Username] = other
	return nil
}

func (store *InMemoryUserStore) Find(username string) (*User, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	user := store.data[username]
	if user == nil {
		return nil, nil
	}
	return user.Clone(), nil
}
