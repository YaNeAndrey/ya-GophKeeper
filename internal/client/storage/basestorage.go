package storage

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"math"
	"os"
	"path"
	"time"
	"ya-GophKeeper/internal/constants/clerror"
	"ya-GophKeeper/internal/content"
)

type BaseStorage struct {
	credentials Credentials
	creditCards CreditCards
	texts       Texts
	files       Files
}

func NewBaseStorage(tempDir string) *BaseStorage {
	return &BaseStorage{
		files: Files{tempDir: tempDir},
	}
}

func (st *BaseStorage) AddNewCreditCard(creditCard *content.CreditCardInfo) error {
	if creditCard.CardNumber == "" || creditCard.CVV == "" || creditCard.Bank == "" || creditCard.ValidThru.IsZero() {
		return fmt.Errorf("AddNewCreditCard : %w", clerror.ErrAllFieldsMustBeFulled)
	} else {
		st.creditCards.stored = append(st.creditCards.stored, *creditCard)
		return nil
	}
}
func (st *BaseStorage) AddNewCredential(credential *content.CredentialInfo) error {
	if credential.Login == "" || credential.Resource == "" || credential.Password == "" {
		return fmt.Errorf("AddNewCredential : %w", clerror.ErrAllFieldsMustBeFulled)
	} else {
		st.credentials.stored = append(st.credentials.stored, *credential)
		return nil
	}
}
func (st *BaseStorage) AddNewText(text *content.TextInfo) error {
	if text.Content == "" {
		return fmt.Errorf("AddNewText : %w", clerror.ErrAllFieldsMustBeFulled)
	} else {
		st.texts.stored = append(st.texts.stored, *text)
		return nil
	}
}
func (st *BaseStorage) AddNewFile(file *content.BinaryFileInfo) error {
	if file.FileName == "" || file.FilePath == "" {
		return fmt.Errorf("AddNewFile : %w", clerror.ErrAllFieldsMustBeFulled)
	}
	tempFilePath := path.Join(st.files.tempDir, uuid.New().String())
	size, err := fileSize(file.FilePath)
	if err != nil {
		return err
	}
	if size > math.MaxUint32 {
		return clerror.ErrMaxFileSizeExceeded
	}
	err = copyFileContents(file.FilePath, tempFilePath)
	if err != nil {
		return clerror.ErrCopyFileProblem
	}

	newFileData := content.BinaryFileInfo{
		ID:               0,
		FileName:         file.FileName,
		FilePath:         tempFilePath,
		Description:      file.Description,
		ModificationTime: time.Now(),
	}
	st.files.stored = append(st.files.stored, newFileData)
	return nil
}

func (st *BaseStorage) RemoveCreditCard(index int) error {
	if st.creditCards.stored == nil {
		return clerror.ErrOutOfRange
	}
	if index > len(st.creditCards.stored) || index < 0 {
		return clerror.ErrOutOfRange
	}
	cardID := st.creditCards.stored[index].ID
	if cardID != 0 {
		st.creditCards.removed = append(st.creditCards.removed, cardID)
	}
	st.creditCards.stored = append(st.creditCards.stored[:index], st.creditCards.stored[index+1:]...)
	return nil
}
func (st *BaseStorage) RemoveCredential(index int) error {
	if st.credentials.stored == nil {
		return clerror.ErrOutOfRange
	}
	if index > len(st.credentials.stored) || index < 0 {
		return clerror.ErrOutOfRange
	}

	credID := st.credentials.stored[index].ID
	if credID != 0 {
		st.credentials.removed = append(st.credentials.removed, credID)
	}
	st.credentials.stored = append(st.credentials.stored[:index], st.credentials.stored[index+1:]...)
	return nil
}
func (st *BaseStorage) RemoveText(index int) error {
	if st.texts.stored == nil {
		return clerror.ErrOutOfRange
	}
	if index > len(st.texts.stored) || index < 0 {
		return clerror.ErrOutOfRange
	}

	textID := st.texts.stored[index].ID
	if textID != 0 {
		st.texts.removed = append(st.texts.removed, textID)
	}
	st.texts.stored = append(st.texts.stored[:index], st.texts.stored[index+1:]...)
	return nil
}
func (st *BaseStorage) RemoveFile(index int) error {
	if st.files.stored == nil {
		return clerror.ErrOutOfRange
	}
	if index > len(st.files.stored) || index < 0 {
		return clerror.ErrOutOfRange
	}

	fileID := st.files.stored[index].ID
	if fileID != 0 {
		st.files.removed = append(st.files.removed, fileID)
	}
	err := os.Remove(st.files.stored[index].FilePath)
	if err != nil {
		log.Println(err)
	}
	st.files.stored = append(st.files.stored[:index], st.files.stored[index+1:]...)
	return nil
}

