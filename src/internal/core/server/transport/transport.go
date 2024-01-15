package transport

import (
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/kodflow/kitsune/src/internal/core/server/transport/proto/generated"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger"
	"google.golang.org/protobuf/proto"
)

func NewRequest(id uuid.UUID) *generated.Request {
	return &generated.Request{
		Id:      id.String(), // Set the Request ID to the UUID string.
		Headers: map[string]*generated.Header{},
	}
}

func NewReponse() *generated.Response {
	return &generated.Response{
		Status:  500,
		Headers: map[string]*generated.Header{},
	}
}

func New() *Exchange {
	id, _ := uuid.NewRandom() // Generate a new random UUID.
	req := NewRequest(id)

	return &Exchange{
		req: req,
	}
}

func (e *Exchange) Wait() {
	if e.res == nil && e.answer == nil {
		e.answer = make(chan struct{}, 1)
		<-e.answer
	}
}

func (e *Exchange) RequestFromTCP(b []byte) {
	// Unmarshal the input byte array into the request struct
	err := proto.Unmarshal(b, e.req)
	if logger.Error(err) {
		e.res.Status = http.StatusBadRequest
		return
	}

	e.res = NewReponse()
}

func (e *Exchange) ResponseFromTCP() []byte {
	e.res.Id = e.req.Id
	b, err := proto.Marshal(e.res)

	if logger.Error(err) {
		return []byte{}
	}

	return b
}

func (e *Exchange) RequestFromHTTP(r *http.Request) {
	e.req.Method = r.Method
	e.req.Endpoint = r.URL.String()
	for k, v := range r.Header {
		e.req.Headers[k] = &generated.Header{Items: v}
	}

	// Read the request body for specific HTTP methods
	if r.Method == "POST" || r.Method == "PATCH" || r.Method == "PUT" {
		// Read the request body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			e.res.Status = http.StatusBadRequest
			return
		}

		e.req.Body = body
	}

	e.res = NewReponse()
}

func (e *Exchange) ResponseFromHTTP(w http.ResponseWriter) {
	w.Header().Set("request-id", e.req.Id)
	w.WriteHeader(int(e.res.Status))
	w.Write(e.res.Body)
}

func (e *Exchange) Request() *generated.Request {
	return e.req
}

func (e *Exchange) Response(res ...*generated.Response) *generated.Response {
	if len(res) > 0 {
		e.res = res[0]
		if e.answer != nil {
			e.answer <- struct{}{}
		}
	}

	return e.res
}

type Exchange struct {
	req    *generated.Request
	res    *generated.Response
	answer chan struct{}
}
