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
			{Type: newString("user"), Label: newString("user")},
			{Type: newString("hostname"), Label: newString("hostname")},
			{Type: newString("os"), Label: newString("os")},
			{Type: newString("kernel"), Label: newString("kernel")},
			{Type: newString("uptime"), Label: newString("uptime")},
			{Type: newString("shell"), Label: newString("shell")},
			{Type: newString("packages"), Label: newString("packages")},
			{Type: newString("memory"), Label: newString("memory")},
			{Type: newString("colors"), Label: newString("colors")},
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