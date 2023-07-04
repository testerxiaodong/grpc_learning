package service

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"grpc_learning/pb"
	"log"
	"sync"
)

var (
	ErrAlreadyExists = errors.New("record already exists")
)

type LaptopStore interface {
	Save(laptop *pb.Laptop) error
	Find(id string) (*pb.Laptop, error)
	Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error
}

type InMemoryLaptopStore struct {
	mutex sync.Mutex
	data  map[string]*pb.Laptop
}

func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

func (i *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	laptop := i.data[id]
	if laptop == nil {
		return nil, nil
	}
	return deepCopy(laptop)
}

func (i *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if i.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}
	// deep copy
	other, err := deepCopy(laptop)
	if err != nil {
		log.Printf("cannot copy laptop data: %v", laptop)
		return err
	}
	i.data[other.Id] = other
	log.Printf("save success, laptopId: %s", laptop.Id)
	return nil
}

func (i *InMemoryLaptopStore) Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	for _, laptop := range i.data {
		//time.Sleep(time.Second)
		log.Println("checking laptop id: ", laptop.GetId())
		if err := ContextError(ctx); err != nil {
			return nil
		}
		if isQualified(filter, laptop) {
			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}
			err = found(other)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isQualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}
	if laptop.GetCpu().GetNumberCores() < filter.GetMinCpuCore() {
		return false
	}
	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}
	if toBit(laptop.GetMemory()) < toBit(filter.GetMinRam()) {
		return false
	}
	return true
}

func toBit(memory *pb.Memory) uint64 {
	value := memory.GetValue()
	switch memory.GetUnit() {
	case pb.Memory_BIT:
		return value
	case pb.Memory_BYTE:
		return value << 3
	case pb.Memory_KILOBYTE:
		return value << 13
	case pb.Memory_MEGABYTE:
		return value << 23
	case pb.Memory_GIGABYTE:
		return value << 33
	case pb.Memory_TERABYTE:
		return value << 43
	default:
		return 0
	}
}

func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, err
	}
	return other, nil
}
