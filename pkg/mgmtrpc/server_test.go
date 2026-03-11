package mgmtrpc

import (
	"context"
	"errors"
	"testing"

	"github.com/vibe-c2/vibe-c2-golang-channel-core/pkg/profile"
)

type memStore struct {
	data map[string]map[string]profile.Profile
}

func newMemStore() *memStore { return &memStore{data: map[string]map[string]profile.Profile{}} }

func (m *memStore) List(_ context.Context, channelID string) ([]profile.Profile, error) {
	ch := m.data[channelID]
	out := make([]profile.Profile, 0, len(ch))
	for _, p := range ch {
		out = append(out, p)
	}
	return out, nil
}
func (m *memStore) Get(_ context.Context, channelID, profileID string) (profile.Profile, error) {
	if ch, ok := m.data[channelID]; ok {
		if p, ok := ch[profileID]; ok {
			return p, nil
		}
	}
	return profile.Profile{}, errors.New("not found")
}
func (m *memStore) Put(_ context.Context, channelID string, p profile.Profile) error {
	if _, ok := m.data[channelID]; !ok {
		m.data[channelID] = map[string]profile.Profile{}
	}
	m.data[channelID][p.ProfileID] = p
	return nil
}
func (m *memStore) Delete(_ context.Context, channelID, profileID string) error {
	if _, ok := m.data[channelID]; ok {
		delete(m.data[channelID], profileID)
	}
	return nil
}

func TestCreateAndValidate(t *testing.T) {
	s := NewServer(newMemStore())
	ctx := context.Background()
	p := profile.Profile{ProfileID: "p1", ChannelType: "http", Enabled: true, Mapping: profile.Mapping{ID: profile.MapField{Target: profile.Target{Location: "body", Key: "id"}}, EncryptedDataIn: profile.MapField{Target: profile.Target{Location: "body", Key: "data_in"}}, EncryptedDataOut: profile.MapField{Target: profile.Target{Location: "body", Key: "data_out"}}}}
	if err := s.CreateProfile(ctx, "c1", p); err != nil {
		t.Fatal(err)
	}
}

func TestActivateNotFound(t *testing.T) {
	s := NewServer(newMemStore())
	if err := s.ActivateProfile(context.Background(), "c1", "x"); err == nil {
		t.Fatal("expected error")
	}
}
