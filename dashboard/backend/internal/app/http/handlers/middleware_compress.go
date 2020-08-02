package handlers

import (
	"compress/flate"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"strings"
)

const (
	gzipEncoding  = "gzip"
	flateEncoding = "deflate"

	acceptEncodingHeader  = "Accept-Encoding"
	contentTypeHeader     = "Content-Type"
	contentLengthHeader   = "Content-Length"
	contentEncodingHeader = "Content-Encoding"
)

type CompressMiddleware struct {
	level  int
	logger Logger
}

func NewCompressMiddleware(level int, logger Logger) *CompressMiddleware {
	if level < gzip.HuffmanOnly || level > gzip.BestCompression {
		level = gzip.DefaultCompression
	}

	return &CompressMiddleware{level: level, logger: logger}
}

func (m *CompressMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// detect what encoding to use
		encoding := m.getReqEncoding(r)

		// always add Accept-Encoding to Vary to prevent intermediate caches corruption
		w.Header().Add("Vary", acceptEncodingHeader)

		// wrap the ResponseWriter with the writer for the chosen encoding
		encWriter := m.getEncWriter(w, encoding)
		if encWriter == nil {
			next.ServeHTTP(w, r)
			return
		}

		defer encWriter.Close()

		w.Header().Set(contentEncodingHeader, encoding)
		r.Header.Del(acceptEncodingHeader)

		hijacker, ok := w.(http.Hijacker)
		if !ok { /* w is not Hijacker... oh well... */
			hijacker = nil
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			flusher = nil
		}

		w = &compressWriter{
			Writer:         encWriter,
			ResponseWriter: w,
			Hijacker:       hijacker,
			Flusher:        flusher,
		}

		next.ServeHTTP(w, r)
	})
}

func (m *CompressMiddleware) getReqEncoding(r *http.Request) (encoding string) {
	for _, curEnc := range strings.Split(r.Header.Get(acceptEncodingHeader), ",") {
		curEnc = strings.TrimSpace(curEnc)

		if curEnc == gzipEncoding || curEnc == flateEncoding {
			return curEnc
		}
	}

	return ""
}

func (m *CompressMiddleware) getEncWriter(w io.Writer, encoding string) io.WriteCloser {
	// wrap the ResponseWriter with the writer for the chosen encoding
	var (
		encWriter io.WriteCloser
		err       error
	)

	switch encoding {
	case gzipEncoding:
		encWriter, err = gzip.NewWriterLevel(w, m.level)
	case flateEncoding:
		encWriter, err = flate.NewWriter(w, m.level)
	default:
		return nil
	}

	if err != nil {
		m.logger.Errorf(context.TODO(), "init encoder failed: %v", err)
		return nil
	}

	return encWriter
}

type compressWriter struct {
	io.Writer
	http.ResponseWriter
	http.Hijacker
	http.Flusher
}

func (w *compressWriter) WriteHeader(c int) {
	w.ResponseWriter.Header().Del(contentLengthHeader)
	w.ResponseWriter.WriteHeader(c)
}

func (w *compressWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *compressWriter) Write(b []byte) (int, error) {
	h := w.ResponseWriter.Header()

	if h.Get(contentTypeHeader) == "" {
		h.Set(contentTypeHeader, http.DetectContentType(b))
	}

	h.Del(contentLengthHeader)

	return w.Writer.Write(b)
}

func (w *compressWriter) Flush() {
	type flusher interface {
		Flush() error
	}

	// Flush compressed data if compressor supports it.
	if f, ok := w.Writer.(flusher); ok {
		_ = f.Flush()
	}

	// Flush HTTP response.
	if w.Flusher != nil {
		w.Flusher.Flush()
	}
}
