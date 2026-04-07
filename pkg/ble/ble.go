package ble

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/pkg/errors"
)

var (
	device = flag.String("device", "default", "implementation of ble")
	du     = flag.Duration("du", 5*time.Second, "scanning duration")
)

var deviceMap map[string]ble.Advertisement

var companyNames map[uint16]string = map[uint16]string{
	0x004c: "Apple",
	0x0006: "Microsoft",
	0x055d: "Valve",
	0x012e: "Assa Abloy",
	0x02e1: "Victron Energy",
}

func Scan() (string, error) {
	deviceMap = make(map[string]ble.Advertisement)
	flag.Parse()

	d, err := dev.NewDevice(*device)
	if err != nil {
		fmt.Printf("can't new device : %s", err)
		return "", err
	}
	ble.SetDefaultDevice(d)

	// Scan for specified durantion, or until interrupted by user.
	fmt.Printf("Scanning for %s...\n", *du)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *du))
	chkErr(ble.Scan(ctx, false, advHandler, nil))
	return printDevices(), nil
}

func advHandler(a ble.Advertisement) {
	deviceMap[a.Addr().String()] = a
}

func printDevices() string {
	sb := strings.Builder{}
	devices := make([]ble.Advertisement, 0)

	for _, a := range deviceMap {
		devices = append(devices, a)
	}

	slices.SortFunc(devices[:], func(a, b ble.Advertisement) int {
		return b.RSSI() - a.RSSI()
	})

	for _, a := range devices {
		if a.Connectable() {
			fmt.Fprintf(&sb, "[%s] C %3d:", a.Addr(), a.RSSI())
		} else {
			fmt.Fprintf(&sb, "[%s] N %3d:", a.Addr(), a.RSSI())
		}
		comma := ""
		if len(a.LocalName()) > 0 {
			fmt.Fprintf(&sb, " Name: %s", a.LocalName())
			comma = ","
		}
		if len(a.Services()) > 0 {
			fmt.Fprintf(&sb, "%s Svcs: %v", comma, a.Services())
			comma = ","
		}
		if len(a.ManufacturerData()) > 0 {
			fmt.Fprintf(&sb, "%s MD: %X", comma, a.ManufacturerData())
		}
		if company := getCompanyName(a); len(company) > 0 {
			fmt.Fprintf(&sb, ", %s", company)
		}
		fmt.Fprintf(&sb, "\n")
	}
	return sb.String()
}

func getCompanyName(a ble.Advertisement) string {
	d := a.ManufacturerData()
	if len(d) < 2 {
		return ""
	}
	name, ok := companyNames[binary.LittleEndian.Uint16(d)]
	if !ok {
		return ""
	}
	return name
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("done\n")
	case context.Canceled:
		fmt.Printf("canceled\n")
	default:
		fmt.Printf("Error during scan: %s\n", err)
	}
}
