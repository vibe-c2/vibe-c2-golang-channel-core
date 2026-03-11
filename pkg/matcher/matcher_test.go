package matcher

import (
	"context"
	"testing"

	"github.com/vibe-c2/vibe-c2-golang-channel-core/pkg/profile"
)

func p(id string, enabled bool, hint string, prio int) profile.Profile {
	var hintField *profile.MapField
	if hint != "" {
		hintField = &profile.MapField{Target: profile.Target{Location: "body", Key: hint}}
	}
	return profile.Profile{
		ProfileID:   id,
		ChannelType: "http",
		Enabled:     enabled,
		Priority:    prio,
		Mapping: profile.Mapping{
			ProfileID:        hintField,
			ID:               profile.MapField{Target: profile.Target{Location: "body", Key: "id"}},
			EncryptedDataIn:  profile.MapField{Target: profile.Target{Location: "body", Key: "in"}},
			EncryptedDataOut: profile.MapField{Target: profile.Target{Location: "body", Key: "out"}},
		},
	}
}

func TestResolveHintSingle(t *testing.T) {
	m := New()
	res, err := m.Resolve(context.Background(), "b", []profile.Profile{p("a", true, "x", 1), p("b", true, "p1", 1)})
	if err != nil {
		t.Fatal(err)
	}
	if res.Profile.ProfileID != "b" {
		t.Fatalf("unexpected profile: %s", res.Profile.ProfileID)
	}
}

func TestResolveNoHint(t *testing.T) {
	m := New()
	if _, err := m.Resolve(context.Background(), "", []profile.Profile{p("a", true, "x", 1)}); err == nil {
		t.Fatal("expected error")
	}
}

func TestEnabledOrdered(t *testing.T) {
	m := New()
	out := m.EnabledOrdered([]profile.Profile{p("a", true, "x", 1), p("b", true, "y", 10), p("c", false, "z", 100)})
	if len(out) != 2 || out[0].ProfileID != "b" {
		t.Fatalf("unexpected order: %+v", out)
	}
}
