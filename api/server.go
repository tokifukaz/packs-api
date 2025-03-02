package api

import (
	"context"
	"encoding/json"
	"mime"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/addit-digital/addcache"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmgorilla"
	"go.elastic.co/apm/module/apmlogrus"

	"packs-api/internal/config"
	"packs-api/internal/utils"
)

func notFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSONError(w, http.StatusNotFound, "page not found")
	}
}

type Server struct {
	srv                    *http.Server
	skipHealthCheckLogging bool
	ObjectIDGenerator      utils.ObjectIDGenerator
	Time                   utils.Time
	Log                    *logrus.Entry
	Cache                  addcache.Cache
}

func NewServer(cfg *config.Config, logger *logrus.Entry) *Server {
	s := new(Server)
	router := mux.NewRouter()
	router.NotFoundHandler = notFoundHandler()
	apmgorilla.Instrument(router)

	router.Use(s.writeSecurityHeaders, s.loggingHandlerWrapper)

	h := handlers.CORS(
		handlers.AllowedOrigins(cfg.AllowedOrigins),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "access-control-allow-origin", "Accept-Encoding", "X-CSRF-Token", "Authorization", "X-API-KEY"}),
	)(router)

	randomObjectIDGenerator := utils.NewRandomObjectIDGenerator()
	realTime := utils.NewRealTime()

	s.srv = &http.Server{Addr: cfg.Addr, Handler: h, ReadHeaderTimeout: 10 * time.Second}
	s.skipHealthCheckLogging = cfg.SkipHealthCheckLogging
	s.Log = logger
	s.ObjectIDGenerator = randomObjectIDGenerator
	s.Time = realTime

	pathPrefix := cfg.PathPrefix

	router.HandleFunc(pathPrefix+"/status", s.Recover(s.HandleStatus(), true)).Methods(http.MethodGet)

	router.HandleFunc(pathPrefix+"/orders", s.HandleCreateOrder(cfg.MongoDB)).Methods(http.MethodPost)
	router.HandleFunc(pathPrefix+"/orders", s.HandleGetAllOrders(cfg.MongoDB)).Methods(http.MethodGet)

	return s
}

func (s *Server) HasContentType(r *http.Request, mimetype string) bool {
	contentType := r.Header.Get("Content-type")
	return compareContentTypes(contentType, mimetype)
}

func compareContentTypes(contentType, mimetype string) bool {
	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == mimetype {
			return true
		}
	}

	return false
}

func (s *Server) ListenAndServe() error {
	s.Log.WithField("addr", s.srv.Addr).Info("server is listening...")
	return s.srv.ListenAndServe()
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	return s.srv.ListenAndServeTLS(certFile, keyFile)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) writeSecurityHeaders(wh http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-store, max-age=0")
		w.Header().Add("Pragma", "no-cache")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Header().Add("X-Frame-Options", "SAMEORIGINS")
		w.Header().Add("X-Content-Type-Options", "nosniff")
		w.Header().Add("X-XSS-Protection", "1; mode=block")

		wh.ServeHTTP(w, r)
	})
}

func (s *Server) WriteJSONError(w http.ResponseWriter, code int, msg string) {
	writeJSONError(w, code, msg)
}

func (s *Server) Recover(next http.HandlerFunc, printstack bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				if printstack {
					debug.PrintStack()
				}

				s.WriteJSONError(w, http.StatusInternalServerError, "internal http error")
			}
		}()

		next(w, r)
	}
}

func writeJSONError(w http.ResponseWriter, c int, msg string) {
	errObject := map[string]interface{}{"error": true, "code": c, "message": msg}
	res, _ := json.Marshal(errObject)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	_, _ = w.Write(res)
}

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (lrw *logResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *logResponseWriter) Write(body []byte) (int, error) {
	lrw.body = body
	return lrw.ResponseWriter.Write(body)
}

func newLogResponseWriter(w http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{ResponseWriter: w}
}

func (s *Server) HandleStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Log.Info("Status API called")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	}
}

func (s *Server) loggingHandlerWrapper(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := s.Time.Now()
		lrw := newLogResponseWriter(w)
		wrappedHandler.ServeHTTP(lrw, r)
		duration := s.Time.Now().Sub(start)

		if s.skipHealthCheckLogging && strings.HasSuffix(r.URL.Path, "/packs-api/status/") {
			return
		}

		var resp interface{}
		//nolint:typecheck
		lrwContentType := lrw.Header().Get("Content-Type")

		if compareContentTypes(lrwContentType, "application/json") {
			_ = json.Unmarshal(lrw.body, &resp)
		}

		fields := logrus.Fields{
			"duration":   duration.String(),
			"statusCode": lrw.statusCode,
			"uri":        r.RequestURI,
			"method":     r.Method,
		}

		if resp != nil {
			fields["response"] = resp
		}

		s.Log.
			WithFields(apmlogrus.TraceContext(r.Context())).
			WithFields(fields).
			Info("Request handling took: ", duration)
	})
}
