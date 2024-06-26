package emulator

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/lab5e/spangw/pkg/gw"
	"github.com/lab5e/spangw/pkg/stdgw"
)

// New creates a new command handler that emulates devices on a gateway
func New() gw.CommandHandler {
	ret := &emulatorHandler{
		devices: make([]device, 0),
		mutex:   &sync.Mutex{},
	}
	go ret.generateUpstream()
	return ret
}

type device struct {
	id       string
	config   map[string]string
	messages []string
}

type emulatorHandler struct {
	gatewayConfig map[string]string
	devices       []device
	deviceCount   int64
	mutex         *sync.Mutex
	upstreamCb    gw.UpstreamMessageFunc
}

func (e *emulatorHandler) UpdateConfig(localID string, config map[string]string) (string, error) {
	e.gatewayConfig = config
	slog.Info("Updated gateway config", "config", config)
	if localID == "" {
		return "1", nil
	}
	return localID, nil
}

func (e *emulatorHandler) RemoveDevice(localID string, deviceID string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	for i, d := range e.devices {
		if d.id == deviceID {
			e.devices = append(e.devices[:i], e.devices[i+1:]...)
			slog.Info("removed device", "localID", localID)
			return nil
		}
	}
	return nil
}

func (e *emulatorHandler) UpdateDevice(localID string, localDeviceID string, config map[string]string) (string, map[string]string, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if localDeviceID == "" {
		e.deviceCount++
		d := device{
			id:       strconv.FormatInt(e.deviceCount, 10),
			config:   config,
			messages: make([]string, 0),
		}
		e.devices = append(e.devices, d)

		config[stdgw.LoraFCntUp] = "99"
		config[stdgw.LoraFCntDn] = "9"
		slog.Info("Added device", "device", d)
		return d.id, config, nil
	}

	for i, d := range e.devices {
		if d.id == localID {
			e.devices[i].config = config
			slog.Info("Updated device", "device", d)
			return d.id, nil, nil
		}
	}
	return localDeviceID, nil, nil
}

func (e *emulatorHandler) DownstreamMessage(localID, localDeviceID, messageID string, payload []byte) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	for i, d := range e.devices {
		if d.id == localDeviceID {
			slog.Info("Got downstream message", "deviceID", d.id, "payloadLength", len(payload))
			e.devices[i].messages = append(e.devices[i].messages, hex.EncodeToString(payload))
			return nil
		}
	}
	return nil
}

func (e *emulatorHandler) UpstreamMessage(upstreamCb gw.UpstreamMessageFunc) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.upstreamCb = upstreamCb
}

func (e *emulatorHandler) generateUpstream() {
	count := 1
	for {
		time.Sleep(30 * time.Second)
		e.mutex.Lock()
		if e.upstreamCb != nil && len(e.devices) > 0 {
			ix := rand.Intn(len(e.devices))
			slog.Info("Generating upstream message for device", "id", e.devices[ix].id)
			e.upstreamCb(e.devices[ix].id, []byte(fmt.Sprintf("msg %d", count)), map[string]string{
				"rssi":   strconv.FormatInt(int64(count), 10),
				"fCntUp": strconv.FormatInt(int64(count), 10),
			})
			count++
		}
		e.mutex.Unlock()
	}
}

func (e *emulatorHandler) Shutdown() {
	// nothing
}
