package webapp

type WebApp struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type WebAppConfig struct {
	Apps []WebApp `yaml:"apps"`
}
