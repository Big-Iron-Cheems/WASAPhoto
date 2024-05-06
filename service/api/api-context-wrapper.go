package api

import (
	"bytes"
	"github.com/Big-Iron-Cheems/WASAPhoto/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type httpRouterHandler func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
func (rt *_router) wrap(fn httpRouterHandler) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		// Log the request details
		logRequestDetails(r, ctx)

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	}
}

func logRequestDetails(r *http.Request, ctx reqcontext.RequestContext) {
	// Log the reqID
	ctx.Logger.Info("Processing request")

	// Log the method and path
	ctx.Logger.Infof("Method: %s, Path: %s", r.Method, r.URL.Path)

	// Log the protocol version, host, remote address, request URI, and content length
	ctx.Logger.Infof("Protocol: %s, Host: %s, RemoteAddr: %s, RequestURI: %s, ContentLength: %d", r.Proto, r.Host, r.RemoteAddr, r.RequestURI, r.ContentLength)

	// If the request was received over a secure channel, log the TLS version and the server name
	if r.TLS != nil {
		ctx.Logger.Infof("TLS Version: %x, ServerName: %s", r.TLS.Version, r.TLS.ServerName)
	}

	// Log the headers
	for name, values := range r.Header {
		for _, value := range values {
			ctx.Logger.Infof("Header: %s: %s", name, value)
		}
	}

	// Log the body
	if r.Body != nil {
		contentType := r.Header.Get("Content-Type")
		if contentType != "" && strings.HasPrefix(contentType, "application/json") {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				ctx.Logger.WithError(err).Error("Error reading request body")
			} else {
				ctx.Logger.Infof("Body: %s", string(bodyBytes))
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		} else if contentType != "" {
			ctx.Logger.Info("Body contains unexpected MIME type, not logging")
		}
	}
}
