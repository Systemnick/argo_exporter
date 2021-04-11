package registrars

import (
	"context"
	"net"
)

type Device struct {
	IP       net.IP
}

type Repository interface {
	List(ctx context.Context) (device []Device, err error)
}
