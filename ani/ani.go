package ani

type ANI struct {
	name string
}

func New(name string) (*ANI, error) {
	e := &ANI{
		name: name,
	}
	return e, nil
}
