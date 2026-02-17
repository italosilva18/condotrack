package setting

import (
	"context"
	"fmt"
	"regexp"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
)

// UseCase handles setting business logic
type UseCase struct {
	settingRepo repository.SettingRepository
}

// NewUseCase creates a new setting use case
func NewUseCase(settingRepo repository.SettingRepository) *UseCase {
	return &UseCase{
		settingRepo: settingRepo,
	}
}

// GetAllSettings returns all settings (with secret values masked)
func (uc *UseCase) GetAllSettings(ctx context.Context) ([]*entity.SettingPublic, error) {
	settings, err := uc.settingRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all settings: %w", err)
	}

	publicSettings := make([]*entity.SettingPublic, len(settings))
	for i, s := range settings {
		publicSettings[i] = s.ToPublic()
	}

	return publicSettings, nil
}

// GetSettingsByCategory returns settings grouped by category
func (uc *UseCase) GetSettingsByCategory(ctx context.Context) ([]*entity.SettingsByCategory, error) {
	settings, err := uc.settingRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all settings: %w", err)
	}

	// Group settings by category
	categoryMap := make(map[string][]*entity.SettingPublic)
	categoryOrder := []string{}

	for _, s := range settings {
		category := string(s.Category)
		if _, exists := categoryMap[category]; !exists {
			categoryOrder = append(categoryOrder, category)
			categoryMap[category] = []*entity.SettingPublic{}
		}
		categoryMap[category] = append(categoryMap[category], s.ToPublic())
	}

	// Build result maintaining order
	result := make([]*entity.SettingsByCategory, 0, len(categoryOrder))
	for _, cat := range categoryOrder {
		label, ok := entity.CategoryLabels[entity.SettingCategory(cat)]
		if !ok {
			label = cat
		}
		result = append(result, &entity.SettingsByCategory{
			Category: cat,
			Label:    label,
			Settings: categoryMap[cat],
		})
	}

	return result, nil
}

// GetSettingByKey returns a setting by its key (with secret value masked)
func (uc *UseCase) GetSettingByKey(ctx context.Context, key string) (*entity.SettingPublic, error) {
	setting, err := uc.settingRepo.GetByKey(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get setting: %w", err)
	}
	if setting == nil {
		return nil, fmt.Errorf("setting not found: %s", key)
	}

	return setting.ToPublic(), nil
}

// GetSettingValue returns the raw value of a setting (for internal use)
func (uc *UseCase) GetSettingValue(ctx context.Context, key string) (string, error) {
	return uc.settingRepo.GetValue(ctx, key)
}

// UpdateSetting updates a single setting value
func (uc *UseCase) UpdateSetting(ctx context.Context, key string, value string) error {
	// Verify setting exists
	setting, err := uc.settingRepo.GetByKey(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to get setting: %w", err)
	}
	if setting == nil {
		return fmt.Errorf("setting not found: %s", key)
	}

	// Validate required fields
	if setting.IsRequired && value == "" {
		return fmt.Errorf("setting %s is required", key)
	}

	// Validate value against regex pattern if defined
	if setting.ValidationRegex != nil && *setting.ValidationRegex != "" && value != "" {
		matched, err := regexp.MatchString(*setting.ValidationRegex, value)
		if err != nil {
			return fmt.Errorf("invalid validation regex for setting %s: %w", key, err)
		}
		if !matched {
			return fmt.Errorf("value for setting %s does not match the required format", key)
		}
	}

	return uc.settingRepo.Update(ctx, key, value)
}

// UpdateSettingByID updates a setting by its ID
func (uc *UseCase) UpdateSettingByID(ctx context.Context, id string, value string) error {
	// Verify setting exists
	setting, err := uc.settingRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get setting: %w", err)
	}
	if setting == nil {
		return fmt.Errorf("setting not found: %s", id)
	}

	// Validate required fields
	if setting.IsRequired && value == "" {
		return fmt.Errorf("setting %s is required", setting.Key)
	}

	return uc.settingRepo.UpdateByID(ctx, id, value)
}

// BulkUpdateSettings updates multiple settings at once
func (uc *UseCase) BulkUpdateSettings(ctx context.Context, settings map[string]string) error {
	// Validate all settings exist and required values are provided
	for key, value := range settings {
		setting, err := uc.settingRepo.GetByKey(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to validate setting %s: %w", key, err)
		}
		if setting == nil {
			return fmt.Errorf("setting not found: %s", key)
		}
		if setting.IsRequired && value == "" {
			return fmt.Errorf("setting %s is required", key)
		}
	}

	return uc.settingRepo.BulkUpdate(ctx, settings)
}

// GetCategories returns all available categories
func (uc *UseCase) GetCategories(ctx context.Context) ([]string, error) {
	return uc.settingRepo.GetCategories(ctx)
}

// GetAsaasConfig returns Asaas configuration for internal use
func (uc *UseCase) GetAsaasConfig(ctx context.Context) (apiKey, apiURL, webhookToken, env string, err error) {
	apiKey, _ = uc.settingRepo.GetValue(ctx, "asaas_api_key")
	apiURL, _ = uc.settingRepo.GetValue(ctx, "asaas_api_url")
	webhookToken, _ = uc.settingRepo.GetValue(ctx, "asaas_webhook_token")
	env, _ = uc.settingRepo.GetValue(ctx, "asaas_env")
	return
}

// GetGeminiAPIKey returns the Gemini API key for internal use
func (uc *UseCase) GetGeminiAPIKey(ctx context.Context) (string, error) {
	return uc.settingRepo.GetValue(ctx, "gemini_api_key")
}

// IsAIEnabled returns whether AI features are enabled
func (uc *UseCase) IsAIEnabled(ctx context.Context) (bool, error) {
	value, err := uc.settingRepo.GetValue(ctx, "ai_enabled")
	if err != nil {
		return false, err
	}
	return value == "true" || value == "1", nil
}
