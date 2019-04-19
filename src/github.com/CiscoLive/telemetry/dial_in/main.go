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
	router1, err := xr.BuildRouter(
		xr.WithUsername("vagrant"),
		xr.WithPassword("vagrant"),
		xr.WithHost("192.0.2.2:57344"),
		xr.WithCert("ems.pem"),
		xr.WithTimeout(60),
	)
	if err != nil {
		log.Fatalf("Target parameters for router are incorrect: %s", err)
	}
    // Determine the ID for first the transaction.
    var id int64 = 1000

	// Define Telemetry parameters
	tConfigInterfaces := &TelemetryConfig{
		SensorGroupID:         "ifcs-group",
		Path:          "Cisco-IOS-XR-infra-statsd-oper:infra-statistics/interfaces/interface/latest/generic-counters",
		SubscriptionID:     "interfaces-dial-in-subs",
		SampleInterval: 5000,
	}

	// Connect to router
	grpcConnection, connectionContext, err := xr.Connect(*router1)
	if err != nil {
		log.Printf("Could not setup a connection to %s, %v", router1.Host, err)

	}


	// Read the OC Telemetry template file
	telemetryTemplate, err := template.ParseFiles(os.Getenv("GOPATH") + "/src/github.com/CiscoLive/telemetry/oc-templates/oc-telemetry.json")

	// 'templateBuf' is an io.Writter to capture the template execution output for the device
	templateBuf := new(bytes.Buffer)
	err = telemetryTemplate.Execute(templateBuf, tConfigInterfaces)
	if err != nil {
		log.Printf("Could not execute telemetry config template for router: %v", err)
		return
	}

	// Apply the template+parameters to the router.
	ri, err := xr.MergeConfig(connectionContext, grpcConnection, templateBuf.String(), id)
	if err != nil {
		log.Fatalf("Failed to config %s: %v\n", router1.Host, err)
		return
	} else {

		fmt.Printf("\nConfig merged on %s -> Request ID: %v, Response ID: %v\n\n", router1.Host, id, ri)
	}

	fmt.Println("Configuration done! Press the return to continue")
    fmt.Scanln()




	// ************* End Step 2

	
	// ************* Step 3 - Create subscription

    id++
	ctx1, cancel := context.WithCancel(connectionContext)

	c := make(chan os.Signal, 1)

	// If no signals are provided, all incoming signals will be relayed to c.
	// Otherwise, just the provided signals will. E.g.: signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt)



	// Telemetry Subscription Name
	p := "interfaces-dial-in-subs"

	var e int64 = 2
	ch, ech, err := xr.GetSubscription(ctx1, grpcConnection, p, id, e)
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

		for tele := range ch {
			log.Printf("***** New message from %v ***** \n", router1.Host)
			fmt.Printf("Undecoded Message: %v\n\n", tele)

            // ************* Step 4 - Decode message content


			message := new(telemetry.Telemetry)

			err = proto.Unmarshal(tele, message)
			if err != nil {
				log.Printf("Could not unmarshall the interface telemetry message for %v: %v\n", router1.Host, err)
			}
			fmt.Printf("======= Decoded Message ===========\n\n")
	
			for _, row := range message.GetDataGpb().GetRow() {
				fmt.Printf("======= New Row \n\n")
				// Get the message content
				content := row.GetContent()
				ifBag := new(interfaces.IfstatsbagGeneric)
				err = proto.Unmarshal(content, ifBag)
				if err != nil {
					log.Fatalf("Could decode Content: %v\n", err)
				}
	
				fmt.Printf("Content: %v\n\n", ifBag)
	
				// ************* Step 5 - Decode message keys
	

				keys := row.GetKeys()
				ifBagKeys := new(interfaces.IfstatsbagGeneric_KEYS)
				err = proto.Unmarshal(keys, ifBagKeys)
				if err != nil {
					log.Fatalf("Could decode keys: %v\n", err)
				}
		
				fmt.Printf("Row Keys: %v\n\n", ifBagKeys)
		
		
	
				// ************* End Step 5
				// ************* Step 6 - Get specific data
	            bytesRcvd := ifBag.GetBytesReceived()
			interfaceName := ifBagKeys.GetInterfaceName()
			fmt.Printf("Interface %v Bytes recieved %v\n", interfaceName, bytesRcvd)


				// ************* End Step 6
	
				fmt.Printf("======= \n\n")
			}
	
	

			// ************* End Step 4

		}



	// ************* End Step 3
}