func (st *BaseStorage) UpdateCreditCards(index int, creditCard *content.CreditCardInfo) error {
	if index > len(st.creditCards.stored) || index < 0 {
		return clerror.ErrOutOfRange
	}
	if creditCard.CardNumber != "" {
		st.creditCards.stored[index].CardNumber = creditCard.CardNumber
	}
	if creditCard.Bank != "" {
		st.creditCards.stored[index].Bank = creditCard.Bank
	}
	if creditCard.CVV != "" {
		st.creditCards.stored[index].CVV = creditCard.CVV
	}
	if !creditCard.ValidThru.IsZero() {
		st.creditCards.stored[index].ValidThru = creditCard.ValidThru
	}
	st.creditCards.stored[index].ModificationTime = creditCard.ModificationTime
	return nil
}
func (st *BaseStorage) UpdateCredentials(index int, credential *content.CredentialInfo) error {
	if index > len(st.credentials.stored) || index < 0 {
		return clerror.ErrOutOfRange
	}
	if credential.Login != "" {
		st.credentials.stored[index].Login = credential.Login
	}
	if credential.Password != "" {
		st.credentials.stored[index].Password = credential.Password
	}
	if credential.Resource != "" {
		st.credentials.stored[index].Resource = credential.Resource
	}
	st.credentials.stored[index].ModificationTime = credential.ModificationTime
	return nil
}
func (st *BaseStorage) UpdateFiles(index int, file *content.BinaryFileInfo) error {
	if index > len(st.files.stored) || index < 0 {
		return clerror.ErrOutOfRange
	}
	tempFilePath := path.Join(st.files.tempDir, uuid.New().String())
	size, err := fileSize(file.FilePath)
	if err != nil {
		return err
	}
	if size > math.MaxUint32 {
		return clerror.ErrMaxFileSizeExceeded
	}
	err = copyFileContents(file.FilePath, tempFilePath)
	if err != nil {
		return clerror.ErrCopyFileProblem
	}
	if file.FileName != "" {
		st.files.stored[index].FileName = file.FileName
	}
	if file.Description != "" {
		st.files.stored[index].Description = file.Description
	}
	st.files.stored[index].ModificationTime = file.ModificationTime
	return nil
}
func (st *BaseStorage) UpdateTexts(index int, text *content.TextInfo) error {
	if index > len(st.texts.stored) || index < 0 {
		return clerror.ErrOutOfRange
	}
	if text.Content != "" {
		st.texts.stored[index].Content = text.Content
	}
	if text.Description != "" {
		st.texts.stored[index].Description = text.Description
	}
	st.texts.stored[index].ModificationTime = text.ModificationTime
	return nil
}

func (st *BaseStorage) GetCreditCardsData() *CreditCards {
	return &st.creditCards
}
func (st *BaseStorage) GetCredentialsData() *Credentials {
	return &st.credentials
}
func (st *BaseStorage) GetFilesData() *Files {
	return &st.files
}
func (st *BaseStorage) GetTextsData() *Texts {
	return &st.texts
}

func (st *BaseStorage) Clear() {
	st.files.Clear()
	st.texts.Clear()
	st.credentials.Clear()
	st.creditCards.Clear()
}

func fileSize(filename string) (int64, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	if fi.IsDir() {
		return 0, clerror.ErrPathIsDir
	}
	return fi.Size(), nil
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
