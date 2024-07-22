package storage

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
	"path"
	"time"
	"ya-GophKeeper/internal/constants/consterror"
	"ya-GophKeeper/internal/content"
)

type BaseStorage struct {
	tempDir     string
	credentials []content.CredentialInfo
	creditCards []content.CreditCardInfo
	texts       []content.TextInfo
	files       []content.BinaryFileInfo
}

func NewBaseStorage(tempDir string) *BaseStorage {
	return &BaseStorage{
		tempDir: tempDir,
	}
}

func (st *BaseStorage) AddNewCreditCard(creditCard *content.CreditCardInfo) error {
	if creditCard.CardNumber == "" || creditCard.CVV == "" || creditCard.Bank == "" || creditCard.ValidThru.IsZero() {
		return fmt.Errorf("AddNewCreditCard : %w", consterror.ErrAllFieldsMustBeFulled)
	} else {
		st.creditCards = append(st.creditCards, *creditCard)
		return nil
	}
}
func (st *BaseStorage) AddNewCredential(credential *content.CredentialInfo) error {
	if credential.Login == "" || credential.Resource == "" || credential.Password == "" {
		return fmt.Errorf("AddNewCredential : %w", consterror.ErrAllFieldsMustBeFulled)
	} else {
		st.credentials = append(st.credentials, *credential)
		return nil
	}
}
func (st *BaseStorage) AddNewFile(file *content.BinaryFileInfo) error {
	if file.BaseFileName == "" || file.FilePath == "" {
		return fmt.Errorf("AddNewFile : %w", consterror.ErrAllFieldsMustBeFulled)
	}
	tempFilePath := path.Join(st.tempDir, uuid.New().String())
	if fileExists(file.FilePath) {
		err := copyFileContents(file.FilePath, tempFilePath)
		if err != nil {
			return consterror.ErrCopyFileProblem
		}
	} else {
		return consterror.ErrFileNotFound
	}

	newFileData := content.BinaryFileInfo{
		ID:               0,
		BaseFileName:     file.BaseFileName,
		FilePath:         tempFilePath,
		Description:      file.Description,
		ModificationTime: time.Now(),
	}
	st.files = append(st.files, newFileData)
	return nil
}
func (st *BaseStorage) AddNewText(text *content.TextInfo) error {
	if text.Content == "" {
		return fmt.Errorf("AddNewText : %w", consterror.ErrAllFieldsMustBeFulled)
	} else {
		st.texts = append(st.texts, *text)
		return nil
	}
}

func (st *BaseStorage) RemoveCreditCard(index int) error {
	if st.creditCards == nil {
		return consterror.ErrOutOfRange
	}
	if index > len(st.creditCards) || index < 0 {
		return consterror.ErrOutOfRange
	}
	st.creditCards = append(st.creditCards[:index], st.creditCards[index+1:]...)
	return nil
}
func (st *BaseStorage) RemoveCredential(index int) error {
	if st.credentials == nil {
		return consterror.ErrOutOfRange
	}
	if index > len(st.credentials) || index < 0 {
		return consterror.ErrOutOfRange
	}
	st.credentials = append(st.credentials[:index], st.credentials[index+1:]...)
	return nil
}
func (st *BaseStorage) RemoveText(index int) error {
	if st.texts == nil {
		return consterror.ErrOutOfRange
	}
	if index > len(st.texts) || index < 0 {
		return consterror.ErrOutOfRange
	}
	st.texts = append(st.texts[:index], st.texts[index+1:]...)
	return nil
}
func (st *BaseStorage) RemoveFile(index int) error {
	if st.files == nil {
		return consterror.ErrOutOfRange
	}
	if index > len(st.files) || index < 0 {
		return consterror.ErrOutOfRange
	}
	st.files = append(st.files[:index], st.files[index+1:]...)
	return nil
}

func (st *BaseStorage) UpdateCreditCards(index int, creditCard *content.CreditCardInfo) error {
	if index > len(st.creditCards) || index < 0 {
		return consterror.ErrOutOfRange
	}
	if creditCard.CardNumber != "" {
		st.creditCards[index].CardNumber = creditCard.CardNumber
	}
	if creditCard.Bank != "" {
		st.creditCards[index].Bank = creditCard.Bank
	}
	if creditCard.CVV != "" {
		st.creditCards[index].CVV = creditCard.CVV
	}
	if !creditCard.ValidThru.IsZero() {
		st.creditCards[index].ValidThru = creditCard.ValidThru
	}
	st.creditCards[index].ModificationTime = creditCard.ModificationTime
	return nil
}
func (st *BaseStorage) UpdateCredentials(index int, credential *content.CredentialInfo) error {
	if index > len(st.credentials) || index < 0 {
		return consterror.ErrOutOfRange
	}
	if credential.Login != "" {
		st.credentials[index].Login = credential.Login
	}
	if credential.Password != "" {
		st.credentials[index].Password = credential.Password
	}
	if credential.Resource != "" {
		st.credentials[index].Resource = credential.Resource
	}
	st.credentials[index].ModificationTime = credential.ModificationTime
	return nil
}
func (st *BaseStorage) UpdateFiles(index int, file *content.BinaryFileInfo) error {
	if index > len(st.files) || index < 0 {
		return consterror.ErrOutOfRange
	}
	if file.FilePath != "" {
		if fileExists(file.FilePath) {
			err := copyFileContents(file.FilePath, st.files[index].FilePath)
			if err != nil {
				return err
			}
		}
	}
	if file.BaseFileName != "" {
		st.files[index].BaseFileName = file.BaseFileName
	}
	if file.Description != "" {
		st.files[index].Description = file.Description
	}
	st.files[index].ModificationTime = file.ModificationTime
	return nil
}
func (st *BaseStorage) UpdateTexts(index int, text *content.TextInfo) error {
	if index > len(st.texts) || index < 0 {
		return consterror.ErrOutOfRange
	}
	if text.Content != "" {
		st.texts[index].Content = text.Content
	}
	if text.Description != "" {
		st.texts[index].Description = text.Description
	}
	st.texts[index].ModificationTime = text.ModificationTime
	return nil
}

func (st *BaseStorage) GetCreditCardData() []content.CreditCardInfo {
	return st.creditCards
}
func (st *BaseStorage) GetCredentials() []content.CredentialInfo {
	return st.credentials
}
func (st *BaseStorage) GetFilesData() []content.BinaryFileInfo {
	return st.files
}
func (st *BaseStorage) GetTextData() []content.TextInfo {
	return st.texts
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func copyFileContents(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return err
}
