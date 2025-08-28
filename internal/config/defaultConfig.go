package config

func GetDefaultConfig() Config {
	return Config{
		Title: &Title{
			Color: newColor("white"),
		},
		Container: &Container{
			MarginLeft: newInt(2),
			MarginRight: newInt(1),
			PaddingLeft: newInt(1),
			PaddingRight: newInt(1),
			BorderColor: newColor("white"),
		},
		Modules: &[]Module{
			{Type: newString("user"), InfoColor: newColor("white"), Label: newString("user"), LabelColor: newColor("white")},
		},
	}
}

func newColor(s string) *Color {
	c := Color(s)
	return &c
}

func newInt(i int) *int {
	return &i
}

func newString(s string) *string {
	return &s
}