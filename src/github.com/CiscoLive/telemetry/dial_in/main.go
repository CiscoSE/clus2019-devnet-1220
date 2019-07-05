/*
Dial-in to specified router in specified subscription
*/

package main

import (
	"fmt"
	// Step 2 imports
	xr "github.com/nleiva/xrgrpc"
	"log"
	"html/template"
	"bytes"
	"os"



	// Step 3 imports
	"os/signal"
	"context"


	// Step 4 imports
	"github.com/golang/protobuf/proto"
	interfaces "github.com/CiscoLive/telemetry/proto/if_generic_counters"
	"github.com/nleiva/xrgrpc/proto/telemetry"


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
	// Create a router instance with the connection parameters
	router1, err := xr.BuildRouter(
		xr.WithUsername("vagrant"),
		xr.WithPassword("vagrant"),
		xr.WithHost("192.0.2.2:57344"),
		xr.WithTimeout(60),
	)
    // Check for any errors
	if err != nil {
		log.Fatalf("Target parameters for router are incorrect: %s", err)
	}
    // Set the ID for first the transaction. Each transaction needs to have an unique ID.
    var id int64 = 1000

	// Define Telemetry parameters
	tConfigInterfaces := &TelemetryConfig{
		SensorGroupID:         "ifcs-group",
		Path:          "Cisco-IOS-XR-infra-statsd-oper:infra-statistics/interfaces/interface/latest/generic-counters",
		SubscriptionID:     "interfaces-dial-in-subs",
		SampleInterval: 5000,
	}

	// Connect to router via gRPC
	grpcConnection, connectionContext, err := xr.Connect(*router1)
    // Check for any errors
	if err != nil {
		log.Printf("Could not setup a connection to %s, %v", router1.Host, err)

	}

	// Read the Telemetry template file to build the configuration to be sent
	telemetryTemplate, err := template.ParseFiles(os.Getenv("GOPATH") + "/src/github.com/CiscoLive/telemetry/oc-templates/oc-telemetry.json")

	// 'templateBuf' is an io.Writter to capture the template execution output for the device
	templateBuf := new(bytes.Buffer)
	err = telemetryTemplate.Execute(templateBuf, tConfigInterfaces)

    // Check for any errors
	if err != nil {
		log.Printf("Could not execute telemetry config template for router: %v", err)
		return
	}

	// Apply the template+parameters to the router using MergeConfig RPC call
	ri, err := xr.MergeConfig(connectionContext, grpcConnection, templateBuf.String(), id)
    // Check for any errors
	if err != nil {
		log.Fatalf("Failed to config %s: %v\n", router1.Host, err)
		return
	} else {
        // No errors found. Print message details
		fmt.Printf("\nConfig merged on %s -> Request ID: %v, Response ID: %v\n\n", router1.Host, id, ri)
	}

    // Print success message
	fmt.Println("Configuration done! Press the return to continue")
    // Wait for user to press return
    fmt.Scanln()


	

	// ************* End Step 2

	
	// ************* Step 3 - Create subscription


    // Increase transaction ID in one
    id++

    // Create connection context
	ctx1, cancel := context.WithCancel(connectionContext)

    // Create channel use to monitor exit signals and errors
	c := make(chan os.Signal, 1)

	// Relay to this channel all interrupt signals
	signal.Notify(c, os.Interrupt)

	// Set telemetry subscription name to request messages from.
	p := "interfaces-dial-in-subs"

    // Set encoding type 2 (Compact Google Protocol Buffers - cGPB)
	var e int64 = 2

    // Execute get subscription RPC in device
	ch, ech, err := xr.GetSubscription(ctx1, grpcConnection, p, id, e)

    // Check for errors
	if err != nil {
		log.Fatalf("Could not setup Telemetry Subscription: %v\n", err)
	}

    // Listen for errors or interrupt signals
	go func() {
		select {
		case <-c:
            // Cancel request (e.g. ctrl+c) from user
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

    // Print connection established message
	fmt.Printf("\nConnected to  %s\n\n", router1.Host)

        // Listen and process all messages comming from the device via "ch"
		for tele := range ch {

            // Print raw message
			log.Printf("***** New message from %v ***** \n", router1.Host)
			fmt.Printf("Undecoded Message: %v\n\n", tele)

            // ************* Step 4 - Decode message content

        // Create a telemetry instance to store decoded information in memory
		message := new(telemetry.Telemetry)

        // Decode/Unmarshal the message headers
		err = proto.Unmarshal(tele, message)

        // Check for any errors
		if err != nil {
			log.Printf("Could not unmarshall the interface telemetry message for %v: %v\n", router1.Host, err)
		}

        // Print results
		fmt.Printf("======= Decoded Message ===========\n\n")

        // Traverse all rows within the message
		for _, row := range message.GetDataGpb().GetRow() {
			fmt.Printf("======= New Row \n\n")
			// Get the message content
			content := row.GetContent()

            // Create a IfstatsbagGeneric instance to store decoded information in memory
			ifBag := new(interfaces.IfstatsbagGeneric)

            // Decode/Unmarshal the message content
			err = proto.Unmarshal(content, ifBag)

            // Check for errors
			if err != nil {
				log.Fatalf("Could decode Content: %v\n", err)
			}

            // Print content
			fmt.Printf("Content: %v\n\n", ifBag)

			// ************* Step 5 - Decode message keys

        // Get the keys for this specific section of the message
		keys := row.GetKeys()

        // Create a IfstatsbagGeneric_KEYS instance to store decoded keys in memory
        ifBagKeys := new(interfaces.IfstatsbagGeneric_KEYS)

        // Decode/Unmarshal message keys
        err = proto.Unmarshal(keys, ifBagKeys)

        // Check for any errors
        if err != nil {
            log.Fatalf("Could decode keys: %v\n", err)
        }

        // Print Result
        fmt.Printf("Row Keys: %v\n\n", ifBagKeys)



			// ************* End Step 5
            // ************* Step 6 - Get specific data

            // Get bytes received from message content
            bytesRcvd := ifBag.GetBytesReceived()

            // Get interface name from message Keys
			interfaceName := ifBagKeys.GetInterfaceName()

            // Print interface name and bytes received
			fmt.Printf("Interface %v Bytes received %v\n", interfaceName, bytesRcvd)


			// ************* End Step 6

            // Print row separator
			fmt.Printf("======= \n\n")
		}



			// ************* End Step 4

		}




	// ************* End Step 3
}
