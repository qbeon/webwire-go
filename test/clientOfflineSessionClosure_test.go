package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClientOfflineSessionClosure tests offline session closure
func TestClientOfflineSessionClosure(t *testing.T) {
	sessionStorage := make(map[string]*wwr.Session)

	currentStep := 1
	var createdSession *wwr.Session

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Request: func(
				_ context.Context,
				conn wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				if currentStep == 2 {
					// Expect the session to be removed
					assert.False(t,
						conn.HasSession(),
						"Expected client to be anonymous",
					)
					return wwr.Payload{}, nil
				}

				// Try to create a new session
				err := conn.CreateSession(nil)
				assert.NoError(t, err)
				return wwr.Payload{}, err
			},
		},
		wwr.ServerOptions{
			SessionManager: &SessionManager{
				// Saves the session
				SessionCreated: func(conn wwr.Connection) error {
					session := conn.Session()
					sessionStorage[session.Key] = session
					return nil
				},
				// Finds session by key
				SessionLookup: func(key string) (
					wwr.SessionLookupResult,
					error,
				) {
					// Expect the key of the created session to be looked up
					assert.Equal(t, createdSession.Key, key)

					assert.Contains(t, sessionStorage, key)
					session := sessionStorage[key]
					// Session found
					return wwr.NewSessionLookupResult(
						session.Creation,                      // Creation
						session.LastLookup,                    // LastLookup
						wwr.SessionInfoToVarMap(session.Info), // Info
					), nil
				},
			},
		},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		TestClientHooks{},
	)

	/*****************************************************************\
		Step 1 - Create session and disconnect
	\*****************************************************************/

	// Create a new session
	reply, err := client.Connection.Request(
		context.Background(),
		[]byte("login"),
		wwr.Payload{
			Encoding: wwr.EncodingBinary,
			Data:     []byte("auth"),
		},
	)
	require.NoError(t, err)
	reply.Close()

	createdSession = client.Connection.Session()

	// Disconnect client without closing the session
	client.Connection.Close()

	// Ensure the session isn't lost
	require.NotEqual(t,
		wwrclt.StatusConnected, client.Connection.Status(),
		"Client is expected to be disconnected",
	)
	require.NotEqual(t,
		"", client.Connection.Session().Key,
		"Session lost after disconnection",
	)

	/*****************************************************************\
		Step 2 - Close session, reconnect and verify
	\*****************************************************************/
	currentStep = 2

	require.NoError(t,
		client.Connection.CloseSession(),
		"Offline session closure failed",
	)

	// Ensure the session is removed locally
	require.Nil(t, client.Connection.Session(), "Session not removed")

	// Reconnect
	require.NoError(t, client.Connection.Connect())

	// Ensure the client is anonymous
	reply, err = client.Connection.Request(
		context.Background(),
		[]byte("verify-restored"),
		wwr.Payload{
			Encoding: wwr.EncodingBinary,
			Data:     []byte("is_restored?"),
		},
	)
	require.NoError(t, err)
	reply.Close()
}
