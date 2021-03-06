package message_test

import (
	"testing"

	"github.com/qbeon/webwire-go/message"
	pld "github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/require"
)

/****************************************************************\
	Constructors
\****************************************************************/

// TestWriteMsgNamelessReq tests WriteMsgNamelessRequest
func TestWriteMsgNamelessReq(t *testing.T) {
	id := genRndMsgIdentifier()
	// sessionKey := generateSessionKey()
	sessionKey := "somesamplesessionkey"

	// Compose encoded message
	// Add type flag
	expected := []byte{message.MsgRequestRestoreSession}
	// Add identifier
	expected = append(expected, id[:]...)
	// Add session key to payload
	expected = append(expected, []byte(sessionKey)...)

	writer := &testWriter{}
	require.NoError(t, message.WriteMsgNamelessRequest(
		writer,
		message.MsgRequestRestoreSession,
		id,
		[]byte(sessionKey),
	))
	require.Equal(t, expected, writer.buf)
	require.True(t, writer.closed)
}

// TestWriteMsgReqBinary tests WriteMsgRequest
// using default binary payload encoding
func TestWriteMsgReqBinary(t *testing.T) {
	id := genRndMsgIdentifier()
	name := genRndName(1, 255)
	payload := pld.Payload{
		Encoding: pld.Binary,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{message.MsgRequestBinary}
	// Add identifier
	expected = append(expected, id[:]...)
	// Add name length flag
	expected = append(expected, byte(len(name)))
	// Add name
	expected = append(expected, []byte(name)...)
	// Add payload
	// (skip header padding byte, not necessary in case of binary encoding)
	expected = append(expected, payload.Data...)

	writer := &testWriter{}
	require.NoError(t, message.WriteMsgRequest(
		writer,
		id,
		name,
		payload.Encoding,
		payload.Data,
		true,
	))
	require.Equal(t, expected, writer.buf)
	require.True(t, writer.closed)
}

// TestWriteMsgReqUtf8 tests WriteMsgRequest using UTF8 payload encoding
func TestWriteMsgReqUtf8(t *testing.T) {
	id := genRndMsgIdentifier()
	name := genRndName(1, 255)
	payload := pld.Payload{
		Encoding: pld.Utf8,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{message.MsgRequestUtf8}
	// Add identifier
	expected = append(expected, id[:]...)
	// Add name length flag
	expected = append(expected, byte(len(name)))
	// Add name
	expected = append(expected, []byte(name)...)
	// Add payload
	// (skip header padding byte, not necessary in case of UTF8 encoding)
	expected = append(expected, payload.Data...)

	writer := &testWriter{}
	require.NoError(t, message.WriteMsgRequest(
		writer,
		id,
		name,
		payload.Encoding,
		payload.Data,
		true,
	))
	require.Equal(t, expected, writer.buf)
	require.True(t, writer.closed)
}

// TestWriteMsgReqUtf16 tests WriteMsgRequest using UTF8 payload encoding
func TestWriteMsgReqUtf16(t *testing.T) {
	id := genRndMsgIdentifier()
	name := genRndName(1, 255)
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte{'r', 0, 'a', 0, 'n', 0, 'd', 0, 'o', 0, 'm', 0},
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{message.MsgRequestUtf16}
	// Add identifier
	expected = append(expected, id[:]...)
	// Add name length flag
	expected = append(expected, byte(len(name)))
	// Add name
	expected = append(expected, []byte(name)...)
	// Add header padding if necessary
	if len(name)%2 != 0 {
		expected = append(expected, byte(0))
	}
	// Add payload
	expected = append(expected, payload.Data...)

	writer := &testWriter{}
	require.NoError(t, message.WriteMsgRequest(
		writer,
		id,
		name,
		payload.Encoding,
		payload.Data,
		true,
	))
	require.Equal(t, expected, writer.buf)
	require.True(t, writer.closed)
}

// TestWriteMsgReqUtf16OddNameLen tests WriteMsgRequest using
// UTF16 payload encoding and a name of odd length
func TestWriteMsgReqUtf16OddNameLen(t *testing.T) {
	id := genRndMsgIdentifier()
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte{'r', 0, 'a', 0, 'n', 0, 'd', 0, 'o', 0, 'm', 0},
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{message.MsgRequestUtf16}
	// Add identifier
	expected = append(expected, id[:]...)
	// Add name length flag
	expected = append(expected, byte(3))
	// Add name of odd length
	expected = append(expected, []byte("odd")...)
	// Add header padding
	expected = append(expected, byte(0))
	// Add payload
	expected = append(expected, payload.Data...)

	writer := &testWriter{}
	require.NoError(t, message.WriteMsgRequest(
		writer,
		id,
		[]byte("odd"),
		payload.Encoding,
		payload.Data,
		true,
	))
	require.Equal(t, expected, writer.buf)
	require.True(t, writer.closed)
}

// TestWriteMsgReplyBinary tests WriteMsgReply
// using default binary payload encoding
func TestWriteMsgReplyBinary(t *testing.T) {
	id := genRndMsgIdentifier()
	payload := pld.Payload{
		Encoding: pld.Binary,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{message.MsgReplyBinary}
	// Add identifier
	expected = append(expected, id[:]...)

	// Add payload
	expected = append(expected, payload.Data...)

	writer := &testWriter{}
	require.NoError(t, message.WriteMsgReply(
		writer,
		id,
		payload.Encoding,
		payload.Data,
	))
	require.Equal(t, expected, writer.buf)
	require.True(t, writer.closed)
}

// TestWriteMsgReplyUtf8 tests WriteMsgReply using UTF8 payload encoding
func TestWriteMsgReplyUtf8(t *testing.T) {
	id := genRndMsgIdentifier()
	payload := pld.Payload{
		Encoding: pld.Utf8,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{message.MsgReplyUtf8}
	// Add identifier
	expected = append(expected, id[:]...)

	// Add payload
	expected = append(expected, payload.Data...)

	writer := &testWriter{}
	require.NoError(t, message.WriteMsgReply(
		writer,
		id,
		payload.Encoding,
		payload.Data,
	))
	require.Equal(t, expected, writer.buf)
	require.True(t, writer.closed)
}

// TestWriteMsgReplyUtf16 tests WriteMsgReply using UTF16 payload encoding
func TestWriteMsgReplyUtf16(t *testing.T) {
	id := genRndMsgIdentifier()
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte{'r', 0, 'a', 0, 'n', 0, 'd', 0, 'o', 0, 'm', 0},
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{message.MsgReplyUtf16}
	// Add identifier
	expected = append(expected, id[:]...)
	// Add header padding byte (necessary in case of a UTF16 encoded reply)
	expected = append(expected, 0)

	// Add payload
	expected = append(expected, payload.Data...)

	writer := &testWriter{}
	require.NoError(t, message.WriteMsgReply(
		writer,
		id,
		payload.Encoding,
		payload.Data,
	))
	require.Equal(t, expected, writer.buf)
	require.True(t, writer.closed)
}

// TestWriteMsgSigBinary tests WriteMsgSignal
// using the default binary encoding
func TestWriteMsgSigBinary(t *testing.T) {
	name := genRndName(1, 255)
	payload := pld.Payload{
		Encoding: pld.Binary,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{message.MsgSignalBinary}
	// Add name length flag
	expected = append(expected, byte(len(name)))
	// Add name
	expected = append(expected, []byte(name)...)
	// Add payload (skip header padding byte in case of binary encoding)
	expected = append(expected, payload.Data...)

	writer := &testWriter{}
	require.NoError(t, message.WriteMsgSignal(
		writer,
		name,
		payload.Encoding,
		payload.Data,
		true,
	))
	require.Equal(t, expected, writer.buf)
	require.True(t, writer.closed)
}

// TestWriteMsgSigUtf8 tests WriteMsgSignal using UTF8 encoding
func TestWriteMsgSigUtf8(t *testing.T) {
	name := genRndName(1, 255)
	payload := pld.Payload{
		Encoding: pld.Utf8,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{message.MsgSignalUtf8}
	// Add name length flag
	expected = append(expected, byte(len(name)))
	// Add name
	expected = append(expected, []byte(name)...)
	// Add payload (skip header padding byte in case of UTF8 encoding)
	expected = append(expected, payload.Data...)

	writer := &testWriter{}
	require.NoError(t, message.WriteMsgSignal(
		writer,
		name,
		payload.Encoding,
		payload.Data,
		true,
	))
	require.Equal(t, expected, writer.buf)
	require.True(t, writer.closed)
}

// TestWriteMsgSigUtf16 tests WriteMsgSignal using UTF16 encoding
func TestWriteMsgSigUtf16(t *testing.T) {
	name := genRndName(1, 255)
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte{'r', 0, 'a', 0, 'n', 0, 'd', 0, 'o', 0, 'm', 0},
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{message.MsgSignalUtf16}
	// Add name length flag
	expected = append(expected, byte(len(name)))
	// Add name
	expected = append(expected, []byte(name)...)
	// Add header padding if necessary
	if len(name)%2 != 0 {
		expected = append(expected, byte(0))
	}
	// Add payload
	expected = append(expected, payload.Data...)

	writer := &testWriter{}
	require.NoError(t, message.WriteMsgSignal(
		writer,
		name,
		payload.Encoding,
		payload.Data,
		true,
	))
	require.Equal(t, expected, writer.buf)
	require.True(t, writer.closed)
}

// TestWriteMsgSigUtf16OddNameLen tests WriteMsgSignal using UTF16 encoding and
// a name of odd length to ensure a header padding byte is used
func TestWriteMsgSigUtf16OddNameLen(t *testing.T) {
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte{'r', 0, 'a', 0, 'n', 0, 'd', 0, 'o', 0, 'm', 0},
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{message.MsgSignalUtf16}
	// Add name length flag
	expected = append(expected, byte(3))
	// Add name of odd length
	expected = append(expected, []byte("odd")...)
	// Add header padding
	expected = append(expected, byte(0))
	// Add payload
	expected = append(expected, payload.Data...)

	writer := &testWriter{}
	require.NoError(t, message.WriteMsgSignal(
		writer,
		[]byte("odd"),
		payload.Encoding,
		payload.Data,
		true,
	))
	require.Equal(t, expected, writer.buf)
	require.True(t, writer.closed)
}

// TestWriteMsgNotifySessionCreated tests WriteMsgNotifySessionCreated
func TestWriteMsgNotifySessionCreated(t *testing.T) {
	// Compose encoded message
	// Write type flag
	expected := []byte{message.MsgNotifySessionCreated}
	// Write session info payload
	expected = append(expected, []byte("session info")...)

	writer := &testWriter{}
	require.NoError(t, message.WriteMsgNotifySessionCreated(
		writer,
		[]byte("session info"),
	))
	require.Equal(t, expected, writer.buf)
	require.True(t, writer.closed)
}

// TestWriteMsgNotifySessionClosed tests WriteMsgNotifySessionClosed
func TestWriteMsgNotifySessionClosed(t *testing.T) {
	// Compose expected message
	expected := []byte{message.MsgNotifySessionClosed}

	writer := &testWriter{}
	require.NoError(t, message.WriteMsgNotifySessionClosed(writer))
	require.Equal(t, expected, writer.buf)
	require.True(t, writer.closed)
}

// TestWriteMsgHeartbeat tests WriteMsgHeartbeat
func TestWriteMsgHeartbeat(t *testing.T) {
	// Compose expected message
	expected := []byte{message.MsgHeartbeat}

	writer := &testWriter{}
	require.NoError(t, message.WriteMsgHeartbeat(writer))
	require.Equal(t, expected, writer.buf)
	require.True(t, writer.closed)
}
