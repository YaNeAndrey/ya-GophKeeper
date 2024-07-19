package http

import (
	"net/http"
	"slices"
	"strings"
	"ya-GophKeeper/internal/server/gzip"
)

func Gzip() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ow := w
			//allAcceptEncodingHeaders := strings.Split(r.Header.Values("Accept-Encoding")[0], ", ")
			var allAcceptEncodingSlice []string
			allAcceptEncodingHeaders := r.Header.Values("Accept-Encoding")
			if len(allAcceptEncodingHeaders) > 0 {
				allAcceptEncodingSlice = strings.Split(allAcceptEncodingHeaders[0], ", ")
			}
			if slices.Contains(allAcceptEncodingSlice, "gzip") {
				cw := gzip.NewCompressWriter(w)
				ow = cw
				defer cw.Close()
			}

			contentEncodings := r.Header.Values("Content-Encoding")
			if slices.Contains(contentEncodings, "gzip") {
				cr, err := gzip.NewCompressReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = cr
				defer cr.Close()
			}
			next.ServeHTTP(ow, r)
		}
		return http.HandlerFunc(fn)
	}
}
