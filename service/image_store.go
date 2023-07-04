package service

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"os"
	"sync"
)

type ImageStore interface {
	Save(laptopId string, imageType string, imageData bytes.Buffer) (string, error)
}

type DiskImageStore struct {
	mutex       sync.Mutex
	data        map[string]*ImageInfo
	imageFolder string
}

type ImageInfo struct {
	LaptopId string
	Type     string
	Path     string
}

func NewDiskImageStore(imageFolder string) ImageStore {
	return &DiskImageStore{
		data:        make(map[string]*ImageInfo),
		imageFolder: imageFolder,
	}
}

func (store *DiskImageStore) Save(laptopId string, imageType string, imageData bytes.Buffer) (string, error) {
	imageId := uuid.NewString()
	imagePath := fmt.Sprintf("%s/%s%s", store.imageFolder, imageId, imageType)
	file, err := os.Create(imagePath)
	if err != nil {
		return "", err
	}
	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", err
	}
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.data[imageId] = &ImageInfo{
		LaptopId: laptopId,
		Type:     imageType,
		Path:     imagePath,
	}
	return imageId, nil
}
