package locale

type Bundle interface {
	NewLocalizer(lang string) Localizer
}

type Localizer interface {
	Localize(tag string, args any, pluralCount int) string
}
