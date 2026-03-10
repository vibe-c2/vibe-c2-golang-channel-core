package profile

import (
	"fmt"
	"strings"

	coreerrors "github.com/vibe-c2/vibe-c2-golang-channel-core/pkg/errors"
)

// Validate performs baseline semantic checks for a single profile.
func Validate(p Profile) error {
	if strings.TrimSpace(p.ProfileID) == "" {
		return coreerrors.New(coreerrors.CodeProfileInvalid, "profile_id is required")
	}
	if strings.TrimSpace(p.ChannelType) == "" {
		return coreerrors.New(coreerrors.CodeProfileInvalid, "channel_type is required")
	}
	if strings.TrimSpace(p.Mapping.ID) == "" {
		return coreerrors.New(coreerrors.CodeProfileInvalid, "mapping.id is required")
	}
	if strings.TrimSpace(p.Mapping.EncryptedData) == "" {
		return coreerrors.New(coreerrors.CodeProfileInvalid, "mapping.encrypted_data is required")
	}
	return nil
}

// ValidateMany validates all profiles and annotates list index on failure.
func ValidateMany(profiles []Profile) error {
	for i := range profiles {
		if err := Validate(profiles[i]); err != nil {
			return coreerrors.Wrap(coreerrors.CodeProfileInvalid, fmt.Sprintf("profile[%d] invalid", i), err)
		}
	}
	return nil
}
