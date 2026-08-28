package mqtt

import (
	"github.com/kuzmichus/neoclaw/pkg/bus"
	"github.com/kuzmichus/neoclaw/pkg/channels"
	"github.com/kuzmichus/neoclaw/pkg/config"
)

func init() {
	channels.RegisterSafeFactory(
		config.ChannelMQTT,
		func(bc *config.Channel, cfg *config.MQTTSettings, b *bus.MessageBus) (channels.Channel, error) {
			return NewMQTTChannel(bc, cfg, b)
		},
	)
}
