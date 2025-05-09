package gen

type Config struct {
	Imports  []Import  `koanf:"imports"`
	Preloads []Preload `koanf:"preloads"`
}

type Preload struct {
	Alias string `koanf:"alias"`
	Path  string `koanf:"path"`
}
type Import struct {
	Alias string `koanf:"alias"`
	Path  string `koanf:"path"`
}
