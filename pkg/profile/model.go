package profile

import (
	"strings"

	"github.com/vibe-c2/vibe-c2-golang-channel-core/pkg/transform"
)

// Profile defines one obfuscation profile for a channel.
type Profile struct {
	ProfileID   string  `json:"profile_id" yaml:"profile_id"`
	ChannelType string  `json:"channel_type" yaml:"channel_type"`
	Enabled     bool    `json:"enabled" yaml:"enabled"`
	Priority    int     `json:"priority" yaml:"priority"`
	Version     int     `json:"version,omitempty" yaml:"version,omitempty"`
	Mapping     Mapping `json:"mapping" yaml:"mapping"`
}

type Mapping struct {
	ProfileID        *MapField      `json:"profile_id,omitempty" yaml:"profile_id,omitempty"`
	ID               MapField       `json:"id,omitempty" yaml:"id,omitempty"`
	EncryptedDataIn  MapField       `json:"encrypted_data_in,omitempty" yaml:"encrypted_data_in,omitempty"`
	EncryptedDataOut MapField       `json:"encrypted_data_out,omitempty" yaml:"encrypted_data_out,omitempty"`
	CombinedIn       *CombinedField `json:"combined_in,omitempty" yaml:"combined_in,omitempty"`
	CombinedOut      *CombinedField `json:"combined_out,omitempty" yaml:"combined_out,omitempty"`
}

type MapField struct {
	Source    string           `json:"source" yaml:"source"`
	Target    Target           `json:"target" yaml:"target"`
	Transform []transform.Spec `json:"transform,omitempty" yaml:"transform,omitempty"`
}

type Target struct {
	Location string `json:"location" yaml:"location"`
	Key      string `json:"key" yaml:"key"`
}

type CombinedField struct {
	Target    Target           `json:"target" yaml:"target"`
	Transform []transform.Spec `json:"transform,omitempty" yaml:"transform,omitempty"`
	Separator string           `json:"separator,omitempty" yaml:"separator,omitempty"`
}

func (c CombinedField) Ref() string {
	loc := strings.TrimSpace(c.Target.Location)
	key := strings.TrimSpace(c.Target.Key)
	if loc == "" {
		return key
	}
	return loc + ":" + key
}

func (m MapField) Ref() string {
	loc := strings.TrimSpace(m.Target.Location)
	key := strings.TrimSpace(m.Target.Key)
	if loc == "" {
		return key
	}
	return loc + ":" + key
}
