package memory

import (
	log "github.com/sirupsen/logrus"
	"os"
	"ya-GophKeeper/internal/client/storage/collection"
	"ya-GophKeeper/internal/constants/clerror"
	"ya-GophKeeper/internal/content"
)

type BaseStorage struct {
	credentials collection.Credentials
	creditCards collection.CreditCards
	texts       collection.Texts
	files       collection.Files
}

func NewBaseStorage() *BaseStorage {
	return &BaseStorage{}
}

func (st *BaseStorage) AddCreditCard(creditCard *content.CreditCardInfo) {
	st.creditCards.Stored = append(st.creditCards.Stored, *creditCard)
}
func (st *BaseStorage) AddCredential(credential *content.CredentialInfo) {
	st.credentials.Stored = append(st.credentials.Stored, *credential)
}
func (st *BaseStorage) AddText(text *content.TextInfo) {
	st.texts.Stored = append(st.texts.Stored, *text)
}
func (st *BaseStorage) AddFile(file *content.BinaryFileInfo) {
	st.files.Stored = append(st.files.Stored, *file)
	/*
		tempFilePath := path.Join(uuid.New().String())
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
		st.files.Stored = append(st.files.Stored, newFileData)
		return nil
	*/
}

func (st *BaseStorage) RemoveCreditCard(index int) error {
	if st.creditCards.Stored == nil {
		return clerror.ErrOutOfRange
	}
	if index >= len(st.creditCards.Stored) || index < 0 {
		return clerror.ErrOutOfRange
	}
	cardID := st.creditCards.Stored[index].ID
	if cardID != 0 {
		st.creditCards.Removed = append(st.creditCards.Removed, cardID)
	}
	st.creditCards.Stored = append(st.creditCards.Stored[:index], st.creditCards.Stored[index+1:]...)
	return nil
}
func (st *BaseStorage) RemoveCredential(index int) error {
	if st.credentials.Stored == nil {
		return clerror.ErrOutOfRange
	}
	if index >= len(st.credentials.Stored) || index < 0 {
		return clerror.ErrOutOfRange
	}

	credID := st.credentials.Stored[index].ID
	if credID != 0 {
		st.credentials.Removed = append(st.credentials.Removed, credID)
	}
	st.credentials.Stored = append(st.credentials.Stored[:index], st.credentials.Stored[index+1:]...)
	return nil
}
func (st *BaseStorage) RemoveText(index int) error {
	if st.texts.Stored == nil {
		return clerror.ErrOutOfRange
	}
	if index >= len(st.texts.Stored) || index < 0 {
		return clerror.ErrOutOfRange
	}

	textID := st.texts.Stored[index].ID
	if textID != 0 {
		st.texts.Removed = append(st.texts.Removed, textID)
	}
	st.texts.Stored = append(st.texts.Stored[:index], st.texts.Stored[index+1:]...)
	return nil
}
func (st *BaseStorage) RemoveFile(index int) error {
	if st.files.Stored == nil {
		return clerror.ErrOutOfRange
	}
	if index >= len(st.files.Stored) || index < 0 {
		return clerror.ErrOutOfRange
	}

	fileID := st.files.Stored[index].ID
	if fileID != 0 {
		st.files.Removed = append(st.files.Removed, fileID)
	}
	err := os.Remove(st.files.Stored[index].FilePath)
	if err != nil {
		log.Println(err)
	}
	st.files.Stored = append(st.files.Stored[:index], st.files.Stored[index+1:]...)
	return nil
}

func (st *BaseStorage) UpdateCreditCards(index int, creditCard *content.CreditCardInfo) error {
	if index > len(st.creditCards.Stored) || index < 0 {
		return clerror.ErrOutOfRange
	}
	if creditCard.CardNumber != "" {
		st.creditCards.Stored[index].CardNumber = creditCard.CardNumber
	}
	if creditCard.Bank != "" {
		st.creditCards.Stored[index].Bank = creditCard.Bank
	}
	if creditCard.CVV != "" {
		st.creditCards.Stored[index].CVV = creditCard.CVV
	}
	if !creditCard.ValidThru.IsZero() {
		st.creditCards.Stored[index].ValidThru = creditCard.ValidThru
	}
	st.creditCards.Stored[index].ModificationTime = creditCard.ModificationTime
	return nil
}
func (st *BaseStorage) UpdateCredentials(index int, credential *content.CredentialInfo) error {
	if index > len(st.credentials.Stored) || index < 0 {
		return clerror.ErrOutOfRange
	}
	if credential.Login != "" {
		st.credentials.Stored[index].Login = credential.Login
	}
	if credential.Password != "" {
		st.credentials.Stored[index].Password = credential.Password
	}
	if credential.Resource != "" {
		st.credentials.Stored[index].Resource = credential.Resource
	}
	st.credentials.Stored[index].ModificationTime = credential.ModificationTime
	return nil
}
func (st *BaseStorage) UpdateFiles(index int, file *content.BinaryFileInfo) error {
	if index > len(st.files.Stored) || index < 0 {
		return clerror.ErrOutOfRange
	}
	/*
			tempFilePath := path.Join(st.files.TempDir, uuid.New().String())
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
	*/
	if file.FileName != "" {
		st.files.Stored[index].FileName = file.FileName
	}
	if file.Description != "" {
		st.files.Stored[index].Description = file.Description
	}
	st.files.Stored[index].ModificationTime = file.ModificationTime
	return nil
}
func (st *BaseStorage) UpdateTexts(index int, text *content.TextInfo) error {
	if index > len(st.texts.Stored) || index < 0 {
		return clerror.ErrOutOfRange
	}
	if text.Content != "" {
		st.texts.Stored[index].Content = text.Content
	}
	if text.Description != "" {
		st.texts.Stored[index].Description = text.Description
	}
	st.texts.Stored[index].ModificationTime = text.ModificationTime
	return nil
}

func (st *BaseStorage) GetCreditCardsData() *collection.CreditCards {
	return &st.creditCards
}
func (st *BaseStorage) GetCredentialsData() *collection.Credentials {
	return &st.credentials
}
func (st *BaseStorage) GetFilesData() *collection.Files {
	return &st.files
}
func (st *BaseStorage) GetTextsData() *collection.Texts {
	return &st.texts
}

func (st *BaseStorage) Clear() {
	st.files.Clear()
	st.texts.Clear()
	st.credentials.Clear()
	st.creditCards.Clear()
}

/*
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
*/
