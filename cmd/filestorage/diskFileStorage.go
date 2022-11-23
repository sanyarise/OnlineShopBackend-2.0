package filestorage

import (
	"os"
)

type OnDiskLocalStorage struct {
	path string
}

func NewOnDiskLocalStorage(path string) *OnDiskLocalStorage {
	d := OnDiskLocalStorage{path: path}
	return &d
}

func (i *OnDiskLocalStorage) PutFile(id string, filename string, file []byte) error {
	if err := os.WriteFile(i.path+id+"-"+filename, file, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (i *OnDiskLocalStorage) GetFile(id string, filename string) ([]byte, error) {
	if f, err := os.ReadFile(i.path + id + "-" + filename); err == nil {
		return f, nil
	}

	return nil, FileNotFoundError{}
}

func (i *OnDiskLocalStorage) GetFileList(id string) ([]FileInStorageInfo, error) {
	result := make([]FileInStorageInfo, 0)
	dir, err := os.ReadDir(i.path)
	if err != nil {
		return nil, err
	}
	for _, v := range dir {
		if !v.IsDir() {
			if info, err := v.Info(); err != nil {
				result = append(result, FileInStorageInfo{
					Id:         id,
					Name:       info.Name(),
					Path:       info.Name(),
					CreateDate: info.ModTime().String(),
					ModifyDate: info.ModTime().String(),
				})
			}
		}
	}

	return result, nil
}
