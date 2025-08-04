package stringcase

type Caser struct {
	display string
}

func NewCaser(s string) Caser {
	return Caser{display: s}
}

func (c Caser) Display() string    { return c.display }
func (c Caser) CamelCase() string  { return ToCamelCase(c.display) }
func (c Caser) KebabCase() string  { return ToKebabCase(c.display) }
func (c Caser) PascalCase() string { return ToPascalCase(c.display) }
func (c Caser) SnakeCase() string  { return ToSnakeCase(c.display) }

type CaserNoun struct {
	Singular Caser
	Plural   Caser
}

func NewCaserNoun(singular, plural string) CaserNoun {
	return CaserNoun{
		Singular: NewCaser(singular),
		Plural:   NewCaser(plural)}
}
