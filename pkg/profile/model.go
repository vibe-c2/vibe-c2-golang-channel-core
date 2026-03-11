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
	ProfileID        *MapField `json:"profile_id,omitempty" yaml:"profile_id,omitempty"`
	ID               MapField  `json:"id" yaml:"id"`
	EncryptedDataIn  MapField  `json:"encrypted_data_in" yaml:"encrypted_data_in"`
	EncryptedDataOut MapField  `json:"encrypted_data_out" yaml:"encrypted_data_out"`
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

func (m MapField) Ref() string {
	loc := strings.TrimSpace(m.Target.Location)
	key := strings.TrimSpace(m.Target.Key)
	if loc == "" {
		return key
	}
	return loc + ":" + key
}
