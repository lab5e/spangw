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

## Test the gateway emulator

This project includes a simple gateway implementation (see the `emulator` package) that 
generates upstream messages for a random device every 30 seconds. 

Create a new collection for your device(s) and gateway via the [Span CLI](https://github.com/lab5e/spancli):

```shell
span col add --tag name:"Gateway demo collection"
```

Add a new gateway on the collection:

```shell
span gw add --name "Emulated gateway" --collection-id={collection id from above}
```

Add a device on the collection and configure it to use the gateway. The configuration properties 
depends on your gateway's needs:

```shell
span dev add --collection-id={collection ID} --gateway-id={gateway ID} --config device:no1 --config foo:bar --config bar:baz
```

Finally, create a certificate for your gateway:

```shell
span cert create --gateway-id=17mjfma79c872g
```

You can now observe the activity events to make sure the gateway connects to Span:

```shell
span activitiy watch --collection-id={collection ID}
```

You can monitor the upstream messages with

```shell
span inbox watch --collection-id={collection ID}
```

Launch the gateway emulator with

```shell
bin/gwemu  --cert-file=cert.crt --key-file=key.pem --state-file=state.json
```

If you inspect the device with `span device get --device-id={device ID} --collection-id={collection ID}` you should see 
an updated device config (the fCntDn and fCntUp fields):

```shell
 Device {device ID}
 Field                                        Value
 lastGatewayId                                {device ID}
 lastTransport                                gateway
 metadata.gateway.gatewayId                   {gateway ID}
 metadata.gateway.lastUpdate                  1715070179623
 collectionId                                 {collection ID}
 config.gateway.17mjfma79c872g.gatewayId      {gateway ID}
 config.gateway.17mjfma79c872g.params.foo     bar
 config.gateway.17mjfma79c872g.params.bar     baz
 config.gateway.17mjfma79c872g.params.device  no1
 config.gateway.17mjfma79c872g.params.fCntDn  9
 config.gateway.17mjfma79c872g.params.fCntUp  99
 firmware.manufacturer
 firmware.modelNumber
 firmware.serialNumber
 firmware.state                               Current
 firmware.stateMessage
 firmware.targetFirmwareId                    0
 firmware.currentFirmwareId                   0
 firmware.firmwareVersion
 tags.name                                    The demo device
 deviceId                                     {device ID}
 lastPayload                                  bXNnIDM=
 lastReceived                                 1715070269492894739
 Config {gateway ID}                          bar:baz,device:no1,fCntDn:9,fCntUp:99,foo:bar
```

