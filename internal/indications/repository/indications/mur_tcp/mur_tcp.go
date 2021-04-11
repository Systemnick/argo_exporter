package mur_tcp

import (
	"fmt"
	"net"
	"time"
)

const murPort = "5000"

type ReConn struct {
	Network        string
	Address        string
	ConnectRetries int
	conn           net.Conn
}

type MUR struct {
	ip         net.IP
	conn       *ReConn
	meterCount uint16
}

func (r *ReConn) Reconnect() error {
	err := r.Close()
	if err != nil {
		return err
	}

	for i := 0; i < r.ConnectRetries; i++ {
		r.conn, err = net.Dial(r.Network, r.Address)
		if err == nil {
			return nil
		}
	}

	return err
}

func (r *ReConn) IsActive() bool {
	return r.conn != nil
}

func (r *ReConn) Close() error {
	if r.conn == nil {
		return nil
	}

	return r.conn.Close()
}

func (r *ReConn) SetTimeout(d time.Duration) {
	deadline := time.Now().Add(d)
	err := r.conn.SetDeadline(deadline)
	if err != nil {
		fmt.Printf("set timeout error: %v\n", err)
		err = r.Reconnect()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		deadline = time.Now().Add(d)
		err := r.conn.SetDeadline(deadline)
		if err != nil {
			fmt.Printf("set timeout error: %v\n", err)
		}
	}
}

// type Repository struct {
// 	conn []net.Conn
// }
//
// func NewRepository(ips []net.IP) *Repository {
// 	connects := make([]net.Conn, len(ips))
// 	for i, ip := range ips {
// 		conn, err := net.Dial("tcp", ip.String()+":"+murPort)
// 		if err != nil {
// 			panic(err)
// 		}
// 		connects[i] = conn
// 	}
// 	return &Repository{conn: connects}
// }
//
// func (r *Repository) GetIndications(ctx context.Context, meter indications.MeterID) (indications map[string]string, err error) {
// 	panic("implement me")
// }
