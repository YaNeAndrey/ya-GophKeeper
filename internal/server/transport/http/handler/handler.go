package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
	"ya-GophKeeper/internal/constants/srverror"
	"ya-GophKeeper/internal/constants/urlsuff"
	"ya-GophKeeper/internal/content"
	"ya-GophKeeper/internal/server/otp"
	"ya-GophKeeper/internal/server/storage/filemanager"
	"ya-GophKeeper/internal/server/transport/http/jwt"
)

type UserInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserCreator interface {
	AddUser(ctx context.Context, login string, password string) error
}

func RegistrationPOST(w http.ResponseWriter, r *http.Request, st UserCreator) {
	user, err := ReadAuthData(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = st.AddUser(ctx, user.Login, user.Password)
	if err != nil {
		if errors.Is(err, srverror.ErrLoginAlreadyTaken) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = SetToken(&w, user.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type UserAuthorizer interface {
	CheckUserPassword(ctx context.Context, login string, password string) (bool, error)
}

func LoginWithPasswordPOST(w http.ResponseWriter, r *http.Request, st UserAuthorizer) {
	user, err := ReadAuthData(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	ok, err := st.CheckUserPassword(ctx, user.Login, user.Password)
	if err != nil {
		if errors.Is(err, srverror.ErrLoginNotFound) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if ok {
		err = SetToken(&w, user.Login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

type UserChanger interface {
	ChangeUserPassword(ctx context.Context, login string, password string) error
}

func ChangePasswordPOST(w http.ResponseWriter, r *http.Request, st UserChanger) {
	claims, ok := jwt.CheckAccess(r)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	newPass, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	err = st.ChangeUserPassword(ctx, claims.Login, string(newPass))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func LoginWithOTP_POST(w http.ResponseWriter, r *http.Request, m *otp.ManagerOTP) {
	user, err := ReadAuthData(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	loginOTP, err := strconv.Atoi(user.Password)
	if err != nil {
		http.Error(w, srverror.ErrIncorrectOTP.Error(), http.StatusInternalServerError)
		return
	}
	ok := m.CheckUserOTP(user.Login, loginOTP)
	if ok {
		err = SetToken(&w, user.Login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

type DataRemover interface {
	RemoveTexts(ctx context.Context, login string, textIDs []int) error
	RemoveCreditCards(ctx context.Context, login string, creditCardIDs []int) error
	RemoveCredentials(ctx context.Context, login string, credentialIDs []int) error
}

func RemoveDataPOST(w http.ResponseWriter, r *http.Request, st DataRemover) {
	claims, ok := jwt.CheckAccess(r)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Incorrect Content-Type. application/json required", http.StatusBadRequest)
	}

	dataType := strings.ToLower(chi.URLParam(r, "Datatype"))

	var rem []int
	err := json.NewDecoder(r.Body).Decode(&rem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	login := claims.Login
	switch dataType {
	case urlsuff.DatatypeCredential:
		err = st.RemoveCredentials(ctx, login, rem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case urlsuff.DatatypeCreditCard:
		err = st.RemoveCreditCards(ctx, login, rem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case urlsuff.DatatypeText:
		err = st.RemoveTexts(ctx, login, rem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	/*case urlsuff.DatatypeFile:
	files, err := st.GetFiles(ctx, login, rem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	RemoveFiles(files)
	err = st.RemoveFiles(ctx, login, rem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	*/
	default:
		http.Error(w, srverror.ErrIncorrectDataTpe.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type FileRemover interface {
	RemoveFiles(ctx context.Context, login string, fileIDs []int) ([]string, error)
}

func RemoveFilesPOST(w http.ResponseWriter, r *http.Request, st FileRemover, fm *filemanager.FileManager) {
	claims, ok := jwt.CheckAccess(r)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Incorrect Content-Type. application/json required", http.StatusBadRequest)
	}

	var rem []int
	err := json.NewDecoder(r.Body).Decode(&rem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	login := claims.Login
	/*
		files, err := st.GetFiles(ctx, login, rem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		filesForRemoving := make([]string, len(files))
		for i, info := range files {
			filesForRemoving[i] = info.FilePath
		}

		fm.RemoveFiles(filesForRemoving)
	*/
	files, err := st.RemoveFiles(ctx, login, rem)
	fm.RemoveFiles(login, files)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type DataInserter interface {
	AddFiles(ctx context.Context, login string, files []content.BinaryFileInfo) ([]content.BinaryFileInfo, error)
	AddTexts(ctx context.Context, login string, texts []content.TextInfo) ([]content.TextInfo, error)
	AddCreditCards(ctx context.Context, login string, creditCards []content.CreditCardInfo) ([]content.CreditCardInfo, error)
	AddCredentials(ctx context.Context, login string, credentials []content.CredentialInfo) ([]content.CredentialInfo, error)
}

func AddDataPOST(w http.ResponseWriter, r *http.Request, st DataInserter) {
	claims, ok := jwt.CheckAccess(r)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Incorrect Content-Type. application/json required", http.StatusBadRequest)
	}

	dataType := strings.ToLower(chi.URLParam(r, "Datatype"))
	login := claims.Login
	if r.Body != http.NoBody {

		ctx := r.Context()
		var respBody []byte
		switch dataType {
		case urlsuff.DatatypeCredential:
			var newCreds []content.CredentialInfo
			err := json.NewDecoder(r.Body).Decode(&newCreds)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			creds, err := st.AddCredentials(ctx, login, newCreds)
			if err != nil || creds == nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			respBody, err = json.Marshal(creds)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case urlsuff.DatatypeCreditCard:
			var newCards []content.CreditCardInfo
			err := json.NewDecoder(r.Body).Decode(&newCards)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			cards, err := st.AddCreditCards(ctx, login, newCards)
			if err != nil || cards == nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			respBody, err = json.Marshal(cards)
		case urlsuff.DatatypeText:
			var newTexts []content.TextInfo
			err := json.NewDecoder(r.Body).Decode(&newTexts)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			texts, err := st.AddTexts(ctx, login, newTexts)
			if err != nil || texts == nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			respBody, err = json.Marshal(texts)
		case urlsuff.DatatypeFile:
			var newFiles []content.BinaryFileInfo
			err := json.NewDecoder(r.Body).Decode(&newFiles)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			files, err := st.AddFiles(ctx, login, newFiles)
			if err != nil || files == nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			respBody, err = json.Marshal(files)
		default:
			http.Error(w, srverror.ErrIncorrectDataTpe.Error(), http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(respBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	//w.WriteHeader(http.StatusOK)
}

/*
	func SyncDataPOST(w http.ResponseWriter, r *http.Request, st storage.StorageRepo) {
		claims, ok := jwt.CheckAccess(r)
		if !ok {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		dataType := strings.ToLower(chi.URLParam(r, "Datatype"))
		syncStep := strings.ToLower(chi.URLParam(r, "StepNumber"))

		_ = dataType
		switch syncStep {
		case urlsuff.SyncFirstStep:
			SyncFirstStep(w, r, claims.Login, dataType, st)
		case urlsuff.SyncSecondStep:
			SyncSecondStep(w, r, claims.Login, dataType, st)
		default:
			http.Error(w, srverror.ErrIncorrectSyncStep.Error(), http.StatusBadRequest)
		}
	}
*/
type DataProvider interface {
	GetCreditCards(ctx context.Context, login string, cardIDs []int) ([]content.CreditCardInfo, error)
	GetCredentials(ctx context.Context, login string, credIDs []int) ([]content.CredentialInfo, error)
	GetFiles(ctx context.Context, login string, fileIDs []int) ([]content.BinaryFileInfo, error)
	GetTexts(ctx context.Context, login string, textIDs []int) ([]content.TextInfo, error)
	GetModtimeWithIDs(ctx context.Context, login string, dataType string) (map[int]time.Time, error)
}

func SyncFirstStep(w http.ResponseWriter, r *http.Request, st DataProvider) {
	claims, ok := jwt.CheckAccess(r)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	dataType := strings.ToLower(chi.URLParam(r, "Datatype"))
	login := claims.Login

	var cliInfo map[int]time.Time
	err := json.NewDecoder(r.Body).Decode(&cliInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	srvInfo, err := st.GetModtimeWithIDs(ctx, login, dataType)
	var dataForSrv []int
	var sendToCli []int
	for id, info := range srvInfo {
		modtime, ok := cliInfo[id]
		if !ok {
			sendToCli = append(sendToCli, id)
			continue
		}
		if info.Before(modtime) {
			dataForSrv = append(dataForSrv, id)
		} else if modtime.Before(info) {
			sendToCli = append(sendToCli, id)
		}
	}
	var removeFromCli []int
	for id := range cliInfo {
		_, ok := srvInfo[id]
		if !ok {
			removeFromCli = append(removeFromCli, id)
			continue
		}
	}

	answer := struct {
		DataForSrv    []int       `json:",omitempty"`
		RemoveFromCli []int       `json:",omitempty"`
		DataForCli    interface{} `json:",omitempty"`
	}{DataForSrv: dataForSrv,
		RemoveFromCli: removeFromCli}
	switch dataType {
	case urlsuff.DatatypeCredential:
		answer.DataForCli, err = st.GetCredentials(ctx, login, sendToCli)
	case urlsuff.DatatypeCreditCard:
		answer.DataForCli, err = st.GetCreditCards(ctx, login, sendToCli)
	case urlsuff.DatatypeText:
		answer.DataForCli, err = st.GetTexts(ctx, login, sendToCli)
	case urlsuff.DatatypeFile:
		answer.DataForCli, err = st.GetFiles(ctx, login, sendToCli)
	}
	respBody, err := json.Marshal(answer)

	//log.Println(respBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(respBody)))
	_, err = w.Write(respBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type DataUpdater interface {
	UpdateFiles(ctx context.Context, login string, files []content.BinaryFileInfo) error
	UpdateTexts(ctx context.Context, login string, texts []content.TextInfo) error
	UpdateCreditCards(ctx context.Context, login string, creditCards []content.CreditCardInfo) error
	UpdateCredentials(ctx context.Context, login string, credentials []content.CredentialInfo) error
}

func SyncSecondStep(w http.ResponseWriter, r *http.Request, st DataUpdater) {
	claims, ok := jwt.CheckAccess(r)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	dataType := strings.ToLower(chi.URLParam(r, "Datatype"))
	login := claims.Login

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Incorrect Content-Type. application/json required", http.StatusBadRequest)
	}
	ctx := r.Context()
	switch dataType {
	case urlsuff.DatatypeCredential:
		var newCreds []content.CredentialInfo
		err := json.NewDecoder(r.Body).Decode(&newCreds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = st.UpdateCredentials(ctx, login, newCreds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case urlsuff.DatatypeCreditCard:
		var newCards []content.CreditCardInfo
		err := json.NewDecoder(r.Body).Decode(&newCards)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = st.UpdateCreditCards(ctx, login, newCards)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case urlsuff.DatatypeText:
		var newTexts []content.TextInfo
		err := json.NewDecoder(r.Body).Decode(&newTexts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = st.UpdateTexts(ctx, login, newTexts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case urlsuff.DatatypeFile:
		var newFiles []content.BinaryFileInfo
		err := json.NewDecoder(r.Body).Decode(&newFiles)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = st.UpdateFiles(ctx, login, newFiles)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "", http.StatusNotFound)
	}
}

type FileWorker interface {
	CheckFileHash(ctx context.Context, fileID int, MD5 string) (bool, error)
	UpdateFilePath(ctx context.Context, login string, fileID int, newFilePath string) error
	RemoveFiles(ctx context.Context, login string, fileIDs []int) ([]string, error)
}

func UploadFilePOST(w http.ResponseWriter, r *http.Request, fm *filemanager.FileManager, st FileWorker) {
	claims, ok := jwt.CheckAccess(r)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	chunk, err := ParseChunk(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	login := claims.Login
	fileName, err := fm.SaveChunk(login, chunk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if fileName != "" {
		fileHash, err := fm.GetFileHash(login, fileName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx := r.Context()
		hashOK, err := st.CheckFileHash(ctx, chunk.FileID, fileHash)
		if err == nil {
			if hashOK {
				fileDownloadLinkWithoutDNS := fmt.Sprintf("/%s/%s/%s", urlsuff.FileOperationDownload, login, fileName)
				err = st.UpdateFilePath(ctx, login, chunk.FileID, fileDownloadLinkWithoutDNS)
				if err != nil {
					http.Error(w, srverror.ErrIncorrectFileHash.Error(), http.StatusInternalServerError)
					return
				}
				w.Write([]byte(fileDownloadLinkWithoutDNS))
				//w.WriteHeader(http.StatusOK)
				return
			} else {
				//TODO: What if db problem?
				files, _ := st.RemoveFiles(ctx, login, []int{chunk.FileID})
				fm.RemoveFiles(login, files)
				http.Error(w, srverror.ErrIncorrectFileHash.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func GenerateOTP_GET(w http.ResponseWriter, r *http.Request, m *otp.ManagerOTP) {
	claims, ok := jwt.CheckAccess(r)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	login := claims.Login

	loginOTP, err := m.GenerateOTP(login)
	if err != nil {
		return
	}

	_, err = w.Write([]byte(strconv.Itoa(loginOTP)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SetToken(w *http.ResponseWriter, login string) error {
	token, err := jwt.BuildJWTStringWithLogin(login)
	if err != nil {
		return err
	}
	http.SetCookie(*w, &http.Cookie{
		Name:  "token",
		Value: token,
	})
	(*w).Header().Set("Content-Type", "application/json")
	return nil
}

func ReadAuthData(r *http.Request) (*UserInfo, error) {
	var user UserInfo
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

/*
func RemoveFiles(files []content.BinaryFileInfo) {
	for _, f := range files {
		e := os.Remove(f.FilePath)
		if e != nil {
			log.Error(e)
		}
	}
}
*/

func ParseChunk(r *http.Request) (*filemanager.Chunk, error) {
	buf := new(bytes.Buffer)

	reader, err := r.MultipartReader()

	if err != nil {
		return nil, err
	}

	if err = getPart("metadata", reader, buf); err != nil {
		return nil, err
	}

	var metadata struct {
		ChunkNumber   uint64
		TotalChunks   uint64
		FileID        int
		TotalFileSize int64
	}

	err = json.Unmarshal(buf.Bytes(), &metadata)
	if err != nil {
		return nil, err
	}
	buf.Reset()

	part, err := reader.NextPart()
	if err != nil {
		return nil, err
	}

	return &filemanager.Chunk{
		ChunkNumber:   metadata.ChunkNumber,
		TotalChunks:   metadata.TotalChunks,
		FileID:        metadata.FileID,
		TotalFileSize: metadata.TotalFileSize,
		Data:          part,
	}, nil
}

func getPart(expectedPart string, reader *multipart.Reader, buf *bytes.Buffer) error {
	part, err := reader.NextPart()
	if err != nil {
		return fmt.Errorf("failed reading %s part %w", expectedPart, err)
	}

	if part.FormName() != expectedPart {
		return fmt.Errorf("invalid form name for part. Expected %s got %s", expectedPart, part.FormName())
	}

	if _, err = io.Copy(buf, part); err != nil {
		return fmt.Errorf("failed copying %s part %w", expectedPart, err)
	}

	return nil
}
