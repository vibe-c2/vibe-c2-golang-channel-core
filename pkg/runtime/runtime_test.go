package runtime

import (
	"context"
	"testing"

	coreerrors "github.com/vibe-c2/vibe-c2-golang-channel-core/pkg/errors"
	protocol "github.com/vibe-c2/vibe-c2-golang-protocol/protocol"
)

type testEnvelope struct {
	source string
	data   map[string]string
}

func (e *testEnvelope) SourceKey() string {
	return e.source
}

func (e *testEnvelope) GetField(location, key string) (string, bool) {
	value, ok := e.data[location+"."+key]
	return value, ok
}

func (e *testEnvelope) SetField(location, key, value string) {
	if e.data == nil {
		e.data = map[string]string{}
	}
	e.data[location+"."+key] = value
}

type testSyncClient struct {
	outbound protocol.OutboundAgentMessage
	err      error
	captured protocol.InboundAgentMessage
}

func (s *testSyncClient) Sync(_ context.Context, in protocol.InboundAgentMessage) (protocol.OutboundAgentMessage, error) {
	s.captured = in
	return s.outbound, s.err
}

func TestRuntimeHandleSuccess(t *testing.T) {
	env := &testEnvelope{data: map[string]string{
		"mapping.profile_id":     "alpha",
		"mapping.id":             "msg-1",
		"mapping.encrypted_data": "blob-in",
	}}
	sync := &testSyncClient{outbound: protocol.OutboundAgentMessage{
		MessageID:     "m-2",
		Type:          protocol.TypeOutboundAgentMessage,
		Version:       protocol.VersionV1,
		Timestamp:     "2026-03-10T15:00:00Z",
		Source:        protocol.SourceInfo{Module: "core", ModuleInstance: "main", Transport: "channel", Tenant: "default"},
		ID:            "msg-2",
		EncryptedData: "blob-out",
	}}
	r := New(sync)

	got, err := r.Handle(context.Background(), env, "telegram")
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	if got.ID != "msg-2" {
		t.Fatalf("unexpected outbound id: %s", got.ID)
	}
	if sync.captured.Type != protocol.TypeInboundAgentMessage {
		t.Fatalf("unexpected inbound type: %s", sync.captured.Type)
	}
	if env.data["mapping.id"] != "msg-2" {
		t.Fatalf("expected envelope mapping.id to be updated")
	}
}

func TestRuntimeHandleMissingRequiredField(t *testing.T) {
	env := &testEnvelope{data: map[string]string{
		"mapping.encrypted_data": "blob-in",
	}}
	r := New(&testSyncClient{})

	_, err := r.Handle(context.Background(), env, "http")
	if err == nil {
		t.Fatal("expected error for missing mapping.id")
	}
	if code := coreerrors.Code(err); code != coreerrors.CodeInvalidInput {
		t.Fatalf("unexpected error code: %s", code)
	}
}

func TestRuntimeHandleInvalidOutbound(t *testing.T) {
	env := &testEnvelope{data: map[string]string{
		"mapping.id":             "msg-1",
		"mapping.encrypted_data": "blob-in",
	}}
	sync := &testSyncClient{outbound: protocol.OutboundAgentMessage{
		MessageID: "m-3",
		Type:      protocol.TypeOutboundAgentMessage,
		Version:   protocol.VersionV1,
		Timestamp: "2026-03-10T15:00:00Z",
		Source:    protocol.SourceInfo{Module: "core", ModuleInstance: "main", Transport: "channel", Tenant: "default"},
		ID:        "msg-2",
		// missing encrypted_data
	}}
	r := New(sync)

	_, err := r.Handle(context.Background(), env, "http")
	if err == nil {
		t.Fatal("expected outbound validation error")
	}
	if code := coreerrors.Code(err); code != coreerrors.CodeCanonicalInvalid {
		t.Fatalf("unexpected error code: %s", code)
	}
}
