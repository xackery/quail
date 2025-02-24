package def

type Arg struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Format      string `yaml:"format"`
	Example     string `yaml:"example"`
}

type Property struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	IsArrayNum  bool       `yaml:"isArrayNum,omitempty"`
	Args        []Arg      `yaml:"args"`
	Properties  []Property `yaml:"properties"`
}

type Definition struct {
	Name        string     `yaml:"name"`
	HasTag      bool       `yaml:"hasTag,omitempty"`
	Description string     `yaml:"description"`
	Properties  []Property `yaml:"properties"`
}
