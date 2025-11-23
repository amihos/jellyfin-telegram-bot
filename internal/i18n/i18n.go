// SPDX-License-Identifier: MIT

package i18n

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Supported languages
var SupportedLanguages = []string{"en", "fa"}

// Default language
const DefaultLanguage = "en"

// Context key for language
type contextKey string

const languageContextKey contextKey = "language"

// InitBundle initializes the i18n bundle with all supported languages
func InitBundle() (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Determine locales directory path
	// Try current directory first, then parent directory (for tests)
	localesDir := findLocalesDir()

	// Load English translations (default)
	enPath := filepath.Join(localesDir, "active.en.toml")
	if _, err := bundle.LoadMessageFile(enPath); err != nil {
		return nil, fmt.Errorf("failed to load English translations from %s: %w", enPath, err)
	}

	// Load Persian translations
	faPath := filepath.Join(localesDir, "active.fa.toml")
	if _, err := bundle.LoadMessageFile(faPath); err != nil {
		return nil, fmt.Errorf("failed to load Persian translations from %s: %w", faPath, err)
	}

	return bundle, nil
}

// findLocalesDir finds the locales directory
func findLocalesDir() string {
	// Try current directory
	if _, err := os.Stat("locales"); err == nil {
		return "locales"
	}

	// Try parent directory (for tests run from test/ directory)
	if _, err := os.Stat("../../locales"); err == nil {
		return "../../locales"
	}

	// Try absolute path from project root
	// This works when running from any subdirectory
	if _, err := os.Stat("/home/huso/jellyfin-telegram-bot/locales"); err == nil {
		return "/home/huso/jellyfin-telegram-bot/locales"
	}

	// Default to "locales" and let it fail if not found
	return "locales"
}

// GetLocalizer returns a localizer for the specified language code
func GetLocalizer(bundle *i18n.Bundle, langCode string) *i18n.Localizer {
	// Normalize language code (handle fa-IR -> fa)
	lang := normalizeLanguageCode(langCode)

	// Create localizer with fallback chain
	return i18n.NewLocalizer(bundle, lang, DefaultLanguage)
}

// T translates a message key without template data
func T(localizer *i18n.Localizer, key string) string {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: key,
	})
	if err != nil {
		// Return key if translation not found (for debugging)
		return key
	}
	return msg
}

// TWithData translates a message key with template data
func TWithData(localizer *i18n.Localizer, key string, data map[string]interface{}) string {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
	})
	if err != nil {
		// Return key if translation not found (for debugging)
		return key
	}
	return msg
}

// DetectLanguage detects the user's language from Telegram language code
// Returns a supported language code or the default language
func DetectLanguage(telegramLangCode string, supportedLanguages []string) string {
	if telegramLangCode == "" {
		return DefaultLanguage
	}

	// Normalize the language code (e.g., "fa-IR" -> "fa")
	normalized := normalizeLanguageCode(telegramLangCode)

	// Check if language is supported
	for _, supported := range supportedLanguages {
		if normalized == supported {
			return normalized
		}
	}

	// Fallback to default language
	return DefaultLanguage
}

// normalizeLanguageCode extracts the base language from a language code
// e.g., "fa-IR" -> "fa", "en-US" -> "en"
func normalizeLanguageCode(langCode string) string {
	if langCode == "" {
		return DefaultLanguage
	}

	// Split by hyphen and take the first part
	parts := strings.Split(langCode, "-")
	return strings.ToLower(parts[0])
}

// IsSupportedLanguage checks if a language code is supported
func IsSupportedLanguage(langCode string) bool {
	normalized := normalizeLanguageCode(langCode)
	for _, supported := range SupportedLanguages {
		if normalized == supported {
			return true
		}
	}
	return false
}

// WithLanguage adds language to context
func WithLanguage(ctx context.Context, langCode string) context.Context {
	return context.WithValue(ctx, languageContextKey, langCode)
}

// GetLanguageFromContext retrieves language from context
func GetLanguageFromContext(ctx context.Context) string {
	if lang, ok := ctx.Value(languageContextKey).(string); ok && lang != "" {
		return lang
	}
	return DefaultLanguage
}

// GetLocalizerFromContext creates a localizer from context language
func GetLocalizerFromContext(ctx context.Context, bundle *i18n.Bundle) *i18n.Localizer {
	lang := GetLanguageFromContext(ctx)
	return GetLocalizer(bundle, lang)
}
