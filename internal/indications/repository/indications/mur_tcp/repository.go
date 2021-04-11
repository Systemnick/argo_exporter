package mur_tcp

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog"

	"github.com/Systemnick/argo_exporter/internal/indications/repository/indications"
)

type Repository struct {
	log *zerolog.Logger
}

func NewRepository(logger *zerolog.Logger) *Repository {
	return &Repository{
		log: logger,
	}
}

func (r Repository) Poll(ip net.IP, metrics *indications.MURMetrics) error {
	timeStart := time.Now()
	msg := "MUR TCP indications poll"
	r.log.Debug().Msg(msg)
	defer func() {
		duration := time.Since(timeStart)
		r.log.Debug().Dur("duration", duration).Msg(msg + " finished")
		metrics.SetRegPollDuration(ip, duration)
	}()

	reg := &MUR{
		ip: ip,
	}

	conn, err := r.connect(ip)
	if err != nil {
		metrics.SetRegStatus(ip, 0)
		return err
	}

	metrics.SetRegStatus(ip, 1)

	reg.conn = conn

	err = r.getIndications(reg, metrics)
	if err != nil {
		return err
	}

	err = r.disconnect(conn)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) connect(ip net.IP) (*ReConn, error) {
	timeStart := time.Now()
	msg := "MUR TCP connect"
	r.log.Debug().Msg(msg)
	defer r.log.Debug().Dur("duration", time.Since(timeStart)).Msg(msg + " finished")

	conn := &ReConn{
		Network:        "tcp",
		Address:        ip.String() + ":" + murPort,
		ConnectRetries: 5,
	}

	err := conn.Reconnect()
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	conn.SetTimeout(15)

	return conn, nil
}

func (r Repository) getIndications(reg *MUR, metrics *indications.MURMetrics) error {
	timeStart := time.Now()
	msg := "MUR TCP get indications"
	r.log.Debug().Msg(msg)
	defer r.log.Debug().Dur("duration", time.Since(timeStart)).Msg(msg + " finished")

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Function getIndications recovered from panic: %+v\n", err)
		}
	}()

	v := reg.GetVersion()
	fmt.Printf("MUR version: %v\n", v)

	reg.Foreplay()

	reg.SetInstantMode()

	regInfo := reg.GetRegistrarInfo()
	if regInfo == nil {
		return errors.New("bad response")
	}

	fmt.Printf("Reqistrar %s adapter count is %d, serial number %d, network address %d\n",
		reg.ip, regInfo.AdapterCount, regInfo.SerialNumber, regInfo.NetworkAddress)

	reg.meterCount = regInfo.AdapterCount

	// reg.Foreplay2()

	// adapter := uint16(46)
	// m := reg.GetAdapterIndications(adapter)
	// fmt.Printf("Adapter %d indications: %v\n", adapter, m)

	reg.Foreplay()

	reg.SetInstantMode()

	reg.SetInstantMode()

	m := reg.GetAllIndications()
	// fmt.Printf("Adapter's indications: %v\n", m)

	metrics.SetMetersIndications(m)

	fmt.Printf("Reqistrar %s poll successfuly finished\n", reg.ip)

	return nil
}

func (r Repository) disconnect(conn *ReConn) error {
	timeStart := time.Now()
	msg := "MUR TCP disconnect"
	r.log.Debug().Msg(msg)
	defer r.log.Debug().Dur("duration", time.Since(timeStart)).Msg(msg + " finished")

	if conn.IsActive() {
		err := conn.Close()
		if err != nil {
			fmt.Printf("closing connection error: %v\n", err)
			return err
		}
	}

	return nil
}
