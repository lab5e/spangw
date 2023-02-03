# Span Gateway library

This is a library that implements the boilerplate gateway code for Span. Implement the `CommandHandler`
interface to add gateway functionality.

There's a sample implementation that will emulate a gateway in the `emulator` package.

The gateway can be launched with a few lines of code: 

```golang
func main() {
	var config gw.Parameters    
    // ... set parameters here

    var myGatewayHandler gw.CommandHandler
    // ..create handler here

	gwHandler, err := gw.Create(config, myGatewayHandler)
	if err != nil {
		log.Printf("Error creating gateway: %v", err)
		os.Exit(2)
	}

	defer gwHandler.Stop()
	if err := gwHandler.Run(); err != nil {
		log.Printf("Could not run the gateway process: %v", err)
		os.Exit(2)
	}
}
```

## Building

Install dependencies with

`make tools``

Build the sample service with

`make``

