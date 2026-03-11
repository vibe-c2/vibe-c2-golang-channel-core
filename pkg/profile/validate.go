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
	if strings.TrimSpace(p.Mapping.ID.Ref()) == "" {
		return coreerrors.New(coreerrors.CodeProfileInvalid, "mapping.id target is required")
	}
	if strings.TrimSpace(p.Mapping.EncryptedDataIn.Ref()) == "" {
		return coreerrors.New(coreerrors.CodeProfileInvalid, "mapping.encrypted_data_in target is required")
	}
	if strings.TrimSpace(p.Mapping.EncryptedDataOut.Ref()) == "" {
		return coreerrors.New(coreerrors.CodeProfileInvalid, "mapping.encrypted_data_out target is required")
	}
	return nil
}

func ValidateMany(profiles []Profile) error {
	for i := range profiles {
		if err := Validate(profiles[i]); err != nil {
			return coreerrors.Wrap(coreerrors.CodeProfileInvalid, fmt.Sprintf("profile[%d] invalid", i), err)
		}
	}
	return nil
}

// ValidateSet validates constraints across all enabled profiles in one channel set.
func ValidateSet(profiles []Profile) error {
	hintSeen := map[string]string{}
	shapeSeen := map[string]string{}
	enabledCount := 0

	for _, p := range profiles {
		if !p.Enabled {
			continue
		}
		enabledCount++

		if p.Mapping.ProfileID != nil {
			hint := strings.ToLower(strings.TrimSpace(p.Mapping.ProfileID.Ref()))
			if hint != "" {
				if prev, ok := hintSeen[hint]; ok {
					return coreerrors.New(coreerrors.CodeProfileAmbiguous, "overlapping enabled mapping.profile_id between "+prev+" and "+p.ProfileID)
				}
				hintSeen[hint] = p.ProfileID
			}
		}

		shape := strings.ToLower(strings.TrimSpace(p.Mapping.ID.Ref())) + "|" + strings.ToLower(strings.TrimSpace(p.Mapping.EncryptedDataIn.Ref()))
		if prev, ok := shapeSeen[shape]; ok {
			return coreerrors.New(coreerrors.CodeProfileAmbiguous, "overlapping enabled mapping shape between "+prev+" and "+p.ProfileID)
		}
		shapeSeen[shape] = p.ProfileID
	}

	if enabledCount == 0 {
		return coreerrors.New(coreerrors.CodeProfileInvalid, "at least one enabled profile is required")
	}
	return nil
}
