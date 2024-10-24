package gen

type Config struct {
	Imports []Import `koanf:"imports"`
}

type Import struct {
	Alias string `koanf:"alias"`
	Path  string `koanf:"path"`
}
