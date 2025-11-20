package types

// FrontendConfig представляет конфигурацию визуала для фронтенда
type FrontendConfig struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Theme      string                 `json:"theme"`      // light, dark, custom
	Colors     map[string]string      `json:"colors"`     // цветовая схема
	Layout     map[string]interface{} `json:"layout"`      // настройки layout
	Components map[string]interface{} `json:"components"` // настройки компонентов
	Branding   map[string]string      `json:"branding"`   // логотип, название и т.д.
	Active     bool                   `json:"active"`
	CreatedAt  string                 `json:"created_at,omitempty"`
	UpdatedAt  string                 `json:"updated_at,omitempty"`
}

