package files

import (
	"TelegramBot/libs/e"
	"TelegramBot/storage"
	"encoding/gob"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

const defaultPerm = 0774

func (storage Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfNil("can't save the data", err) }()
	fPath := filepath.Join(storage.basePath, page.UserName)
	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}
	fName, err := fileName(page)
	if err != nil {
		return err
	}
	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}
	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfNil("can't save the data", err) }()

	path := filepath.Join(s.basePath, userName)
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(files))

	file := files[randomIndex]

	return s.decodePage(filepath.Join(path, file.Name()))
}
func (s Storage) Remove(p *storage.Page) error {

	fileName, err := fileName(p)
	if err != nil {
		return err
	}

	filePath := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}
func (s Storage) IsExist(page *storage.Page) (condition bool, err error) {
	defer func() { err = e.WrapIfNil("can't define the existence", err) }()
	fileName, err := fileName(page)
	if err != nil {
		return false, err
	}
	filePath := filepath.Join(s.basePath, page.UserName, fileName)
	switch _, err := os.Stat(filePath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, err
	}
	return true, nil

}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't open the file", err)
	}
	var storagePage storage.Page

	if err := gob.NewDecoder(file).Decode(&storagePage); err != nil {
		return nil, e.Wrap("can't decode the file", err)
	}
	return &storagePage, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
