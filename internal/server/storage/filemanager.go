package storage

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
)

type FileManager struct {
	storageDir string
}

type Chunk struct {
	ChunkNumber   uint64
	TotalChunks   uint64
	FileID        int
	TotalFileSize int64
	Data          io.Reader
}

func InitFileManager(storageDir string) *FileManager {
	return &FileManager{storageDir: storageDir}
}

func (fm *FileManager) SaveChunk(chunk *Chunk) error {
	if err := os.MkdirAll(path.Join(fm.storageDir, strconv.Itoa(chunk.FileID)), 02750); err != nil {
		return err
	}

	if err := fm.StoreChunk(chunk); err != nil {
		return err
	}

	if chunk.ChunkNumber == (chunk.TotalChunks - 1) {
		err := fm.buildFileFromChunks(strconv.Itoa(chunk.FileID), chunk.ChunkNumber)
		if err != nil {
			return err
		}

		err = os.RemoveAll(path.Join(fm.storageDir, strconv.Itoa(chunk.FileID)))
		if err != nil {
			return err
		}
	}

	return nil
}

func (fm *FileManager) StoreChunk(chunk *Chunk) error {
	chunkFile, err := os.Create(path.Join(fm.storageDir, strconv.Itoa(chunk.FileID), strconv.FormatUint(chunk.ChunkNumber, 10)))
	if err != nil {
		return err
	}

	if _, err = io.CopyN(chunkFile, chunk.Data, 5*1024*1024); err != nil && err != io.EOF {
		return err
	}
	return nil
}

func (fm *FileManager) buildFileFromChunks(fileName string, maxChunkNumber uint64) error {
	if err := checkChunks(path.Join(fm.storageDir, fileName), maxChunkNumber); err != nil {
		return err
	}
	fullFile, err := os.OpenFile(fmt.Sprintf("%s\\%s.file", fm.storageDir, fileName), os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	for i := uint64(0); i <= maxChunkNumber; i++ {
		err = appendChunk(path.Join(fm.storageDir, fileName), strconv.FormatUint(i, 10), fullFile)
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
