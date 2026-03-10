package protocol

import (
	"fmt"
	"strings"
)

// InboundAgentMessage is the canonical request format sent to C2.
type InboundAgentMessage struct {
	ChannelID     string `json:"channel_id" yaml:"channel_id"`
	ProfileID     string `json:"profile_id,omitempty" yaml:"profile_id,omitempty"`
	ID            string `json:"id" yaml:"id"`
	EncryptedData string `json:"encrypted_data" yaml:"encrypted_data"`
}

// OutboundAgentMessage is the canonical response format received from C2.
type OutboundAgentMessage struct {
	ChannelID     string `json:"channel_id" yaml:"channel_id"`
	ProfileID     string `json:"profile_id,omitempty" yaml:"profile_id,omitempty"`
	ID            string `json:"id" yaml:"id"`
	EncryptedData string `json:"encrypted_data" yaml:"encrypted_data"`
}

func ValidateInboundAgentMessage(m InboundAgentMessage) error {
	if strings.TrimSpace(m.ChannelID) == "" {
		return fmt.Errorf("inbound channel_id is required")
	}
	if strings.TrimSpace(m.ID) == "" {
		return fmt.Errorf("inbound id is required")
	}
	if strings.TrimSpace(m.EncryptedData) == "" {
		return fmt.Errorf("inbound encrypted_data is required")
	}
	return nil
}

func ValidateOutboundAgentMessage(m OutboundAgentMessage) error {
	if strings.TrimSpace(m.ChannelID) == "" {
		return fmt.Errorf("outbound channel_id is required")
	}
	if strings.TrimSpace(m.ID) == "" {
		return fmt.Errorf("outbound id is required")
	}
	if strings.TrimSpace(m.EncryptedData) == "" {
		return fmt.Errorf("outbound encrypted_data is required")
	}
	return nil
}
