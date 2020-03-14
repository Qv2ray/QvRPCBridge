package main

//go:generate protoc --go_out=plugins=grpc:. command/v2ray_api_commands.proto

import "C" //required
import (
	"context"
	"fmt"
	cmd "github.com/Qv2ray/QvRPCBridge/command"
	"google.golang.org/grpc"
	"time"
)

//required
func main() {}

var client *grpc.ClientConn

//export Dial
func Dial(addr *C.char, timeout uint32) (errMsg *C.char) {
	return C.CString(_Dial(C.GoString(addr), timeout))
}

//export GetStats
func GetStats(name *C.char, timeout uint32) (value int64) {
	return _GetStats(C.GoString(name), timeout)
}

func _Dial(addr string, timeout uint32) (errMsg string) {
	if client != nil {
		_ = client.Close()
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	errPipe := make(chan error)
	goodPipe := make(chan struct{})
	go func() {
		conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
		select {
		case <-ctx.Done():
			return
		default:
			if err != nil {
				errPipe <- err
				return
			}
			client = conn
			close(goodPipe)
		}
	}()
	select {
	case <-time.After(time.Duration(timeout) * time.Millisecond):
		cancelFunc()
		return fmt.Sprintf("dial timed out %dms at addr `%s`", timeout, addr)
	case err := <-errPipe:
		return fmt.Sprintf("dial failed: %s", err.Error())
	case <-goodPipe:
		return ""
	}
}

func _GetStats(name string, timeout uint32) (value int64) {
	if client == nil {
		return -999
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	errChan := make(chan error)
	goodChan := make(chan int64)
	go func() {
		stats, err := cmd.NewStatsServiceClient(client).GetStats(ctx, &cmd.GetStatsRequest{
			Name:   name,
			Reset_: true,
		})
		select {
		case <-ctx.Done():
			return
		default:
			if err != nil {
				errChan <- err
				return
			}
			value = stats.Stat.Value
			close(goodChan)
		}
	}()
	select {
	case <-time.After(time.Duration(timeout) * time.Millisecond):
		cancelFunc()
		return -666
	case <-errChan:
		return -1
	case <-goodChan:
		return
	}
}
