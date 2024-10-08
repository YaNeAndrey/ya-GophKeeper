package filemanager

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/google/uuid"
)

type FileManager struct {
	storageDir string
	chunkSize  int64
}

type Chunk struct {
	ChunkNumber   uint64
	TotalChunks   uint64
	FileID        int
	TotalFileSize int64
	Data          io.Reader
}

func InitFileManager(storageDir string, chunkSize int64) *FileManager {
	return &FileManager{storageDir: storageDir, chunkSize: chunkSize}
}

func (fm *FileManager) GetFileHash(subdir string, fileName string) (string, error) {
	fullFilePath := path.Join(fm.storageDir, subdir, fileName)
	return checksumMD5(fullFilePath)
}

func (fm *FileManager) RemoveFiles(subdir string, fileNames []string) {
	storagePath := path.Join(fm.storageDir, subdir)
	for _, file := range fileNames {
		fullFilePath := path.Join(storagePath, file)
		err := os.Remove(fullFilePath)
		if err != nil {
			log.Println(err)
		}
	}
}

func (fm *FileManager) SaveChunk(subdir string, chunk *Chunk) (string, error) {
	if err := os.MkdirAll(path.Join(fm.storageDir, strconv.Itoa(chunk.FileID)), 02750); err != nil {
		return "", err
	}

	if err := fm.StoreChunk(chunk); err != nil {
		return "", err
	}

	if chunk.ChunkNumber == (chunk.TotalChunks - 1) {
		newFileName := uuid.New().String()
		err := fm.buildFileFromChunks(subdir, newFileName, chunk.FileID, chunk.ChunkNumber)
		if err != nil {
			return "", err
		}

		err = os.RemoveAll(path.Join(fm.storageDir, strconv.Itoa(chunk.FileID)))
		if err != nil {
			return "", err
		}

		return newFileName, nil
	}

	return "", nil
}

func (fm *FileManager) StoreChunk(chunk *Chunk) error {
	chunkFile, err := os.Create(path.Join(fm.storageDir, strconv.Itoa(chunk.FileID), strconv.FormatUint(chunk.ChunkNumber, 10)))
	if err != nil {
		return err
	}

	defer chunkFile.Close()
	if _, err = io.CopyN(chunkFile, chunk.Data, fm.chunkSize); err != nil && err != io.EOF {
		return err
	}
	return nil
}

func (fm *FileManager) buildFileFromChunks(subdir string, fileName string, fileID int, maxChunkNumber uint64) error {
	if err := checkChunks(path.Join(fm.storageDir, strconv.Itoa(fileID)), maxChunkNumber); err != nil {
		return err
	}
	storageDir := path.Join(fm.storageDir, subdir)
	if err := os.MkdirAll(storageDir, 02750); err != nil {
		return err
	}

	fullFile, err := os.OpenFile( /*fmt.Sprintf("%s\\%s", storageDir, fileName)*/ path.Join(storageDir, fileName), os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer fullFile.Close()

	for i := uint64(0); i <= maxChunkNumber; i++ {
		err = appendChunk(path.Join(fm.storageDir, strconv.Itoa(fileID)), strconv.FormatUint(i, 10), fullFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func appendChunk(uploadDir string, chunkFile string, fullFile *os.File) error {
	src, err := os.Open(path.Join(uploadDir, chunkFile))
	if err != nil {
		return err
	}
	defer src.Close()
	if _, err := io.Copy(fullFile, src); err != nil {
		return err
	}

	return nil
}

func checkChunks(dir string, maxChunkNumber uint64) error {
	for i := uint64(0); i <= maxChunkNumber; i++ {
		if _, err := os.Stat(path.Join(dir, strconv.FormatUint(i, 10))); err != nil {
			return fmt.Errorf("%w : %d", err, i)
		}
	}
	return nil
}

func checksumMD5(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
