package webwire

import (
	"fmt"
	"time"

	"github.com/qbeon/webwire-go/transport"
)

func (srv *server) writeConfMessage(sock transport.Socket) error {
	writer, err := sock.GetWriter()
	if err != nil {
		return fmt.Errorf(
			"couldn't get writer for configuration message: %s",
			err,
		)
	}

	if _, err := writer.Write(srv.configMsg); err != nil {
		return fmt.Errorf("couldn't write configuration message: %s", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("couldn't close writer: %s", err)
	}

	return nil
}

func (srv *server) handleConnection(
	connectionOptions ConnectionOptions,
	userAgent []byte,
	sock transport.Socket,
) {
	// Send server configuration message
	if err := srv.writeConfMessage(sock); err != nil {
		srv.errorLog.Println("couldn't write config message: ", err)
		if closeErr := sock.Close(); closeErr != nil {
			srv.errorLog.Println("couldn't close socket: ", closeErr)
		}
		return
	}

	if err := sock.SetReadDeadline(
		time.Now().Add(srv.options.ReadTimeout),
	); err != nil {
		srv.errorLog.Printf("couldn't set read deadline: %s", err)
		return
	}

	// Register connected client
	connection := newConnection(
		sock,
		userAgent,
		srv,
		connectionOptions,
	)

	srv.connectionsLock.Lock()
	srv.connections = append(srv.connections, connection)
	srv.connectionsLock.Unlock()

	// Call hook on successful connection
	srv.impl.OnClientConnected(connection)

	for {
		// Get a message buffer
		msg := srv.messagePool.Get()

		if !connection.IsActive() {
			msg.Close()
			connection.Close()
			srv.impl.OnClientDisconnected(connection, nil)
			break
		}

		// Await message
		if err := sock.Read(msg); err != nil {
			msg.Close()

			if err.IsAbnormalCloseErr() {
				srv.warnLog.Printf("abnormal closure error: %s", err)
			}

			connection.Close()
			srv.impl.OnClientDisconnected(connection, err)
			break
		}

		// Parse & handle the message
		if err := srv.handleMessage(connection, msg); err != nil {
			srv.errorLog.Print("message handler failed: ", err)
		}
		msg.Close()
	}
}
