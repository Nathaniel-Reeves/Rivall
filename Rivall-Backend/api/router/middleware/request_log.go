package middleware

import (
	"net/http"

	"github.com/mitchellh/colorstring"
	"github.com/rs/zerolog/log"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func RequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var s string
		switch r.Method {
		case http.MethodGet:
			s = colorstring.Color("[green]GET")
		case http.MethodPost:
			s = colorstring.Color("[yellow]POST")
		case http.MethodPut:
			s = colorstring.Color("[blue]PUT")
		case http.MethodDelete:
			s = colorstring.Color("[red]DELETE")
		default:
			s = "UNKNOWN"
		}

		msg := "<-- " + s + " REQUEST " + r.URL.Path

		// If the request is an upgrade to websocket request, skip logging
		if r.Header.Get("Upgrade") == "websocket" {
			next.ServeHTTP(w, r)
			return
		}

		// Log the request address, method, path, and headers
		log.Info().
			// Str("headers", r.Header.Get("Content-Type")).
			Msg(msg)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		// Get the status code from the response
		statusCode := lrw.statusCode

		var sc string
		switch statusCode {
		case http.StatusOK:
			sc = colorstring.Color("[green]200 " + http.StatusText(statusCode))
		case http.StatusCreated:
			sc = colorstring.Color("[green]201 " + http.StatusText(statusCode))
		case http.StatusAccepted:
			sc = colorstring.Color("[green]202 " + http.StatusText(statusCode))
		case http.StatusNoContent:
			sc = colorstring.Color("[green]204 " + http.StatusText(statusCode))
		case http.StatusMovedPermanently:
			sc = colorstring.Color("[green]301 " + http.StatusText(statusCode))
		case http.StatusFound:
			sc = colorstring.Color("[green]302 " + http.StatusText(statusCode))
		case http.StatusSeeOther:
			sc = colorstring.Color("[green]303 " + http.StatusText(statusCode))
		case http.StatusNotModified:
			sc = colorstring.Color("[green]304 " + http.StatusText(statusCode))
		case http.StatusUseProxy:
			sc = colorstring.Color("[green]305 " + http.StatusText(statusCode))
		case http.StatusTemporaryRedirect:
			sc = colorstring.Color("[green]307 " + http.StatusText(statusCode))
		case http.StatusBadRequest:
			sc = colorstring.Color("[yellow]400 " + http.StatusText(statusCode))
		case http.StatusUnauthorized:
			sc = colorstring.Color("[yellow]401 " + http.StatusText(statusCode))
		case http.StatusForbidden:
			sc = colorstring.Color("[yellow]403 " + http.StatusText(statusCode))
		case http.StatusNotFound:
			sc = colorstring.Color("[yellow]404 " + http.StatusText(statusCode))
		case http.StatusMethodNotAllowed:
			sc = colorstring.Color("[yellow]405 " + http.StatusText(statusCode))
		case http.StatusNotAcceptable:
			sc = colorstring.Color("[yellow]406 " + http.StatusText(statusCode))
		case http.StatusProxyAuthRequired:
			sc = colorstring.Color("[yellow]407 " + http.StatusText(statusCode))
		case http.StatusRequestTimeout:
			sc = colorstring.Color("[yellow]408 " + http.StatusText(statusCode))
		case http.StatusConflict:
			sc = colorstring.Color("[yellow]409 " + http.StatusText(statusCode))
		case http.StatusGone:
			sc = colorstring.Color("[yellow]410 " + http.StatusText(statusCode))
		case http.StatusLengthRequired:
			sc = colorstring.Color("[yellow]411 " + http.StatusText(statusCode))
		case http.StatusPreconditionFailed:
			sc = colorstring.Color("[yellow]412 " + http.StatusText(statusCode))
		case http.StatusRequestEntityTooLarge:
			sc = colorstring.Color("[yellow]413" + http.StatusText(statusCode))
		case http.StatusRequestURITooLong:
			sc = colorstring.Color("[yellow]414 " + http.StatusText(statusCode))
		case http.StatusUnsupportedMediaType:
			sc = colorstring.Color("[yellow]415 " + http.StatusText(statusCode))
		case http.StatusRequestedRangeNotSatisfiable:
			sc = colorstring.Color("[yellow]416 " + http.StatusText(statusCode))
		case http.StatusExpectationFailed:
			sc = colorstring.Color("[yellow]417 " + http.StatusText(statusCode))
		case http.StatusTeapot:
			sc = colorstring.Color("[yellow]418 " + http.StatusText(statusCode))
		case http.StatusMisdirectedRequest:
			sc = colorstring.Color("[yellow]421 " + http.StatusText(statusCode))
		case http.StatusUnprocessableEntity:
			sc = colorstring.Color("[yellow]422 " + http.StatusText(statusCode))
		case http.StatusLocked:
			sc = colorstring.Color("[yellow]423 " + http.StatusText(statusCode))
		case http.StatusFailedDependency:
			sc = colorstring.Color("[yellow]424 " + http.StatusText(statusCode))
		case http.StatusUpgradeRequired:
			sc = colorstring.Color("[yellow]426 " + http.StatusText(statusCode))
		case http.StatusPreconditionRequired:
			sc = colorstring.Color("[yellow]428 " + http.StatusText(statusCode))
		case http.StatusTooManyRequests:
			sc = colorstring.Color("[yellow]429 " + http.StatusText(statusCode))
		case http.StatusRequestHeaderFieldsTooLarge:
			sc = colorstring.Color("[yellow]431 " + http.StatusText(statusCode))
		case http.StatusUnavailableForLegalReasons:
			sc = colorstring.Color("[yellow]451 " + http.StatusText(statusCode))
		case http.StatusInternalServerError:
			sc = colorstring.Color("[red]500 " + http.StatusText(statusCode))
		case http.StatusNotImplemented:
			sc = colorstring.Color("[red]501 " + http.StatusText(statusCode))
		case http.StatusBadGateway:
			sc = colorstring.Color("[red]502 " + http.StatusText(statusCode))
		case http.StatusServiceUnavailable:
			sc = colorstring.Color("[red]503 " + http.StatusText(statusCode))
		case http.StatusGatewayTimeout:
			sc = colorstring.Color("[red]504 " + http.StatusText(statusCode))
		case http.StatusHTTPVersionNotSupported:
			sc = colorstring.Color("[red]505 " + http.StatusText(statusCode))
		case http.StatusVariantAlsoNegotiates:
			sc = colorstring.Color("[red]506 " + http.StatusText(statusCode))
		case http.StatusInsufficientStorage:
			sc = colorstring.Color("[red]507 " + http.StatusText(statusCode))
		case http.StatusLoopDetected:
			sc = colorstring.Color("[red]508 " + http.StatusText(statusCode))
		case http.StatusNotExtended:
			sc = colorstring.Color("[red]510 " + http.StatusText(statusCode))
		case http.StatusNetworkAuthenticationRequired:
			sc = colorstring.Color("[red]511 " + http.StatusText(statusCode))
		default:
			sc = "UNKNOWN"
		}

		msg = "--> " + sc + " RESPONSE " + r.URL.Path

		// Log the response status code and headers
		log.Info().
			// Str("headers", w.Header().Get("Content-Type")).
			Msg(msg)
	})
}
