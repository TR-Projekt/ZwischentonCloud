package servertools

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Middleware(traceLogger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			requestStart := time.Now()
			requestID := r.Header.Get("X-Request-ID")
			tlog := traceLogger.With().Str("request_id", requestID).Logger()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {

				requestEnd := time.Now()
				status := ww.Status()

				// Recover and record stack traces in case of a panic to global logger
				if rec := recover(); rec != nil {
					log.Error().
						Interface("recover_info", rec).
						Str("request_id", requestID).
						Bytes("debug_stack", debug.Stack()).
						Msg("Recovered from panicking routine")
					RespondError(ww, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				}

				if status < 300 {
					// log successfull requests at trace lvl to trace logger
					tlog.Trace().
						Fields(map[string]interface{}{
							"url":        r.Host + r.URL.Path,
							"method":     r.Method,
							"status":     status,
							"latency_ms": float64(requestEnd.Sub(requestStart).Nanoseconds()) / 1000000.0,
							"bytes_out":  ww.BytesWritten(),
						}).
						Msg("Incoming request")
				} else {
					// log failed requests at debug lvl to global logger
					log.Debug().
						Fields(map[string]interface{}{
							"request_id": requestID,
							"url":        r.Host + r.URL.Path,
							"method":     r.Method,
							"status":     status,
							"latency_ms": float64(requestEnd.Sub(requestStart).Nanoseconds()) / 1000000.0,
							"bytes_out":  ww.BytesWritten(),
						}).
						Msg("Incoming request")
				}
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

func InitializeGlobalLogger(logfile string, console bool) {

	logFile, err := NewRollingFile(logfile)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to allocate new rolling log file for global logger")
	}

	var writers []io.Writer
	if console {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	writers = append(writers, logFile)

	multiWriter := io.MultiWriter(writers...)
	logger := zerolog.New(multiWriter).With().Timestamp().Logger()
	log.Logger = logger
}

func TraceLogger(logfile string) *zerolog.Logger {

	logFile, err := NewRollingFile(logfile)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to allocate new rolling log file for trace logger")
	}

	var writers []io.Writer
	writers = append(writers, logFile)

	multiWriter := io.MultiWriter(writers...)
	var traceLogger = zerolog.New(multiWriter).With().Timestamp().Str("type", "access").Logger()

	return &traceLogger
}

func NewRollingFile(file string) (io.Writer, error) {

	dir, _ := filepath.Split(file)
	if err := os.MkdirAll(dir, 0744); err != nil {
		return nil, errors.New("Can't create log directory at:'" + dir + "' with error: " + err.Error())
	}

	return &lumberjack.Logger{
		Filename:   file,
		MaxBackups: 10, // files
		MaxSize:    50, // megabytes
		MaxAge:     31, // days
	}, nil
}
