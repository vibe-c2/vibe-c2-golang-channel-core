package runtime

import (
	"context"
	"strings"

	coreerrors "github.com/vibe-c2/vibe-c2-golang-channel-core/pkg/errors"
	protocol "github.com/vibe-c2/vibe-c2-golang-protocol"
)

// Runtime orchestrates canonical message handling for channel adapters.
type Runtime struct {
	SyncClient SyncClient
}

func New(syncClient SyncClient) *Runtime {
	return &Runtime{SyncClient: syncClient}
}

// Handle converts envelope fields into canonical inbound message, validates it,
// syncs with C2, validates outbound canonical message, and writes mapped fields back.
func (r *Runtime) Handle(ctx context.Context, envelope TransportEnvelope, channelID string) (protocol.OutboundAgentMessage, error) {
	if r == nil || r.SyncClient == nil {
		return protocol.OutboundAgentMessage{}, coreerrors.New(coreerrors.CodeInvalidInput, "runtime sync client is required")
	}
	if envelope == nil {
		return protocol.OutboundAgentMessage{}, coreerrors.New(coreerrors.CodeInvalidInput, "transport envelope is required")
	}
	if strings.TrimSpace(channelID) == "" {
		return protocol.OutboundAgentMessage{}, coreerrors.New(coreerrors.CodeInvalidInput, "channelID is required")
	}

	inbound, err := inboundFromEnvelope(envelope, channelID)
	if err != nil {
		return protocol.OutboundAgentMessage{}, err
	}
	if err := protocol.ValidateInboundAgentMessage(inbound); err != nil {
		return protocol.OutboundAgentMessage{}, coreerrors.Wrap(coreerrors.CodeCanonicalInvalid, "invalid inbound canonical message", err)
	}

	outbound, err := r.SyncClient.Sync(ctx, inbound)
	if err != nil {
		return protocol.OutboundAgentMessage{}, err
	}
	if err := protocol.ValidateOutboundAgentMessage(outbound); err != nil {
		return protocol.OutboundAgentMessage{}, coreerrors.Wrap(coreerrors.CodeCanonicalInvalid, "invalid outbound canonical message", err)
	}

	envelope.SetField("mapping", "id", outbound.ID)
	envelope.SetField("mapping", "encrypted_data", outbound.EncryptedData)
	if outbound.ProfileID != "" {
		envelope.SetField("mapping", "profile_id", outbound.ProfileID)
	}

	return outbound, nil
}

func inboundFromEnvelope(envelope TransportEnvelope, channelID string) (protocol.InboundAgentMessage, error) {
	id, err := requiredEnvelopeField(envelope, "mapping", "id")
	if err != nil {
		return protocol.InboundAgentMessage{}, err
	}
	encryptedData, err := requiredEnvelopeField(envelope, "mapping", "encrypted_data")
	if err != nil {
		return protocol.InboundAgentMessage{}, err
	}

	profileID, _ := envelope.GetField("mapping", "profile_id")

	return protocol.InboundAgentMessage{
		ChannelID:     channelID,
		ProfileID:     strings.TrimSpace(profileID),
		ID:            id,
		EncryptedData: encryptedData,
	}, nil
}

func requiredEnvelopeField(envelope TransportEnvelope, location, key string) (string, error) {
	value, ok := envelope.GetField(location, key)
	if !ok || strings.TrimSpace(value) == "" {
		return "", coreerrors.New(coreerrors.CodeInvalidInput, "missing required field: "+location+"."+key)
	}
	return strings.TrimSpace(value), nil
}
