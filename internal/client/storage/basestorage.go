package storage

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"time"
	"ya-GophKeeper/internal/client/clerror"
	"ya-GophKeeper/internal/content"
)

type BaseStorage struct {
	tempDir            string
	credentials        []content.CredentialInfo
	creditCards        []content.CreditCardInfo
	texts              []content.TextInfo
	files              []content.BinaryFileInfo
	contentForRemoving RemovedContent
}

type RemovedContent struct {
	credentials []int
	creditCards []int
	texts       []int
	files       []int
}

func NewBaseStorage(tempDir string) *BaseStorage {
	return &BaseStorage{
		tempDir: tempDir,
	}
}

func (st *BaseStorage) AddNewCreditCard(creditCard *content.CreditCardInfo) error {
	if creditCard.CardNumber == "" || creditCard.CVV == "" || creditCard.Bank == "" || creditCard.ValidThru.IsZero() {
		return fmt.Errorf("AddNewCreditCard : %w", clerror.ErrAllFieldsMustBeFulled)
	} else {
		st.creditCards = append(st.creditCards, *creditCard)
		return nil
	}
}
func (st *BaseStorage) AddNewCredential(credential *content.CredentialInfo) error {
	if credential.Login == "" || credential.Resource == "" || credential.Password == "" {
		return fmt.Errorf("AddNewCredential : %w", clerror.ErrAllFieldsMustBeFulled)
	} else {
		st.credentials = append(st.credentials, *credential)
		return nil
	}
}
func (st *BaseStorage) AddNewFile(file *content.BinaryFileInfo) error {
	if file.FileName == "" || file.FilePath == "" {
		return fmt.Errorf("AddNewFile : %w", clerror.ErrAllFieldsMustBeFulled)
	}
	tempFilePath := path.Join(st.tempDir, uuid.New().String())
	if fileExists(file.FilePath) {
		err := copyFileContents(file.FilePath, tempFilePath)
		if err != nil {
			return clerror.ErrCopyFileProblem
		}
	} else {
		return clerror.ErrFileNotFound
	}

	newFileData := content.BinaryFileInfo{
		ID:               0,
		FileName:         file.FileName,
		FilePath:         tempFilePath,
		Description:      file.Description,
		ModificationTime: time.Now(),
	}
	st.files = append(st.files, newFileData)
	return nil
}
func (st *BaseStorage) AddNewText(text *content.TextInfo) error {
	if text.Content == "" {
		return fmt.Errorf("AddNewText : %w", clerror.ErrAllFieldsMustBeFulled)
	} else {
		st.texts = append(st.texts, *text)
		return nil
	}
}

func (st *BaseStorage) RemoveCreditCard(index int) error {
	if st.creditCards == nil {
		return clerror.ErrOutOfRange
	}
	if index > len(st.creditCards) || index < 0 {
		return clerror.ErrOutOfRange
	}
	cardID := st.creditCards[index].ID
	if cardID != 0 {
		st.contentForRemoving.creditCards = append(st.contentForRemoving.creditCards, cardID)
	}
	st.creditCards = append(st.creditCards[:index], st.creditCards[index+1:]...)
	return nil
}
func (st *BaseStorage) RemoveCredential(index int) error {
	if st.credentials == nil {
		return clerror.ErrOutOfRange
	}
	if index > len(st.credentials) || index < 0 {
		return clerror.ErrOutOfRange
	}

	credID := st.credentials[index].ID
	if credID != 0 {
		st.contentForRemoving.credentials = append(st.contentForRemoving.credentials, credID)
	}
	st.credentials = append(st.credentials[:index], st.credentials[index+1:]...)
	return nil
}
func (st *BaseStorage) RemoveText(index int) error {
	if st.texts == nil {
		return clerror.ErrOutOfRange
	}
	if index > len(st.texts) || index < 0 {
		return clerror.ErrOutOfRange
	}

	textID := st.texts[index].ID
	if textID != 0 {
		st.contentForRemoving.credentials = append(st.contentForRemoving.texts, textID)
	}
	st.texts = append(st.texts[:index], st.texts[index+1:]...)
	return nil
}
func (st *BaseStorage) RemoveFile(index int) error {
	if st.files == nil {
		return clerror.ErrOutOfRange
	}
	if index > len(st.files) || index < 0 {
		return clerror.ErrOutOfRange
	}

	fileID := st.files[index].ID
	if fileID != 0 {
		st.contentForRemoving.files = append(st.contentForRemoving.files, fileID)
	}
	err := os.Remove(st.files[index].FilePath)
	if err != nil {
		log.Println(err)
	}
	st.files = append(st.files[:index], st.files[index+1:]...)
	return nil
}

func (st *BaseStorage) UpdateCreditCards(index int, creditCard *content.CreditCardInfo) error {
	if index > len(st.creditCards) || index < 0 {
		return clerror.ErrOutOfRange
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
		return clerror.ErrOutOfRange
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
		return clerror.ErrOutOfRange
	}
	if file.FilePath != "" {
		if fileExists(file.FilePath) {
			err := copyFileContents(file.FilePath, st.files[index].FilePath)
			if err != nil {
				return err
			}
		}
	}
	if file.FileName != "" {
		st.files[index].FileName = file.FileName
	}
	if file.Description != "" {
		st.files[index].Description = file.Description
	}
	st.files[index].ModificationTime = file.ModificationTime
	return nil
}
func (st *BaseStorage) UpdateTexts(index int, text *content.TextInfo) error {
	if index > len(st.texts) || index < 0 {
		return clerror.ErrOutOfRange
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
