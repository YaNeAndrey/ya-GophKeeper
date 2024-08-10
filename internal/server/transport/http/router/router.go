package router

import (
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"ya-GophKeeper/internal/server/otp"
	"ya-GophKeeper/internal/server/storage"
	"ya-GophKeeper/internal/server/transport/http/handler"
)

func InitRouter(st storage.StorageRepo, m *otp.ManagerOTP) http.Handler {
	r := chi.NewRouter()
	r.NotFound(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	})

	//fs := http.FileServer(http.Dir("temp file dir"))

	logger := log.New()
	logger.SetLevel(log.InfoLevel)

	r.Route("/", func(r chi.Router) {
		r.Post("/registration", func(rw http.ResponseWriter, req *http.Request) {
			handler.RegistrationPOST(rw, req, st)
		})
		r.Route("/login", func(r chi.Router) {
			r.Post("/otp", func(rw http.ResponseWriter, req *http.Request) {
				handler.LoginWithOTP_POST(rw, req, m)
			})
			r.Post("/passwd", func(rw http.ResponseWriter, req *http.Request) {
				handler.LoginWithPasswordPOST(rw, req, st)
			})
		})

		r.Get("/otp", func(rw http.ResponseWriter, req *http.Request) {
			handler.GenerateOTP_GET(rw, req, m)
		})

		r.Post("/remove/{Datatype}", func(rw http.ResponseWriter, req *http.Request) {
			handler.RemoveDataPOST(rw, req, st)
		})

		r.Post("/add/{Datatype}", func(rw http.ResponseWriter, req *http.Request) {
			handler.AddNewDataPOST(rw, req, st)
		})

		r.Post("/sync/{Datatype}/{StepNumber}", func(rw http.ResponseWriter, req *http.Request) {
			handler.SyncDataPOST(rw, req, st)
		})

		//r.Handle("/files/", fs)
	})
	return r
}
