package filestorage

import (
	"os"
	"time"
)

type FileInStorageInfo struct {
	Id         string `json:"Id"`
	Name       string `json:"Name"`
	Path       string `json:"Path"`
	CreateDate string `json:"CreateDate"`
	ModifyDate string `json:"ModifyDate"`
}

type FileStorager interface {
	GetFileList(id string) ([]FileInStorageInfo, error)
	GetFile(id string, filename string) ([]byte, error)
	PutFile(id string, filename string, file []byte) error
}

type ImMemoryLocalStorage struct {
	files []FileInStorageInfo
}
type FileNotFoundError struct{}

func (f FileNotFoundError) Error() string {
	return "File not found"
}

type WriteFileError struct{}

func (f WriteFileError) Error() string {
	return "WriteFileError"
}

func (i *ImMemoryLocalStorage) PutFile(id string, filename string, file []byte) error {
	if i.files == nil {
		i.files = make([]FileInStorageInfo, 0)
	}
	if err := os.WriteFile(id+"-"+filename, file, os.ModePerm); err != nil {
		return err
	}
	i.files = append(i.files, FileInStorageInfo{
		Id:         id,
		Name:       filename,
		Path:       id + "-" + filename,
		CreateDate: time.Now().GoString(),
		ModifyDate: time.Now().GoString(),
	})
	return nil
}

func (i *ImMemoryLocalStorage) GetFile(id string, filename string) ([]byte, error) {
	for _, file := range i.files {
		if file.Id == id && file.Name == filename {
			if f, err := os.ReadFile(file.Path); err == nil {
				return f, nil
			}
		}
	}
	return nil, FileNotFoundError{}
}

func (i *ImMemoryLocalStorage) GetFileList(id string) ([]FileInStorageInfo, error) {
	result := make([]FileInStorageInfo, 0, len(i.files))
	for _, v := range i.files {
		if v.Id == id {
			result = append(result, v)
		}
	}

	return result, nil
}
