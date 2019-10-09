package service

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/linux"
)

var (
	device = flag.String("device", "default", "implementation of ble")
	du     = flag.Duration("du", 300*time.Second, "advertising duration, 0 for indefinitely")
)

type GATTServer struct {
}

// Start yay
func (*GATTServer) Start(s *Server) {
	flag.Parse()

	d, err := linux.NewDevice()
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}
	ble.SetDefaultDevice(d)

	testSvc := ble.NewService(TestSvcUUID)
	testSvc.AddCharacteristic(NewPlayChar(s))
	testSvc.AddCharacteristic(NewVolumeChar(s))
	testSvc.AddCharacteristic(NewLoopChar(s))
	testSvc.AddCharacteristic(NewEQChar(s))
	testSvc.AddCharacteristic(NewMusicChar(s))
	testSvc.AddCharacteristic(NewTestChar(s))

	if err := ble.AddService(testSvc); err != nil {
		log.Fatalf("can't add service: %s", err)
	}

	// Advertise for specified durantion, or until interrupted by user.
	fmt.Printf("Advertising for %s...\n", *du)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *du))
	chkErr(ble.AdvertiseNameAndServices(ctx, "KPlay", testSvc.UUID), s)
}

func chkErr(err error, s *Server) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("done\n")
	case context.Canceled:
		fmt.Printf("canceled\n")
	default:
		log.Fatalf(err.Error())
	}

}
