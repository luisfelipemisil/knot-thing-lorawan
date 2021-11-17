package knot

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	pb "github.com/brocaar/chirpstack-api/go/v3/as/integration"
	"github.com/brocaar/chirpstack-application-server/internal/config"
	"github.com/brocaar/chirpstack-application-server/internal/integration/knot/entities"
	"github.com/brocaar/chirpstack-application-server/internal/integration/knot/network"
	"github.com/brocaar/chirpstack-application-server/internal/integration/marshaler"
	"github.com/brocaar/chirpstack-application-server/internal/integration/models"
)

// Integration implements an KNoT integration.
type Integration struct {
	protocol Protocol
}

var deviceChan = make(chan entities.Device)
var msgChan = make(chan network.InMsg)

// New creates a new KNoT integration.
func New(m marshaler.Type, conf config.IntegrationKNoTConfig) (*Integration, error) {
	var err error
	i := Integration{}

	i.protocol, err = newProtocol(conf, deviceChan, msgChan)
	if err != nil {
		return nil, errors.Wrap(err, "new knot protocol")
	}

	return &i, nil
}

// Formatting all the information needed to configure a knot device
func formatDevice(DevEui []byte, deviceName string, config map[string]string) entities.Device {
	device := entities.Device{}
	DevEUI_str := []byte("")
	configFrame := entities.Config{}
	var str string
	//get all settings from all sensors
	for i := 1; i < 100; i++ {
		str = fmt.Sprintf("sensor%d", i)
		if confData, ok := config[str]; ok {
			json.Unmarshal([]byte(string(confData)), &configFrame)
			device.Config = append(device.Config, configFrame)
		} else {
			break
		}
	}
	fmt.Println(device.Config)

	for _, v := range DevEui {
		DevEUI_str = strconv.AppendInt(DevEUI_str, int64(v), 16)
	}
	device.ID = string(DevEUI_str)
	device.Name = deviceName
	return device
}

// HandleUplinkEvent sends an UplinkEvent.
func (i *Integration) HandleUplinkEvent(ctx context.Context, _ models.Integration, vars map[string]string, pl pb.UplinkEvent) error {

	deviceChan <- formatDevice(pl.DevEui, pl.DeviceName, pl.Tags)

	log.WithFields(log.Fields{"event": "uplink"}).Info("New uplink")

	return nil
}

// HandleJoinEvent sends a JoinEvent.
func (i *Integration) HandleJoinEvent(ctx context.Context, _ models.Integration, vars map[string]string, pl pb.JoinEvent) error {

	deviceChan <- formatDevice(pl.DevEui, pl.DeviceName, pl.Tags)

	log.WithFields(log.Fields{"event": "join"}).Info("New join")

	return nil
}

// HandleAckEvent sends an AckEvent.
func (i *Integration) HandleAckEvent(ctx context.Context, _ models.Integration, vars map[string]string, pl pb.AckEvent) error {
	return nil
}

// HandleErrorEvent sends an ErrorEvent.
func (i *Integration) HandleErrorEvent(ctx context.Context, _ models.Integration, vars map[string]string, pl pb.ErrorEvent) error {
	return nil
}

// HandleStatusEvent sends a StatusEvent.
func (i *Integration) HandleStatusEvent(ctx context.Context, _ models.Integration, vars map[string]string, pl pb.StatusEvent) error {
	return nil
}

// HandleLocationEvent sends a LocationEvent.
func (i *Integration) HandleLocationEvent(ctx context.Context, _ models.Integration, vars map[string]string, pl pb.LocationEvent) error {
	return nil
}

// HandleTxAckEvent sends a TxAckEvent.
func (i *Integration) HandleTxAckEvent(ctx context.Context, _ models.Integration, vars map[string]string, pl pb.TxAckEvent) error {
	return nil
}

// HandleIntegrationEvent sends an IntegrationEvent.
func (i *Integration) HandleIntegrationEvent(ctx context.Context, _ models.Integration, vars map[string]string, pl pb.IntegrationEvent) error {
	return nil
}

// DataDownChan returns nil
func (i *Integration) DataDownChan() chan models.DataDownPayload {
	return nil
}

// Close closes the integration.
func (i *Integration) Close() error {
	return i.protocol.Close()
}
