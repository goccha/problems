package problems

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/goccha/http-constants/pkg/headers"
	"github.com/goccha/http-constants/pkg/mimetypes"
	"github.com/goccha/logging/log"
)

type Renderer interface {
	JSON(ctx context.Context, w http.ResponseWriter)
	XML(ctx context.Context, w http.ResponseWriter)
}

func setHeader(ctx context.Context, w http.ResponseWriter, status int, mimetype string) {
	w.Header().Set(headers.ContentType, mimetype)
	if status > 0 {
		w.WriteHeader(status)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func WriteJson(ctx context.Context, w http.ResponseWriter, status int, v interface{}) {
	setHeader(ctx, w, status, mimetypes.ProblemJson)
	if ex, ok := v.(Extendable); ok {
		if ex.Extended() {
			v, _ = ex.Map()
		}
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.EmbedObject(ctx, log.Warn(ctx).Err(err)).Send()
	}
}

func WriteXml(ctx context.Context, w http.ResponseWriter, status int, v interface{}) {
	setHeader(ctx, w, status, mimetypes.ProblemXml)
	if ex, ok := v.(Extendable); ok {
		if ex.Extended() {
			v, _ = ex.Map()
		}
	}
	if err := xml.NewEncoder(w).Encode(v); err != nil {
		log.EmbedObject(ctx, log.Warn(ctx).Err(err)).Send()
	}
}
