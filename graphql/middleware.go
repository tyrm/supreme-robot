package graphql

import (
	"net/http"
	"time"
)

const layoutCombined = "02/Jan/2006:15:04:05 -0700"

type responseWriterX struct {
	http.ResponseWriter
	status     int
	bodyLength int
}

func (w *responseWriterX) Write(b []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(b)
	w.bodyLength += n
	return
}

func (w *responseWriterX) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.status = status
	return
}

func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wx := &responseWriterX{
			ResponseWriter: w,
			status:         200,
			bodyLength:     0,
		}

		// everything from here is json
		w.Header().Set("Content-Type", "application/json")

		// Do Request
		next.ServeHTTP(wx, r)

		// Log It
		duration := time.Since(start)
		logger.Debugf("%s - %s [%s] \"%s %s %s\" %d %d \"%s\" \"%s\" rt=%d",
			r.RemoteAddr,
			"-",
			start.Format(layoutCombined),
			r.Method,
			r.URL.Path,
			r.Proto,
			wx.status,
			wx.bodyLength,
			r.Referer(),
			r.UserAgent(),
			duration.Milliseconds(),
		)
	})
}
