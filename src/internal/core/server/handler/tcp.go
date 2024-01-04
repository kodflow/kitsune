package handler

import (
	"time"

	"github.com/kodflow/kitsune/src/internal/core/server/router"
	"github.com/kodflow/kitsune/src/internal/core/server/transport"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/metrics"
	"google.golang.org/protobuf/proto"
)

var tcpm = metrics.GetAverage("tcp/req", time.Second)

// TCPHandler handles TCP requests by unmarshalling, processing, and marshalling responses.
// It is responsible for converting raw byte data into a structured request, processing it
// using a router, and then returning the structured response as byte data. It handles
// errors at each step by returning an empty response in case of failure.
//
// Parameters:
// - b: []byte Raw byte array representing a TCP request.
//
// Returns:
// - []byte: Processed response as a byte array. Returns an empty response in case of errors.
func TCPHandler(b []byte) []byte {
	tcpm.Hit()
	//metrics.GetCounter("tcp/req")
	// Initialize a new transport request and response
	req, res := transport.New()

	// Unmarshal the input byte array into the request struct
	err := proto.Unmarshal(b, req)

	// Set the process ID in the response to match the request
	res.Pid = req.Pid

	// Return an empty response if there's an error in unmarshalling
	if err != nil {
		return transport.Empty
	}

	// Resolve the request using the router and update the response
	err = router.Resolve(req, res)

	// Return an empty response if there's an error in processing the request
	if err != nil {
		return transport.Empty
	}

	// Marshal the response back into a byte array
	b, err = proto.Marshal(res)

	// Return an empty response if there's an error in marshalling
	if err != nil {
		return transport.Empty
	}

	// Return the processed response
	return b
}
