/*
Tests
*/

package main

import (
	"fmt"
	// Step 2 imports
	xr "github.com/nleiva/xrgrpc"
	"log"
	"os"



	// Step 3 imports
    "os/signal"
	"context"


)

type TelemetryConfig struct {
	SensorGroupID  string
	Path           string
	SubscriptionID string
	SampleInterval int
}

func main() {  
	fmt.Print("\nWelcome to the telemetry app\n")
	// ************* Step 2 - Configure telemetry
	router1, err := xr.BuildRouter(
		xr.WithUsername("vagrant"),
		xr.WithPassword("vagrant"),
		xr.WithHost("192.0.2.2:57344"),
		xr.WithCert("../ems.pem"),
		xr.WithTimeout(60),
	)
	if err != nil {
		log.Fatalf("Target parameters for router are incorrect: %s", err)
	}
    // Determine the ID for first the transaction.
    var id int64 = 1000
	// Connect to router
	grpcConnection, connectionContext, err := xr.Connect(*router1)
	if err != nil {
		log.Printf("Could not setup a connection to %s, %v", router1.Host, err)

	}

	ctx1, cancel := context.WithCancel(connectionContext)

	c := make(chan os.Signal, 1)

	// If no signals are provided, all incoming signals will be relayed to c.
	// Otherwise, just the provided signals will. E.g.: signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt)



	// Telemetry Subscription Name
	p := "interfaces-dial-in-subs-test"

	var e int64 = 2
	_, ech, err := xr.GetSubscription(ctx1, grpcConnection, p, id, e)
	if err != nil {
		log.Fatalf("Could not setup Telemetry Subscription: %v\n", err)
	}

	go func() {
		select {
		case <-c:
			fmt.Printf("\nManually cancelled the session to %v\n\n", router1.Host)
			cancel()
			return
		case <-ctx1.Done():
			// Timeout: "context deadline exceeded"
			err = ctx1.Err()
			fmt.Printf("\ngRPC session timed out after %v seconds: %v\n\n", router1.Timeout, err.Error())
			return
		case err = <-ech:
			// Session canceled: "context canceled"
			fmt.Printf("\ngRPC session to %v failed: %v\n\n", router1.Host, err.Error())
			return
		}
	}()

	fmt.Printf("\nConnected to  %s\n\n", router1.Host)
}
