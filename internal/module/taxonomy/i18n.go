package taxonomy

import "strings"

func normalizeLanguage(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	switch lang {
	case "en", "zh":
		return lang
	default:
		return "zh"
	}
}
