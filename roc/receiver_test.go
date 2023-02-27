package roc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func makeReceiverConfig() ReceiverConfig {
	return ReceiverConfig{
		FrameSampleRate:  44100,
		FrameChannels:    ChannelSetStereo,
		FrameEncoding:    FrameEncodingPcmFloat,
		ClockSource:      ClockInternal,
		ResamplerProfile: ResamplerProfileDisable,
	}
}

func TestReceiver_Open(t *testing.T) {
	tests := []struct {
		name    string
		config  ReceiverConfig
		wantErr error
	}{
		{
			name:    "ok",
			config:  makeReceiverConfig(),
			wantErr: nil,
		},
		{
			name:    "invalid config",
			config:  ReceiverConfig{},
			wantErr: newNativeErr("roc_receiver_open()", -1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := OpenContext(ContextConfig{})

			require.NoError(t, err)
			require.NotNil(t, ctx)

			receiver, err := OpenReceiver(ctx, tt.config)

			if tt.wantErr == nil {
				require.NoError(t, err)
				require.NotNil(t, receiver)

				err = receiver.Close()
				require.NoError(t, err)
			} else {
				require.Equal(t, tt.wantErr, err)
				require.Nil(t, receiver)
			}

			err = ctx.Close()
			require.NoError(t, err)
		})
	}
}

func TestReceiver_OpenWithNilContext(t *testing.T) {
	receiver, err := OpenReceiver(nil, makeReceiverConfig())

	require.Equal(t, errors.New("context is nil"), err)
	require.Nil(t, receiver)
}

func TestReceiver_OpenWithClosedContext(t *testing.T) {
	ctx, err := OpenContext(ContextConfig{})
	ctx.Close()
	receiver, err := OpenReceiver(ctx, makeReceiverConfig())

	require.Equal(t, errors.New("context is closed"), err)
	require.Nil(t, receiver)
}

func TestReceiver_BindErrors(t *testing.T) {
	tests := []struct {
		name     string
		slot     Slot
		iface    Interface
		endpoint *Endpoint
		wantErr  error
	}{
		{
			name:    "nil endpoint",
			slot:    SlotDefault,
			iface:   InterfaceAudioSource,
			wantErr: errors.New("endpoint is nil"),
		},
		{
			name:     "invalid endpoint",
			slot:     SlotDefault,
			iface:    InterfaceAudioSource,
			endpoint: &Endpoint{},
			wantErr:  newNativeErr("roc_endpoint_set_protocol()", -1),
		},
		{
			name:     "incompatible iface and endpoint",
			slot:     SlotDefault,
			iface:    InterfaceAudioControl,
			endpoint: &Endpoint{Protocol: ProtoRtsp, Host: "192.168.0.1"},
			wantErr:  newNativeErr("roc_receiver_bind()", -1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := OpenContext(ContextConfig{})
			require.NoError(t, err)
			require.NotNil(t, ctx)

			receiver, err := OpenReceiver(ctx, makeReceiverConfig())
			require.NoError(t, err)
			require.NotNil(t, receiver)

			err = receiver.Bind(tt.slot, tt.iface, tt.endpoint)

			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestReceiver_Closed(t *testing.T) {
	tests := []struct {
		name      string
		operation func(receiver *Receiver) error
	}{
		{
			name: "SetMulticastGroup after close",
			operation: func(receiver *Receiver) error {
				return receiver.SetMulticastGroup(SlotDefault, InterfaceAudioSource, "127.0.0.1")
			},
		},
		{
			name: "Bind after close",
			operation: func(receiver *Receiver) error {
				return receiver.Bind(SlotDefault, InterfaceAudioSource, nil)
			},
		},
		{
			name: "ReadFloats after close",
			operation: func(receiver *Receiver) error {
				recFloats := make([]float32, 2)
				return receiver.ReadFloats(recFloats)
			},
		},
	}
	for _, tt := range tests {
		ctx, err := OpenContext(ContextConfig{})
		require.NoError(t, err)
		require.NotNil(t, ctx)

		receiver, err := OpenReceiver(ctx, makeReceiverConfig())
		require.NoError(t, err)
		require.NotNil(t, receiver)

		err = receiver.Close()
		require.NoError(t, err)

		require.Equal(t, errors.New("receiver is closed"), tt.operation(receiver))
	}
}
