package gw

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/lab5e/spangw/pkg/pb/gateway/v1"
)

const keepAliveInterval = time.Second * 60

// GatewayProcess is the command processor for gateways
type GatewayProcess struct {
	Stream           gateway.UserGatewayService_ControlStreamClient
	Commands         CommandHandler
	stateFile        string
	upstreamRequests chan *gateway.ControlStreamRequest
}

// NewGatewayProcess creates a new stream processor
func NewGatewayProcess(stateFile string, stream gateway.UserGatewayService_ControlStreamClient, commands CommandHandler) *GatewayProcess {
	return &GatewayProcess{
		Stream:    stream,
		Commands:  commands,
		stateFile: stateFile,
	}
}

// Run runs the stream processor. It will not return unless an error occurs
func (sp *GatewayProcess) Run() error {
	state, err := NewStateFromFile(sp.stateFile)
	if err != nil {
		return nil
	}
	downstreamResponses := make(chan *gateway.ControlStreamResponse)
	upstreamRequests := make(chan *gateway.ControlStreamRequest)
	errorCh := make(chan error)
	defer close(errorCh)
	defer sp.Commands.Shutdown()

	// Timestamp to keep track of last activity on the stream(s)
	lastMessage := time.Now()

	// Start two goroutines to handle the up- and downstream messages from Span. The upstream request
	// channel sends requests from the gateway to Span and the downstream response channel is used
	// to send messages back up to Span (messages, configuration changes and so on.)
	go func() {
		for msg := range upstreamRequests {
			err := sp.Stream.Send(msg)
			if err != nil {
				errorCh <- err
				return
			}
			lastMessage = time.Now()
		}
	}()

	go func() {
		defer close(downstreamResponses)
		for {
			msg, err := sp.Stream.Recv()
			if err != nil {
				errorCh <- err
				return
			}
			downstreamResponses <- msg
			lastMessage = time.Now()
		}
	}()

	// Send a configuration request as the first message to get the updated configuration
	upstreamRequests <- &gateway.ControlStreamRequest{
		Msg: &gateway.ControlStreamRequest_Config{},
	}

	sp.Commands.UpstreamMessage(func(localDeviceID string, payload []byte, metadata map[string]string) {
		deviceID := state.GetReverseMapping(localDeviceID)
		if deviceID == "" {
			slog.Error("Can't process upstream message. The local device ID is unnknown. Ignoring", "localID", localDeviceID)
			return
		}
		upstreamRequests <- &gateway.ControlStreamRequest{
			Msg: &gateway.ControlStreamRequest_UpstreamMessage{
				UpstreamMessage: &gateway.UpstreamMessage{
					DeviceId: deviceID,
					Payload:  payload,
					Metadata: metadata,
				},
			},
		}
	})
	for {
		select {
		case res := <-downstreamResponses:
			switch msg := res.Msg.(type) {
			case *gateway.ControlStreamResponse_KeepaliveResponse:
				break

			case *gateway.ControlStreamResponse_GatewayUpdate:
				state.GatewayID = msg.GatewayUpdate.GatewayId
				localID, err := sp.Commands.UpdateConfig(state.LocalID, msg.GatewayUpdate.Config)
				if err != nil {
					slog.Error("Error updating gateway configuration", "error", err)
					continue
				}
				state.LocalID = localID
				state.Save(sp.stateFile)

			case *gateway.ControlStreamResponse_DeviceRemoved:
				if err := sp.Commands.RemoveDevice(state.LocalID, state.GetMapping(msg.DeviceRemoved.DeviceId)); err != nil {
					slog.Error("Error removing device ", "deviceID", msg.DeviceRemoved.DeviceId, "error", err)
					continue
				}
				state.RemoveMapping(msg.DeviceRemoved.DeviceId)

			case *gateway.ControlStreamResponse_DeviceUpdate:
				if state.LocalID == "" {
					slog.Error("No local ID is set for the device. Will ignore the update command")
					continue
				}
				localDeviceID, newConfig, err := sp.Commands.UpdateDevice(state.LocalID, state.GetMapping(msg.DeviceUpdate.DeviceId), msg.DeviceUpdate.Config)
				if err != nil {
					slog.Error("Error updating device", "deviceID", msg.DeviceUpdate.DeviceId, "error", err)
					continue
				}
				state.SetMapping(msg.DeviceUpdate.DeviceId, localDeviceID)
				state.Save(sp.stateFile)
				if len(newConfig) > 0 {
					upstreamRequests <- &gateway.ControlStreamRequest{
						Msg: &gateway.ControlStreamRequest_DeviceUpdate{
							DeviceUpdate: &gateway.DeviceUpdate{
								DeviceId: msg.DeviceUpdate.DeviceId,
								Config:   newConfig,
							},
						},
					}
				}

			case *gateway.ControlStreamResponse_DownstreamMessage:
				if err := sp.Commands.DownstreamMessage(
					state.LocalID,
					state.GetMapping(msg.DownstreamMessage.DeviceId),
					msg.DownstreamMessage.MessageId,
					msg.DownstreamMessage.Payload); err != nil {
					slog.Error("Error sending message to device", "messageID", msg.DownstreamMessage.MessageId, "deviceID", msg.DownstreamMessage.DeviceId, "error", err)
				}

			default:
				slog.Warn("Unknown message from server", "type", fmt.Sprintf("%T", res.Msg))
			}

		case err := <-errorCh:
			if sp.upstreamRequests != nil {
				close(sp.upstreamRequests)
			}
			return err

		case <-time.After(10 * time.Second):
			// Check for timeout, send keepalive if time is > keepAliveInterval
			if time.Since(lastMessage) > keepAliveInterval {
				upstreamRequests <- &gateway.ControlStreamRequest{
					Msg: &gateway.ControlStreamRequest_Keepalive{},
				}
			}
		}
	}
}

// Stop closes the gateway process
func (sp *GatewayProcess) Stop() {
	sp.Stream.CloseSend()
}
