package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"ya-GophKeeper/internal/constants/srverror"
	"ya-GophKeeper/internal/constants/urlsuff"
	"ya-GophKeeper/internal/content"
	"ya-GophKeeper/internal/server/otp"
	"ya-GophKeeper/internal/server/storage"
	"ya-GophKeeper/internal/server/transport/http/jwt"
)

type UserInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func RegistrationPOST(w http.ResponseWriter, r *http.Request, st storage.StorageRepo) {
	user, err := ReadAuthDate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err = st.AddNewUser(ctx, user.Login, user.Password)
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

func LoginWithPasswordPOST(w http.ResponseWriter, r *http.Request, st storage.StorageRepo) {
	user, err := ReadAuthDate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := context.Background()
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

func ChangePasswordPOST(w http.ResponseWriter, r *http.Request, st storage.StorageRepo) {
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
	err = st.ChangeUserPassword(context.Background(), claims.Login, string(newPass))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func LoginWithOTP_POST(w http.ResponseWriter, r *http.Request, m *otp.ManagerOTP) {
	user, err := ReadAuthDate(r)
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

func RemoveDataPOST(w http.ResponseWriter, r *http.Request, st storage.StorageRepo) {
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
	login := claims.Login
	switch dataType {
	case urlsuff.DatatypeCredential:
		err = st.RemoveCredentials(context.Background(), login, rem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case urlsuff.DatatypeCreditCard:
		err = st.RemoveCreditCards(context.Background(), login, rem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case urlsuff.DatatypeText:
		err = st.RemoveTexts(context.Background(), login, rem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case urlsuff.DatatypeFile:
		files, err := st.GetFiles(context.Background(), login, rem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		RemoveFiles(files)
		err = st.RemoveFiles(context.Background(), login, rem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, srverror.ErrIncorrectDataTpe.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func AddNewDataPOST(w http.ResponseWriter, r *http.Request, st storage.StorageRepo) {
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

		var respBody []byte
		switch dataType {
		case urlsuff.DatatypeCredential:
			var newCreds []content.CredentialInfo
			err := json.NewDecoder(r.Body).Decode(&newCreds)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			creds, err := st.AddNewCredentials(context.Background(), login, newCreds)
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
			cards, err := st.AddNewCreditCards(context.Background(), login, newCards)
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
			texts, err := st.AddNewTexts(context.Background(), login, newTexts)
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
			files, err := st.AddNewFiles(context.Background(), login, newFiles)
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
	w.WriteHeader(http.StatusOK)
}

func UploadFilePOST(w http.ResponseWriter, r *http.Request, fm *storage.FileManager) {
	chunk, err := ParseChunk(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = fm.SaveChunk(chunk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

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
	case "1":
		SyncFirstStep(w, r, claims.Login, dataType, st)
	case "2":
		SyncSecondStep(w, r, claims.Login, dataType, st)
	default:
		http.Error(w, srverror.ErrIncorrectSyncStep.Error(), http.StatusBadRequest)
	}
}

func SyncFirstStep(w http.ResponseWriter, r *http.Request, login string, dataType string, st storage.StorageRepo) {
	var cliInfo map[int]time.Time
	err := json.NewDecoder(r.Body).Decode(&cliInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	srvInfo, err := st.GetModtimeWithIDs(context.Background(), login, dataType)
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
	_, err = w.Write(respBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func SyncSecondStep(w http.ResponseWriter, r *http.Request, login string, dataType string, st storage.StorageRepo) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Incorrect Content-Type. application/json required", http.StatusBadRequest)
	}
	switch dataType {
	case urlsuff.DatatypeCredential:
		var newCreds []content.CredentialInfo
		err := json.NewDecoder(r.Body).Decode(&newCreds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = st.UpdateCredentials(context.Background(), login, newCreds)
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
		err = st.UpdateCreditCards(context.Background(), login, newCards)
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
		err = st.UpdateTexts(context.Background(), login, newTexts)
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
		err = st.UpdateFiles(context.Background(), login, newFiles)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "", http.StatusNotFound)
	}
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

func ReadAuthDate(r *http.Request) (*UserInfo, error) {
	var user UserInfo
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func RemoveFiles(files []content.BinaryFileInfo) {
	for _, f := range files {
		e := os.Remove(f.FilePath)
		if e != nil {
			log.Error(e)
		}
	}
}

func ParseChunk(r *http.Request) (*storage.Chunk, error) {
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

	return &storage.Chunk{
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
