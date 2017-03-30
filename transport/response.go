package transport

import (
	"github.com/gospackler/metre/logging"
	zmq "github.com/pebbe/zmq4"
	"go.uber.org/zap"
)

type Process interface {
	GetResponse(string) string
}

type RespConn struct {
	Conn *Connection
}

func NewRespConn(uri string) (*RespConn, error) {
	conn, err := NewConnection()
	if err != nil {
		return nil, err
	}

	err = conn.Connect(uri, zmq.REP)
	if err != nil {
		return nil, err
	}

	return &RespConn{
		Conn: conn,
	}, nil
}

// Call this function from a goroutine
func (r *RespConn) Listen(process Process, id int) error {
	// FIXME : Probably stream the errors and log it in if the server
	// continues to crash with listen errors.
	for {
		//  Wait for next request from client
		req, err := r.Conn.Sock.Recv(0)
		if err != nil {
			return err
		}
		resp := process.GetResponse(req)
		logging.Logger.Debug("processed response from run",
			zap.Int("id", id),
			zap.String("response", resp),
		)
		// Send reply back to client
		_, err = r.Conn.Sock.Send(resp, zmq.DONTWAIT)
		if err != nil {
			return err
		}
	}
}

func (r *RespConn) Close() {
	r.Conn.Close()
}
