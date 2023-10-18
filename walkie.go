package walkie

import (
	"io"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

type Walkie interface {
	// Writer returns a new default writer.
	Writer() io.Writer

	// ToFile Logs to file.
	ToFile() gin.HandlerFunc
}

type logger struct {
	*slog.Logger
}

// New returns a new Walkie instance.
func New() Walkie {
	return &logger{slog.Default().WithGroup("gin")}
}

func (l *logger) Writer() io.Writer {
	return writeFunc(func(data []byte) (int, error) {
		l.Debug("Debug Data", "data", string(data))
		return 0, nil
	})
}

func (l *logger) ToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()

		duration := endTime.Sub(startTime)
		logw := l.With(
			"method", c.Request.Method,
			"duration", duration,
			"request-uri", c.Request.RequestURI,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"client-ip", c.ClientIP(),
			"referrer", c.Request.Referer(),
			"request-id", c.Writer.Header().Get("Request-Id"),
		)

		if c.Writer.Status() >= 500 {
			logw.Error("Failed request", "errors", c.Errors)
			return
		}

		logw.Info("Successful request")
	}
}

// writeFunc convert func to io.Writer.
type writeFunc func([]byte) (int, error)

func (fn writeFunc) Write(data []byte) (int, error) {
	return fn(data)
}
