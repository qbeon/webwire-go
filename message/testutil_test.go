package message_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/qbeon/webwire-go/message"
	pld "github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// testWriter is an io.WriteCloser implementation for testing purposes
type testWriter struct {
	closed bool
	buf    []byte
}

// Write implements the io.WriteCloser interface
func (tw *testWriter) Write(p []byte) (int, error) {
	if tw.buf == nil {
		tw.buf = make([]byte, len(p))
		copy(tw.buf, p)
	} else {
		tw.buf = append(tw.buf, p...)
	}
	return len(p), nil
}

// Close implements the io.WriteCloser interface
func (tw *testWriter) Close() error {
	tw.closed = true
	return nil
}

func tryParse(t *testing.T, encoded []byte) (*message.Message, error) {
	msg := message.NewMessage(uint32(len(encoded)))
	typeDetermined, err := msg.ReadBytes(encoded)
	require.True(t, typeDetermined, "Couldn't determine message type")
	return msg, err
}

func tryParseNoErr(t *testing.T, encoded []byte) *message.Message {
	msg, err := tryParse(t, encoded)
	require.NoError(t, err)
	return msg
}

// genRndMsgIdentifier returns a randomly generated message id
func genRndMsgIdentifier() (randomIdentifier []byte) {
	randomIdentifier = make([]byte, 8)
	rand.Read(randomIdentifier)
	return randomIdentifier
}

// genRndByteString returns a randomly generated byte string
func genRndByteString(min, max, power uint) []byte {
	if min%power != 0 {
		panic(fmt.Errorf(
			"Expected minimum byte string length to be power %d but was %d",
			power,
			min,
		))
	}
	if max%power != 0 {
		panic(fmt.Errorf(
			"Expected maximum byte string length to be power %d but was %d",
			power,
			max,
		))
	}

	if min > max {
		panic(fmt.Errorf(
			"Invalid genRndByteString parameters: %d | %d",
			min,
			max,
		))
	}
	const letters = " !\"#$%&'()*+,-./0123456789:;<=>?@" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"[\\]^_`" +
		"abcdefghijklmnopqrstuvwxyz" +
		"{|}~"

	// Determine length
	randomLength := min
	if max != min {
		randomLength = min + uint(rand.Intn(int(max-min)))
	}

	// Check power
	delta := randomLength % power
	if delta != 0 {
		randomLength += delta
	}

	if randomLength < 1 {
		return []byte{}
	}

	// Generate
	str := make([]byte, randomLength)
	for i := range str {
		str[i] = letters[rand.Intn(len(letters))]
	}
	return str
}

// genRndName returns a randomly generated byte-string
func genRndName(min, max uint) []byte {
	if max > 255 || min > max {
		panic(fmt.Errorf("Invalid genRndName parameters: %d | %d", min, max))
	}
	return genRndByteString(min, max, 1)
}

// rndRequestMsg returns a randomly generated request message,
// its randomly generated id, name and payload
func rndRequestMsg(
	messageType byte,
	minNameLen uint,
	maxNameLen uint,
	minPayloadLen uint,
	maxPayloadLen uint,
) (
	encodedMessage []byte,
	id []byte,
	name []byte,
	payload pld.Payload,
) {
	id = genRndMsgIdentifier()
	name = genRndName(1, 255)

	messageLength := 10 + len(name) + len(payload.Data)
	payloadEncoding := pld.Binary
	switch messageType {
	case message.MsgRequestBinary:
	case message.MsgRequestUtf8:
		payloadEncoding = pld.Utf8
	case message.MsgRequestUtf16:
		panic(fmt.Errorf(
			"Consider using rndRequestMsgUtf16" +
				"for UTF16 encoded request messages",
		))
	default:
		panic(fmt.Errorf(
			"Invalid type: %d for request message",
			messageType,
		))
	}

	payload = pld.Payload{
		Encoding: payloadEncoding,
		Data:     genRndByteString(minPayloadLen, maxPayloadLen, 1),
	}

	// Compose encoded message
	encodedMessage = make([]byte, 0, messageLength)

	// Add type flag
	encodedMessage = append(encodedMessage, messageType)

	// Add identifier
	encodedMessage = append(encodedMessage, id...)

	// Add name length flag
	encodedMessage = append(encodedMessage, byte(len(name)))
	// Add name
	encodedMessage = append(encodedMessage, []byte(name)...)

	// Add payload
	encodedMessage = append(encodedMessage, payload.Data...)

	return encodedMessage, id, name, payload
}

// rndRequestMsg returns a randomly generated request message,
// its randomly generated id, name and payload
func rndRequestMsgUtf16(
	minNameLen uint,
	maxNameLen uint,
	minPayloadLen uint,
	maxPayloadLen uint,
) (
	encodedMessage []byte,
	id []byte,
	name []byte,
	payload pld.Payload,
) {
	id = genRndMsgIdentifier()
	name = genRndName(1, 255)
	payload = pld.Payload{
		Encoding: pld.Utf16,
		Data:     genRndByteString(minPayloadLen, maxPayloadLen, 2),
	}

	messageLength := 10 + len(name) + len(payload.Data)
	useHeaderPaddingByte := false

	// Add header padding byte only when the name length is odd
	if len(name)%2 != 0 {
		useHeaderPaddingByte = true
		messageLength++
	}

	// Compose encoded message
	encodedMessage = make([]byte, 0, messageLength)

	// Add type flag
	encodedMessage = append(encodedMessage, message.MsgRequestUtf16)

	// Add identifier
	encodedMessage = append(encodedMessage, id...)

	// Add name length flag
	encodedMessage = append(encodedMessage, byte(len(name)))

	// Add name
	encodedMessage = append(encodedMessage, []byte(name)...)

	// Add header padding byte only when the name length is odd
	if useHeaderPaddingByte {
		encodedMessage = append(encodedMessage, byte(0))
	}

	// Add payload
	encodedMessage = append(encodedMessage, payload.Data...)

	return encodedMessage, id, name, payload
}

