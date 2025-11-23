package integration

import (
	"context"
	"testing"

	"jellyfin-telegram-bot/internal/i18n"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestI18nBundleInitialization tests that the i18n bundle initializes correctly
func TestI18nBundleInitialization(t *testing.T) {
	bundle, err := i18n.InitBundle()
	require.NoError(t, err, "Bundle initialization should not fail")
	assert.NotNil(t, bundle, "Bundle should not be nil")
}

// TestI18nGetLocalizer tests creating localizers for different languages
func TestI18nGetLocalizer(t *testing.T) {
	bundle, err := i18n.InitBundle()
	require.NoError(t, err)

	tests := []struct {
		name     string
		langCode string
		wantErr  bool
	}{
		{
			name:     "English localizer",
			langCode: "en",
			wantErr:  false,
		},
		{
			name:     "Persian localizer",
			langCode: "fa",
			wantErr:  false,
		},
		{
			name:     "Default fallback for unsupported language",
			langCode: "de",
			wantErr:  false, // Should fallback to English
		},
		{
			name:     "Empty language code fallback",
			langCode: "",
			wantErr:  false, // Should fallback to English
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localizer := i18n.GetLocalizer(bundle, tt.langCode)
			assert.NotNil(t, localizer, "Localizer should not be nil")
		})
	}
}

// TestI18nTranslationKeys tests that critical translation keys exist
func TestI18nTranslationKeys(t *testing.T) {
	bundle, err := i18n.InitBundle()
	require.NoError(t, err)

	// Test both English and Persian
	languages := []string{"en", "fa"}

	// Critical message keys that must exist
	criticalKeys := []string{
		"welcome.message",
		"help.message",
		"error.generic",
		"notification.movie.header",
		"notification.episode.header",
		"command.start.description",
		"command.recent.description",
		"command.search.description",
		"command.mutedlist.description",
		"command.language.description",
	}

	for _, lang := range languages {
		localizer := i18n.GetLocalizer(bundle, lang)

		for _, key := range criticalKeys {
			t.Run(lang+"_"+key, func(t *testing.T) {
				translation := i18n.T(localizer, key)
				assert.NotEmpty(t, translation, "Translation for %s in %s should not be empty", key, lang)
				assert.NotEqual(t, key, translation, "Translation should not return the key itself")
			})
		}
	}
}

// TestI18nFallbackChain tests language fallback behavior
func TestI18nFallbackChain(t *testing.T) {
	bundle, err := i18n.InitBundle()
	require.NoError(t, err)

	// Get localizer for unsupported language (should fallback to English)
	localizer := i18n.GetLocalizer(bundle, "de")
	translation := i18n.T(localizer, "welcome.message")

	assert.NotEmpty(t, translation, "Should fallback to English for unsupported language")
}

// TestI18nWithTemplateData tests translation with template data
func TestI18nWithTemplateData(t *testing.T) {
	bundle, err := i18n.InitBundle()
	require.NoError(t, err)

	tests := []struct {
		name     string
		langCode string
		key      string
		data     map[string]interface{}
	}{
		{
			name:     "English with name",
			langCode: "en",
			key:      "search.no_results",
			data:     map[string]interface{}{"Query": "Interstellar"},
		},
		{
			name:     "Persian with name",
			langCode: "fa",
			key:      "search.no_results",
			data:     map[string]interface{}{"Query": "اینترستلار"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localizer := i18n.GetLocalizer(bundle, tt.langCode)
			translation := i18n.TWithData(localizer, tt.key, tt.data)
			assert.NotEmpty(t, translation, "Translation with data should not be empty")
			// Check that the data was interpolated (contains the query value)
			if query, ok := tt.data["Query"].(string); ok {
				assert.Contains(t, translation, query, "Translation should contain interpolated data")
			}
		})
	}
}

// TestI18nLanguageDetection tests language detection and storage
func TestI18nLanguageDetection(t *testing.T) {
	// This test will verify the language detection logic
	// It simulates what happens when a user interacts with the bot

	tests := []struct {
		name               string
		telegramLangCode   string
		expectedLang       string
		supportedLanguages []string
	}{
		{
			name:               "English user",
			telegramLangCode:   "en",
			expectedLang:       "en",
			supportedLanguages: []string{"en", "fa"},
		},
		{
			name:               "Persian user",
			telegramLangCode:   "fa",
			expectedLang:       "fa",
			supportedLanguages: []string{"en", "fa"},
		},
		{
			name:               "Persian with country code",
			telegramLangCode:   "fa-IR",
			expectedLang:       "fa",
			supportedLanguages: []string{"en", "fa"},
		},
		{
			name:               "Unsupported language fallback",
			telegramLangCode:   "de",
			expectedLang:       "en",
			supportedLanguages: []string{"en", "fa"},
		},
		{
			name:               "Empty language code fallback",
			telegramLangCode:   "",
			expectedLang:       "en",
			supportedLanguages: []string{"en", "fa"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detected := i18n.DetectLanguage(tt.telegramLangCode, tt.supportedLanguages)
			assert.Equal(t, tt.expectedLang, detected, "Language detection should return expected language")
		})
	}
}

// TestI18nUserPreferencePersistence tests language preference storage
// This is an integration test that requires database setup
func TestI18nUserPreferencePersistence(t *testing.T) {
	// Skip in short mode as this requires database
	if testing.Short() {
		t.Skip("Skipping database integration test in short mode")
	}

	// This test will be implemented after database changes are made
	// It will test:
	// 1. Setting user language preference
	// 2. Retrieving user language preference
	// 3. Fallback behavior when no preference is set
	t.Skip("Database integration test - to be implemented with database changes")
}

// TestI18nContextualLocalization tests getting localizer from context
func TestI18nContextualLocalization(t *testing.T) {
	bundle, err := i18n.InitBundle()
	require.NoError(t, err)

	ctx := context.Background()

	// Test setting and getting language from context
	ctx = i18n.WithLanguage(ctx, "fa")
	lang := i18n.GetLanguageFromContext(ctx)
	assert.Equal(t, "fa", lang, "Should retrieve Persian from context")

	// Test with English
	ctx = i18n.WithLanguage(ctx, "en")
	lang = i18n.GetLanguageFromContext(ctx)
	assert.Equal(t, "en", lang, "Should retrieve English from context")

	// Test with empty context (should fallback to default)
	ctx = context.Background()
	lang = i18n.GetLanguageFromContext(ctx)
	assert.Equal(t, "en", lang, "Should fallback to English for empty context")

	// Test getting localizer from context
	ctx = i18n.WithLanguage(ctx, "fa")
	localizer := i18n.GetLocalizerFromContext(ctx, bundle)
	assert.NotNil(t, localizer, "Should get localizer from context")
}
