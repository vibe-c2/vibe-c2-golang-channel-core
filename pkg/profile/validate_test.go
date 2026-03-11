package profile

import "testing"

func mf(loc, key string) MapField { return MapField{Target: Target{Location: loc, Key: key}} }

func TestValidateSetRequiresEnabled(t *testing.T) {
	profiles := []Profile{{ProfileID: "a", ChannelType: "http", Enabled: false, Mapping: Mapping{ID: mf("body", "id"), EncryptedDataIn: mf("body", "in"), EncryptedDataOut: mf("body", "out")}}}
	if err := ValidateSet(profiles); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateSetOverlapDetection(t *testing.T) {
	profiles := []Profile{
		{ProfileID: "a", ChannelType: "http", Enabled: true, Mapping: Mapping{ProfileID: &MapField{Target: Target{Location: "body", Key: "hint"}}, ID: mf("body", "id"), EncryptedDataIn: mf("body", "in"), EncryptedDataOut: mf("body", "out")}},
		{ProfileID: "b", ChannelType: "http", Enabled: true, Mapping: Mapping{ProfileID: &MapField{Target: Target{Location: "body", Key: "hint"}}, ID: mf("body", "id2"), EncryptedDataIn: mf("body", "in2"), EncryptedDataOut: mf("body", "out2")}},
	}
	if err := ValidateSet(profiles); err == nil {
		t.Fatal("expected overlap error")
	}
}
