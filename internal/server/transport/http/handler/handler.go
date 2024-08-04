package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"strings"
	"ya-GophKeeper/internal/constants/urlsuff"
	"ya-GophKeeper/internal/content"
	"ya-GophKeeper/internal/server/otp"
	"ya-GophKeeper/internal/server/srverror"
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
		http.Error(w, err.Error(), http.StatusUnauthorized)
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

func LoginWithOTP_POST(w http.ResponseWriter, r *http.Request, m *otp.ManagerOTP) {
	user, err := ReadAuthDate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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
		}
	case urlsuff.DatatypeCreditCard:
		err = st.RemoveCreditCards(context.Background(), login, rem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case urlsuff.DatatypeText:
		err = st.RemoveTexts(context.Background(), login, rem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case urlsuff.DatatypeFile:
		files, err := st.RemoveFiles(context.Background(), login, rem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		RemoveFiles(files)
	default:
		http.Error(w, srverror.ErrIncorrectDataTpe.Error(), http.StatusBadRequest)
	}
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
	switch dataType {
	case urlsuff.DatatypeCredential:
		var newCreds []content.CredentialInfo
		err := json.NewDecoder(r.Body).Decode(&newCreds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = st.AddNewCredentials(context.Background(), login, newCreds)
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
		err = st.AddNewCreditCards(context.Background(), login, newCards)
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
		err = st.AddNewTexts(context.Background(), login, newTexts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case urlsuff.DatatypeFile:
		AddNewFilePOST()
	default:
		http.Error(w, srverror.ErrIncorrectDataTpe.Error(), http.StatusBadRequest)
	}
}

func AddNewFilePOST() {
}

func SyncDataPOST(w http.ResponseWriter, r *http.Request, st storage.StorageRepo) {
	claims, ok := jwt.CheckAccess(r)
	_ = claims
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Incorrect Content-Type. application/json required", http.StatusBadRequest)
	}

	dataType := strings.ToLower(chi.URLParam(r, "Datatype"))
	syncStep := strings.ToLower(chi.URLParam(r, "StepNumber"))

	_ = dataType
	switch syncStep {
	case "1":
		//FirstStep()
	case "2":
		//SecondStep()
	default:
		http.Error(w, srverror.ErrIncorrectSyncStep.Error(), http.StatusBadRequest)
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

	_ = loginOTP
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

func RemoveFiles(files []string) {
	for _, f := range files {
		e := os.Remove(f)
		if e != nil {
			log.Error(e)
		}
	}
}
