package def

type Arg struct {
	Name        string `yaml:"name"`
	Comment     string `yaml:"comment"`
	Description string `yaml:"description"`
	Format      string `yaml:"format"`
	Example     string `yaml:"example"`
}

type Property struct {
	Name        string     `yaml:"name"`
	Comment     string     `yaml:"comment"`
	Description string     `yaml:"description"`
	Args        []Arg      `yaml:"args"`
	Properties  []Property `yaml:"properties"`
}

type Definition struct {
	Name        string     `yaml:"name"`
	HasTag      bool       `yaml:"hasTag,omitempty"`
	Comment     string     `yaml:"comment"`
	Description string     `yaml:"description"`
	Properties  []Property `yaml:"properties"`
}
