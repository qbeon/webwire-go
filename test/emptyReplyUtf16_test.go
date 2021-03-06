package test

import (
	"context"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/require"
)

// TestEmptyReplyUtf16 tests returning empty UTF16 replies
func TestEmptyReplyUtf16(t *testing.T) {
	// Initialize webwire server given only the request
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Request: func(
				_ context.Context,
				_ wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Return empty reply
				return wwr.Payload{Encoding: wwr.EncodingUtf16}, nil
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	sock, _ := setup.NewClientSocket()

	// Send request and await an empty binary reply
	reply := request(t, sock, 64, []byte("r"), payload.Payload{})
	require.Equal(t, message.MsgReplyUtf16, reply.MsgType)
	require.Equal(t, payload.Utf16, reply.MsgPayload.Encoding)
	require.Equal(t, []byte(nil), reply.MsgPayload.Data)
}
