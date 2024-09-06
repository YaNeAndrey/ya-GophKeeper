package router

import (
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"ya-GophKeeper/internal/server/otp"
	"ya-GophKeeper/internal/server/storage"
	"ya-GophKeeper/internal/server/transport/http/handler"
)

func InitRouter(st storage.StorageRepo, m *otp.ManagerOTP, fm *storage.FileManager) http.Handler {
	r := chi.NewRouter()
	r.NotFound(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	})

	//fs := filesystem.NoListFileSystem{Base: http.Dir("C:\\Users\\pc\\Documents\\GoYandex\\ya-GophKeeper")}
	fs := http.FileServer(http.Dir("."))
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

		r.Post("/changepass", func(rw http.ResponseWriter, req *http.Request) {
			handler.ChangePasswordPOST(rw, req, st)
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

		r.Post("/file/upload", func(rw http.ResponseWriter, req *http.Request) {
			handler.UploadFilePOST(rw, req, fm)
		})
		/*
			r.Route("/file/", func(r chi.Router) {

					r.Post("/upload", func(rw http.ResponseWriter, req *http.Request) {
						handler.UploadFilePOST(rw, req)
					})

				r.Handle("/download/*", noDirListing(http.StripPrefix("/files/download/", fs)))
			})
		*/

		r.Handle("/download/*", noDirListing(http.StripPrefix("/files/download/", fs)))
	})
	return r
}

func noDirListing(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
