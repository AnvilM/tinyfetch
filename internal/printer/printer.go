package printer

import (
	"fmt"
	"io/ioutil"
	"strings"
	"tinyfetch/internal/config"
	"tinyfetch/internal/types"
	"tinyfetch/internal/utils/logger"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

func Ptr[T any](v T) *T { return &v }

func Print(cfg config.Config) {
	printTitle(cfg.Title)
	if cfg.Modules != nil && cfg.Container != nil {
		printInfo(*cfg.Modules, *cfg.Container)
	}
}

func getColor(col config.Color) *color.Color {
	switch col {
	case "white":
		return color.New(color.FgWhite)
	case "black":
		return color.New(color.FgBlack)
	case "red":
		return color.New(color.FgRed)
	case "green":
		return color.New(color.FgGreen)
	case "blue":
		return color.New(color.FgBlue)
	case "yellow":
		return color.New(color.FgYellow)
	case "magenta":
		return color.New(color.FgMagenta)
	case "cyan":
		return color.New(color.FgCyan)
	default:
		return color.New(color.Reset)
	}
}

func printTitle(title *config.Title) {
	if title == nil || title.FilePath == nil || title.Color == nil {
		return
	}

	filePath, err := homedir.Expand(*title.FilePath)
	if err != nil {
		logger.Fatal("Error expanding path: %v\n", err)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Fatal("Error reading file: %v\n", err)
	}

	c := getColor(*title.Color)
	c.Print(string(data))
}

func printInfo(modules []config.Module, container config.Container) {

	marginLeft := 0
	if container.MarginLeft != nil {
		marginLeft = *container.MarginLeft
	}
	marginRight := 0
	if container.MarginRight != nil {
		marginRight = *container.MarginRight
	}
	paddingLeft := 0
	if container.PaddingLeft != nil {
		paddingLeft = *container.PaddingLeft
	}
	paddingRight := 0
	if container.PaddingRight != nil {
		paddingRight = *container.PaddingRight
	}

	maxLen := 0
	for _, module := range modules {
		labelLen := 0
		if module.Label != nil {
			labelLen = len(*module.Label)
		}

		length := labelLen
		if module.Prefix != nil {
			length += utf8.RuneCountInString(*module.Prefix)
		}

		if length > maxLen {
			maxLen = length
		}
	}

	boxWidth := maxLen + paddingLeft + paddingRight

	boxBorderTop := strings.Repeat(" ", marginLeft) + getColor(*container.BorderColor).Sprint("╭") + getColor(*container.BorderColor).Sprint(strings.Repeat("─", boxWidth)) + getColor(*container.BorderColor).Sprint("╮")
	boxBorderBottom := strings.Repeat(" ", marginLeft) + getColor(*container.BorderColor).Sprint("╰") + getColor(*container.BorderColor).Sprint(strings.Repeat("─", boxWidth)) + getColor(*container.BorderColor).Sprint("╯")

	fmt.Println(boxBorderTop)

	// Print modules
	for i, module := range modules {

		if i == len(modules)-1 && len(modules) > 1 {
			fmt.Println(strings.Repeat(" ", marginLeft) + getColor(*container.BorderColor).Sprint("├") + getColor(*container.BorderColor).Sprint(strings.Repeat("─", boxWidth)) + getColor(*container.BorderColor).Sprint("┤"))
		}

		fmt.Print(strings.Repeat(" ", marginLeft) + getColor(*container.BorderColor).Sprint("│") + strings.Repeat(" ", paddingLeft))

		// Prefix
		if module.Prefix != nil && module.PrefixColor != nil {
			getColor(*module.PrefixColor).Print(*module.Prefix)
		}

		// Label
		if module.Label != nil && module.LabelColor != nil {
			getColor(*module.LabelColor).Print(*module.Label)
		} else if module.Label != nil {
			fmt.Print(*module.Label)
		}

		labelLen := 0
		if module.Label != nil {
			labelLen = len(*module.Label)
		}
		prefixMargin := 0
		if module.Prefix != nil {
			prefixMargin =  utf8.RuneCountInString(*module.Prefix)
		}

		spaces := maxLen - labelLen - prefixMargin + paddingRight
		fmt.Print(strings.Repeat(" ", spaces))

		fmt.Print(getColor(*container.BorderColor).Sprint("│") + strings.Repeat(" ", marginRight))

		// Info
		if module.Type != nil && module.InfoColor != nil {
			getColor(*module.InfoColor).Print(types.GetTypeInfo(*module.Type) + "\n")
		} else if module.Type != nil {
			fmt.Print(types.GetTypeInfo(*module.Type) + "\n")
		} else {
			fmt.Print("\n")
		}
	}

	fmt.Println(boxBorderBottom)
}
