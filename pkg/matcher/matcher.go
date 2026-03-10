package matcher

import (
	"context"

	"github.com/vibe-c2/vibe-c2-golang-channel-core/pkg/profile"
)

// Matcher contains profile selection helpers.
type Matcher struct{}

func New() *Matcher {
	return &Matcher{}
}

// SelectHintFirst performs hint-first selection and falls back to the
// highest-priority enabled default profile.
func (m *Matcher) SelectHintFirst(_ context.Context, _ string, hintProfileID string, candidates []profile.Profile) (profile.Profile, bool, error) {
	for _, p := range candidates {
		if !p.Enabled {
			continue
		}
		if hintProfileID != "" && (p.ProfileID == hintProfileID || p.Mapping.ProfileID == hintProfileID) {
			return p, true, nil
		}
	}

	var selected profile.Profile
	found := false
	for _, p := range candidates {
		if !p.Enabled || !p.DefaultFallback {
			continue
		}
		if !found || p.Priority > selected.Priority {
			selected = p
			found = true
		}
	}

	return selected, found, nil
}
