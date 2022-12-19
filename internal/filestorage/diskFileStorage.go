package filestorage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

type FileStorager interface {
	GetFileList() ([]FileInStorageInfo, error)
	PutItemImage(id string, filename string, file []byte) (string, error)
	PutCategoryImage(id string, filename string, file []byte) (string, error)
	DeleteItemImage(id string, filename string) error
	DeleteCategoryImage(id string, filename string) error
	DeleteCategoryImageById(id string) error
}

type FileInStorageInfo struct {
	Name       string `json:"Name"`
	Path       string `json:"Path"`
	CreateDate string `json:"CreateDate"`
	ModifyDate string `json:"ModifyDate"`
}

type OnDiskLocalStorage struct {
	serverURL string
	path      string
	logger    *zap.Logger
}

func NewOnDiskLocalStorage(url string, path string, logger *zap.Logger) *OnDiskLocalStorage {
	logger.Debug("Enter in NewOnDiskLocalStorage()")
	d := OnDiskLocalStorage{serverURL: url, path: path, logger: logger}
	return &d
}

func (imagestorage *OnDiskLocalStorage) PutItemImage(id string, filename string, file []byte) (filePath string, err error) {
	imagestorage.logger.Debug("Enter in filestorage PutItemImage()")
	_, err = os.Stat(imagestorage.path + "items/" + id)
	if os.IsNotExist(err) {
		err = os.Mkdir(imagestorage.path+"items/"+id, 0700)
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
	imagestorage.logger.Debug("Enter in filestorage PutCategoryImage()")
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
	imagestorage.logger.Debug("Enter in filestorage DeleteItemImage()")
	imagestorage.logger.Debug(fmt.Sprintf("name of deleting image: %s", filename))
	err := os.Remove(imagestorage.path + "items/" + id + "/" + filename)
	if err != nil {
		imagestorage.logger.Debug(fmt.Sprintf("error on delete file: %v", err))
		return fmt.Errorf("error on delete file: %w", err)
	}
	return nil
}

func (imagestorage *OnDiskLocalStorage) DeleteCategoryImage(id string, filename string) error {
	imagestorage.logger.Debug("Enter in filestorage DeleteCategoryImage()")
	imagestorage.logger.Debug(fmt.Sprintf("name of deleting image: %s", filename))
	err := os.Remove(imagestorage.path + "categories/" + id + "/" + filename)
	if err != nil {
		imagestorage.logger.Debug(fmt.Sprintf("error on delete file: %v", err))
		return fmt.Errorf("error on delete file: %w", err)
	}
	return nil
}

func (imagestorage *OnDiskLocalStorage) DeleteCategoryImageById(id string) error {
	imagestorage.logger.Debug("Enter in filestorage DeleteCategoryImageById()")
	imagestorage.logger.Debug(fmt.Sprintf("path is: %s", imagestorage.path+"categories/"+id))
	err := filepath.Walk(imagestorage.path+"categories/"+id, func(path string, info os.FileInfo, err error) error {
		imagestorage.logger.Debug(fmt.Sprintf("isDir: %t, name: %s", info.IsDir(), info.Name()))
		if !info.IsDir() {
			if strings.Contains(info.Name(), ".jpeg") || strings.Contains(info.Name(), ".png") || strings.Contains(info.Name(), ".jpg") {
				imagestorage.logger.Debug(fmt.Sprintf("find file, name is: %s", info.Name()))
				err := imagestorage.DeleteCategoryImage(id, info.Name())
				if err != nil {
					imagestorage.logger.Error(fmt.Sprintf("error on delete category image: %v", err))
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		imagestorage.logger.Error(fmt.Sprintf("error on filepath.Walk: %v", err))
		return err
	}
	imagestorage.logger.Debug(fmt.Sprintf("images of category with id %s deleted success", id))
	return nil
}

func (imagestorage *OnDiskLocalStorage) GetFileList() ([]FileInStorageInfo, error) {
	imagestorage.logger.Debug("Enter in filestorage GetFileList()")
	result := make([]FileInStorageInfo, 0)
	err := filepath.Walk(imagestorage.path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			result = append(result, FileInStorageInfo{
				Name:       info.Name(),
				Path:       path,
				CreateDate: info.ModTime().String(),
				ModifyDate: info.ModTime().String(),
			})
		}
		return nil
	})
	if err != nil {
		imagestorage.logger.Error(fmt.Sprintf("error no filepath.Walk: %v", err))
		return nil, err
	}
	return result, nil
}
