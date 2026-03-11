package profile

import "testing"

func TestValidateSetRequiresSingleFallback(t *testing.T) {
	profiles := []Profile{
		{ProfileID: "a", ChannelType: "http", Enabled: true, DefaultFallback: true, Mapping: Mapping{ID: "id", EncryptedData: "data"}},
		{ProfileID: "b", ChannelType: "http", Enabled: true, DefaultFallback: true, Mapping: Mapping{ID: "id", EncryptedData: "data"}},
	}
	if err := ValidateSet(profiles); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateSetOverlapDetection(t *testing.T) {
	profiles := []Profile{
		{ProfileID: "a", ChannelType: "http", Enabled: true, DefaultFallback: true, Mapping: Mapping{ProfileID: "hint", ID: "id", EncryptedData: "data"}},
		{ProfileID: "b", ChannelType: "http", Enabled: true, DefaultFallback: false, Mapping: Mapping{ProfileID: "hint", ID: "id", EncryptedData: "data"}},
	}
	if err := ValidateSet(profiles); err == nil {
		t.Fatal("expected overlap error")
	}
}
