package filestorage

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

type FileStorager interface {
	GetFileList(id string) ([]FileInStorageInfo, error)
	GetFile(id string, filename string) ([]byte, error)
	PutItemImage(id string, filename string, file []byte) (string, error)
	PutCategoryImage(id string, filename string, file []byte) (string, error)
	DeleteItemImage(id string, filename string) error
	DeleteCategoryImage(id string, filename string) error
}

type OnDiskLocalStorage struct {
	serverURL string
	path      string
	logger    *zap.Logger
}

func NewOnDiskLocalStorage(url string, path string, logger *zap.Logger) *OnDiskLocalStorage {
	d := OnDiskLocalStorage{serverURL: url, path: path, logger: logger}
	return &d
}

func (imagestorage *OnDiskLocalStorage) PutItemImage(id string, filename string, file []byte) (filePath string, err error) {
	imagestorage.logger.Debug("Enter in filestorage PutItemImage()")
	_, err = os.Stat(imagestorage.path + "items/" + id)
	if os.IsNotExist(err) {
		err = os.Mkdir(imagestorage.path + "items/" + id, 0700)
		if err != nil {
			imagestorage.logger.Debug(fmt.Sprintf("error on create dir for save image %v", err))
			return "", fmt.Errorf("error on create dir for save image: %w", err)
		}
	}
	filePath = imagestorage.path + "items/" + id + "/" + filename
	if err := os.WriteFile(filePath, file, os.ModePerm); err != nil {
		imagestorage.logger.Debug(fmt.Sprintf("error on filestorage put file: %v", err))
		return "", fmt.Errorf("error on filestorage put file: %w", err)
	}
	urlPath := imagestorage.serverURL + "/files/items/" + id + "/" + filename
	imagestorage.logger.Debug(urlPath)
	return urlPath, nil
}

func (imagestorage *OnDiskLocalStorage) PutCategoryImage(id string, filename string, file []byte) (filePath string, err error) {
	imagestorage.logger.Debug("Enter in filestorage PutCategoryFile()")
	_, err = os.Stat(imagestorage.path + "categories/" + id)
	if os.IsNotExist(err) {
		err = os.Mkdir(imagestorage.path+"categories/"+id, 0700)
		if err != nil {
			imagestorage.logger.Debug(fmt.Sprintf("error on create dir for save image %v", err))
			return "", fmt.Errorf("error on create dir for save image: %w", err)
		}
	}
	filePath = imagestorage.path + "categories/" + id + "/" + filename
	if err := os.WriteFile(filePath, file, os.ModePerm); err != nil {
		imagestorage.logger.Debug(fmt.Sprintf("error on filestorage put file: %v", err))
		return "", fmt.Errorf("error on filestorage put file: %w", err)
	}
	urlPath := imagestorage.serverURL + "/files/categories/" + id + "/" + filename
	imagestorage.logger.Debug(urlPath)
	return urlPath, nil
}

func (imagestorage *OnDiskLocalStorage) DeleteItemImage(id string, filename string) error {
	err := os.Remove(imagestorage.path + "items/" + id + "/" + filename)
	if err != nil {
		imagestorage.logger.Debug(fmt.Sprintf("error on delete file: %v", err))
		return fmt.Errorf("error on delete file: %w", err)
	}
	return nil
}

func (imagestorage *OnDiskLocalStorage) DeleteCategoryImage(id string, filename string) error {
	err := os.Remove(imagestorage.path + "categories/" + id + "/" + filename)
	if err != nil {
		imagestorage.logger.Debug(fmt.Sprintf("error on delete file: %v", err))
		return fmt.Errorf("error on delete file: %w", err)
	}
	return nil
}

func (imagestorage *OnDiskLocalStorage) GetFile(id string, filename string) ([]byte, error) {
	if f, err := os.ReadFile(imagestorage.path + id + "-" + filename); err == nil {
		return f, nil
	}

	return nil, FileNotFoundError{}
}

func (imagestorage *OnDiskLocalStorage) GetFileList(id string) ([]FileInStorageInfo, error) {
	result := make([]FileInStorageInfo, 0)
	dir, err := os.ReadDir(imagestorage.path)
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
