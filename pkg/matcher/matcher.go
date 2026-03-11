package matcher

import (
	"context"
	"sort"
	"strings"

	coreerrors "github.com/vibe-c2/vibe-c2-golang-channel-core/pkg/errors"
	"github.com/vibe-c2/vibe-c2-golang-channel-core/pkg/profile"
)

type Matcher struct{}

func New() *Matcher { return &Matcher{} }

type MatchSource string

const (
	MatchSourceHint       MatchSource = "hint"
	MatchSourceBruteforce MatchSource = "bruteforce"
	MatchSourceNone       MatchSource = "none"
)

type Resolution struct {
	Profile profile.Profile
	Source  MatchSource
}

// Resolve resolves only by explicit hint. Brute-force candidate iteration is done by channel module.
func (m *Matcher) Resolve(_ context.Context, hintProfileID string, candidates []profile.Profile) (Resolution, error) {
	hintProfileID = strings.TrimSpace(hintProfileID)
	if hintProfileID == "" {
		return Resolution{Source: MatchSourceNone}, coreerrors.New(coreerrors.CodeProfileNotFound, "no hint provided")
	}

	hits := make([]profile.Profile, 0, 1)
	for _, p := range candidates {
		if !p.Enabled {
			continue
		}
		if p.ProfileID == hintProfileID {
			hits = append(hits, p)
		}
	}
	switch len(hits) {
	case 1:
		return Resolution{Profile: hits[0], Source: MatchSourceHint}, nil
	case 0:
		return Resolution{Source: MatchSourceNone}, coreerrors.New(coreerrors.CodeProfileNotFound, "hint did not match enabled profiles")
	default:
		return Resolution{Source: MatchSourceNone}, coreerrors.New(coreerrors.CodeProfileAmbiguous, "hint matches multiple enabled profiles")
	}
}

// EnabledOrdered returns enabled profiles ordered by priority desc, profile_id asc.
func (m *Matcher) EnabledOrdered(candidates []profile.Profile) []profile.Profile {
	out := make([]profile.Profile, 0, len(candidates))
	for _, p := range candidates {
		if p.Enabled {
			out = append(out, p)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Priority != out[j].Priority {
			return out[i].Priority > out[j].Priority
		}
		return out[i].ProfileID < out[j].ProfileID
	})
	return out
}
