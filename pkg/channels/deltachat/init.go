package deltachat

import (
	"github.com/kuzmichus/neoclaw/pkg/bus"
	"github.com/kuzmichus/neoclaw/pkg/channels"
	"github.com/kuzmichus/neoclaw/pkg/config"
)

func init() {
	channels.RegisterFactory(
		config.ChannelDeltaChat,
		func(channelName, channelType string, cfg *config.Config, b *bus.MessageBus) (channels.Channel, error) {
			bc := cfg.Channels[channelName]
			if bc == nil || !bc.Enabled {
				return nil, nil
			}
			decoded, err := bc.GetDecoded()
			if err != nil {
				return nil, err
			}
			c, ok := decoded.(*config.DeltaChatSettings)
			if !ok {
				return nil, channels.ErrSendFailed
			}
			ch, err := NewDeltaChatChannel(bc, c, b)
			if err != nil {
				return nil, err
			}
			if channelName != config.ChannelDeltaChat {
				ch.SetName(channelName)
			}
			return ch, nil
		},
	)
}