// rndReplyMsg returns a randomly generated reply message,
// its randomly generated id and payload
func rndReplyMsg(
	messageType byte,
	minPayloadLen uint,
	maxPayloadLen uint,
) (
	encodedMessage []byte,
	id []byte,
	payload pld.Payload,
) {
	id = genRndMsgIdentifier()

	messageLength := 9 + len(payload.Data)
	payloadEncoding := pld.Binary
	switch messageType {
	case message.MsgReplyBinary:
	case message.MsgReplyUtf8:
		payloadEncoding = pld.Utf8
	case message.MsgReplyUtf16:
		panic(fmt.Errorf(
			"Consider using rndReplyMsgUtf16 for UTF16 encoded messages",
		))
	default:
		panic(fmt.Errorf(
			"Invalid type: %d for request message",
			messageType,
		))
	}

	payload = pld.Payload{
		Encoding: payloadEncoding,
		Data:     genRndByteString(minPayloadLen, maxPayloadLen, 1),
	}

	// Compose encoded message
	encodedMessage = make([]byte, 0, messageLength)

	// Add type flag
	encodedMessage = append(encodedMessage, messageType)

	// Add identifier
	encodedMessage = append(encodedMessage, id...)

	// Add payload
	encodedMessage = append(encodedMessage, payload.Data...)

	return encodedMessage, id, payload
}

// rndReplyMsgUtf16 returns a randomly generated UTF16 encoded reply message,
// its randomly generated id and payload
func rndReplyMsgUtf16(
	minPayloadLen uint,
	maxPayloadLen uint,
) (
	encodedMessage []byte,
	id []byte,
	payload pld.Payload,
) {
	id = genRndMsgIdentifier()
	payload = pld.Payload{
		Encoding: pld.Utf16,
		Data:     genRndByteString(minPayloadLen, maxPayloadLen, 2),
	}

	messageLength := 10 + len(payload.Data)

	// Compose encoded message
	encodedMessage = make([]byte, 0, messageLength)

	// Add type flag
	encodedMessage = append(encodedMessage, message.MsgReplyUtf16)

	// Add identifier
	encodedMessage = append(encodedMessage, id...)

	// Add header padding byte
	encodedMessage = append(encodedMessage, byte(0))

	// Add payload
	encodedMessage = append(encodedMessage, payload.Data...)

	return encodedMessage, id, payload
}

// rndSignalMsg returns a randomly generated binary encoded signal message,
// and its randomly generated payload
func rndSignalMsg(
	messageType byte,
	minNameLen uint,
	maxNameLen uint,
	minPayloadLen uint,
	maxPayloadLen uint,
) (
	encodedMessage []byte,
	name []byte,
	payload pld.Payload,
) {
	name = genRndName(1, 255)

	messageLength := 2 + len(name) + len(payload.Data)
	payloadEncoding := pld.Binary
	switch messageType {
	case message.MsgSignalBinary:
	case message.MsgSignalUtf8:
		payloadEncoding = pld.Utf8
	case message.MsgSignalUtf16:
		panic(fmt.Errorf(
			"Consider using rndSignalMsgUtf16" +
				"for UTF16 encoded signal messages",
		))
	default:
		panic(fmt.Errorf(
			"Invalid type: %d for signal message",
			messageType,
		))
	}

	payload = pld.Payload{
		Encoding: payloadEncoding,
		Data:     genRndByteString(minPayloadLen, maxPayloadLen, 1),
	}

	// Compose encoded message
	encodedMessage = make([]byte, 0, messageLength)

	// Add type flag
	encodedMessage = append(encodedMessage, messageType)

	// Add name length flag
	encodedMessage = append(encodedMessage, byte(len(name)))

	// Add name
	encodedMessage = append(encodedMessage, []byte(name)...)

	// Add payload
	encodedMessage = append(encodedMessage, payload.Data...)

	return encodedMessage, name, payload
}

// rndSignalMsgUtf16 returns a randomly generated signal message,
// its randomly generated name and payload
func rndSignalMsgUtf16(
	minNameLen uint,
	maxNameLen uint,
	minPayloadLen uint,
	maxPayloadLen uint,
) (
	encodedMessage []byte,
	name []byte,
	payload pld.Payload,
) {
	name = genRndName(1, 255)
	payload = pld.Payload{
		Encoding: pld.Utf16,
		Data:     genRndByteString(minPayloadLen, maxPayloadLen, 2),
	}

	messageLength := 2 + len(name) + len(payload.Data)
	useHeaderPaddingByte := false

	// Add header padding byte only when the name length is odd
	if len(name)%2 != 0 {
		useHeaderPaddingByte = true
		messageLength++
	}

	// Compose encoded message
	encodedMessage = make([]byte, 0, messageLength)

	// Add type flag
	encodedMessage = append(encodedMessage, message.MsgSignalUtf16)

	// Add name length flag
	encodedMessage = append(encodedMessage, byte(len(name)))

	// Add name
	encodedMessage = append(encodedMessage, []byte(name)...)

	// Add header padding byte only when the name length is odd
	if useHeaderPaddingByte {
		encodedMessage = append(encodedMessage, byte(0))
	}

	// Add payload
	encodedMessage = append(encodedMessage, payload.Data...)

	return encodedMessage, name, payload
}
