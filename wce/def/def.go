package def

type Arg struct {
	Name        string `yaml:"name"`
	Note        string `yaml:"note"`
	Description string `yaml:"description"`
	Format      string `yaml:"format"`
	Example     string `yaml:"example"`
}

type Property struct {
	Name        string     `yaml:"name"`
	Note        string     `yaml:"note"`
	Description string     `yaml:"description"`
	Args        []Arg      `yaml:"args"`
	Properties  []Property `yaml:"properties"`
}

type Definition struct {
	Name        string     `yaml:"name"`
	HasTag      bool       `yaml:"hasTag,omitempty"`
	Note        string     `yaml:"note"`
	Description string     `yaml:"description"`
	Properties  []Property `yaml:"properties"`
}
