package entity

import "testing"

func TestSetting_ToPublic_NormalValue(t *testing.T) {
	val := "my-api-key"
	desc := "API key for service"
	s := &Setting{
		ID:       "s1",
		Key:      "api_key",
		Value:    &val,
		Type:     SettingTypeString,
		Category: CategoryGeneral,
		Label:    "API Key",
		Description: &desc,
		IsSecret: false,
		IsRequired: true,
		DisplayOrder: 1,
	}

	pub := s.ToPublic()

	if pub.Value != "my-api-key" {
		t.Errorf("ToPublic().Value = %q, want %q", pub.Value, "my-api-key")
	}
	if !pub.HasValue {
		t.Error("ToPublic().HasValue should be true")
	}
	if pub.Description != "API key for service" {
		t.Errorf("ToPublic().Description = %q, want %q", pub.Description, "API key for service")
	}
	if pub.ID != "s1" {
		t.Errorf("ToPublic().ID = %q, want %q", pub.ID, "s1")
	}
	if pub.Key != "api_key" {
		t.Errorf("ToPublic().Key = %q, want %q", pub.Key, "api_key")
	}
}

func TestSetting_ToPublic_SecretValue(t *testing.T) {
	val := "super-secret-key-123"
	s := &Setting{
		ID:       "s2",
		Key:      "secret_key",
		Value:    &val,
		Type:     SettingTypeSecret,
		Category: CategoryPayment,
		Label:    "Secret",
		IsSecret: true,
	}

	pub := s.ToPublic()

	if pub.Value != "********" {
		t.Errorf("ToPublic().Value for secret = %q, want %q", pub.Value, "********")
	}
	if !pub.HasValue {
		t.Error("ToPublic().HasValue for secret should be true")
	}
}

func TestSetting_ToPublic_NilValue(t *testing.T) {
	s := &Setting{
		ID:       "s3",
		Key:      "empty_key",
		Value:    nil,
		Type:     SettingTypeString,
		Category: CategoryGeneral,
		Label:    "Empty",
	}

	pub := s.ToPublic()

	if pub.Value != "" {
		t.Errorf("ToPublic().Value for nil = %q, want empty string", pub.Value)
	}
	if pub.HasValue {
		t.Error("ToPublic().HasValue for nil should be false")
	}
}

func TestSetting_ToPublic_EmptyValue(t *testing.T) {
	val := ""
	s := &Setting{
		ID:       "s4",
		Key:      "blank_key",
		Value:    &val,
		Type:     SettingTypeString,
		Category: CategoryGeneral,
		Label:    "Blank",
	}

	pub := s.ToPublic()

	if pub.HasValue {
		t.Error("ToPublic().HasValue for empty string should be false")
	}
}

func TestSetting_ToPublic_NilDescription(t *testing.T) {
	val := "test"
	s := &Setting{
		ID:          "s5",
		Key:         "no_desc",
		Value:       &val,
		Type:        SettingTypeString,
		Category:    CategoryGeneral,
		Label:       "No Desc",
		Description: nil,
	}

	pub := s.ToPublic()

	if pub.Description != "" {
		t.Errorf("ToPublic().Description for nil = %q, want empty string", pub.Description)
	}
}

func TestSetting_ToPublic_SecretNilValue(t *testing.T) {
	s := &Setting{
		ID:       "s6",
		Key:      "unset_secret",
		Value:    nil,
		Type:     SettingTypeSecret,
		Category: CategoryPayment,
		Label:    "Unset Secret",
		IsSecret: true,
	}

	pub := s.ToPublic()

	if pub.Value != "" {
		t.Errorf("ToPublic().Value for nil secret = %q, want empty string", pub.Value)
	}
	if pub.HasValue {
		t.Error("ToPublic().HasValue for nil secret should be false")
	}
}
