package mur_tcp

// import (
// 	"fmt"
// 	"net"
// 	"sync"
//
// 	"github.com/Systemnick/argo_exporter/internal/indications/repository/indications"
// )
//
// type FabricState uint8
//
// const (
// 	StateNew FabricState = iota
// 	StateRunning
// 	StateFinished
// )
//
// type RepositoryFabric struct {
// 	regs    []*MUR
// 	state   FabricState
// 	onReady []func()
// 	mx      sync.Mutex
// }
//
// func NewRepositoryFabric() *RepositoryFabric {
// 	return &RepositoryFabric{}
// }
//
// func (r *RepositoryFabric) New(ip net.IP, meterCount uint16) {
// 	r.regs = append(r.regs, &MUR{
// 		ip:         ip,
// 		meterCount: meterCount,
// 	})
// }
//
// func (r *RepositoryFabric) OpenConnections(onReady func()) {
// 	r.mx.Lock()
// 	defer r.mx.Unlock()
//
// 	r.onReady = append(r.onReady, onReady)
//
// 	if r.state == StateFinished {
// 		go onReady()
// 		return
// 	}
//
// 	if r.state == StateRunning {
// 		return
// 	}
//
// 	r.state = StateRunning
//
// 	wg := sync.WaitGroup{}
// 	for _, reg := range r.regs {
// 		wg.Add(1)
// 		go func(reg *MUR) {
// 			defer wg.Done()
//
// 			conn := &ReConn{
// 				Network:        "tcp",
// 				Address:        reg.ip.String() + ":" + murPort,
// 				ConnectRetries: 5,
// 			}
//
// 			err := conn.Reconnect()
// 			if err != nil {
// 				fmt.Printf("%v\n", err)
// 				return
// 			}
//
// 			reg.conn = conn
// 		}(reg)
// 	}
//
// 	go func() {
// 		wg.Wait()
// 		for _, f := range r.onReady {
// 			go f()
// 		}
// 	}()
// }
//
// func (r *RepositoryFabric) CloseConnections(onReady func()) {
// 	r.mx.Lock()
// 	defer r.mx.Unlock()
//
// 	r.onReady = append(r.onReady, onReady)
//
// 	if r.state == StateFinished {
// 		go onReady()
// 		return
// 	}
//
// 	if r.state == StateRunning {
// 		return
// 	}
//
// 	r.state = StateRunning
//
// 	wg := sync.WaitGroup{}
// 	for _, reg := range r.regs {
// 		wg.Add(1)
// 		go func(reg *MUR) {
// 			defer wg.Done()
//
// 			if reg.conn.IsActive() {
// 				err := reg.conn.Close()
// 				if err != nil {
// 					fmt.Printf("closing connection error: %v\n", err)
// 					return
// 				}
// 			}
// 		}(reg)
// 	}
//
// 	go func() {
// 		wg.Wait()
// 		for _, f := range r.onReady {
// 			go f()
// 		}
// 	}()
// }
//
// func (r *RepositoryFabric) Poll(metrics *indications.MURMetrics, onReady func()) {
// 	wg := &sync.WaitGroup{}
//
// 	for _, reg := range r.regs {
// 		wg.Add(1)
// 		go func(reg *MUR) {
// 			defer wg.Done()
//
// 			if !reg.conn.IsActive() {
// 				metrics.SetRegStatus(reg.ip, 0)
// 				return
// 			}
//
// 			metrics.SetRegStatus(reg.ip, 1)
//
// 			v := reg.GetVersion()
// 			fmt.Printf("MUR version: %v\n", v)
//
// 			reg.Foreplay()
//
// 			reg.SetInstantMode()
//
// 			regInfo := reg.GetRegistrarInfo()
// 			fmt.Printf("Reqistrar %s adapter count is %d, serial number %d, network address %d\n",
// 				reg.ip, regInfo.AdapterCount, regInfo.SerialNumber, regInfo.NetworkAddress)
//
// 			reg.meterCount = regInfo.AdapterCount
//
// 			// reg.Foreplay2()
//
// 			// adapter := uint16(46)
// 			// m := reg.GetAdapterIndications(adapter)
// 			// fmt.Printf("Adapter %d indications: %v\n", adapter, m)
//
// 			reg.Foreplay()
//
// 			reg.SetInstantMode()
//
// 			reg.SetInstantMode()
//
// 			m := reg.GetAllIndications()
// 			// fmt.Printf("Adapter's indications: %v\n", m)
//
// 			metrics.SetMetersIndications(m)
//
// 			fmt.Printf("Reqistrar %s poll successfuly finished\n", reg.ip)
// 		}(reg)
// 	}
//
// 	wg.Wait()
//
// 	fmt.Printf("Poll successfuly finished\n")
// 	onReady()
// }
