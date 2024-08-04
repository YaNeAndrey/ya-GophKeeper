package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"ya-GophKeeper/internal/constants/urlsuff"
	"ya-GophKeeper/internal/content"
	"ya-GophKeeper/internal/server/otp"
	"ya-GophKeeper/internal/server/srverror"
	"ya-GophKeeper/internal/server/storage"
)

func RegistrationPOST(w http.ResponseWriter, r *http.Request, st storage.StorageRepo) {
}

func LoginWithPasswordPOST(w http.ResponseWriter, r *http.Request, st storage.StorageRepo) {
}

func LoginWithOTP_POST(w http.ResponseWriter, r *http.Request, m *otp.ManagerOTP) {
}

func RemoveDataPOST(w http.ResponseWriter, r *http.Request, st storage.StorageRepo) {
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

	switch dataType {
	case urlsuff.DatatypeCredential:
		err = st.RemoveCredentials(context.Background(), "", rem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case urlsuff.DatatypeCreditCard:
		err = st.RemoveCreditCards(context.Background(), "", rem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case urlsuff.DatatypeText:
		err = st.RemoveTexts(context.Background(), "", rem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case urlsuff.DatatypeFile:
		files, err := st.RemoveFiles(context.Background(), "", rem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		RemoveFiles(files)
	default:
		http.Error(w, srverror.ErrIncorrectDataTpe.Error(), http.StatusBadRequest)
	}
}

func AddNewDataPOST(w http.ResponseWriter, r *http.Request, st storage.StorageRepo) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Incorrect Content-Type. application/json required", http.StatusBadRequest)
	}

	dataType := strings.ToLower(chi.URLParam(r, "Datatype"))
	switch dataType {
	case urlsuff.DatatypeCredential:
		var newCreds []content.CredentialInfo
		err := json.NewDecoder(r.Body).Decode(&newCreds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = st.AddNewCredentials(context.Background(), "", newCreds)
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
		err = st.AddNewCreditCards(context.Background(), "", newCards)
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
		err = st.AddNewTexts(context.Background(), "", newTexts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case urlsuff.DatatypeFile:
	default:
		http.Error(w, srverror.ErrIncorrectDataTpe.Error(), http.StatusBadRequest)
	}
}

func AddNewFilePOST(w http.ResponseWriter, r *http.Request, st storage.StorageRepo) {
}

func SyncDataPOST(w http.ResponseWriter, r *http.Request, st storage.StorageRepo) {
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
}

func RemoveFiles(files []string) {
	for _, f := range files {
		e := os.Remove(f)
		if e != nil {
			log.Error(e)
		}
	}
}
