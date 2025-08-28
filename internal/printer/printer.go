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

// Ptr helper для указателей
func Ptr[T any](v T) *T { return &v }

// Print выводит заголовок и список модулей в рамке с цветами
func Print(cfg config.Config) {
	printTitle(cfg.Title)
	if cfg.Modules != nil && cfg.Container != nil {
		printInfo(*cfg.Modules, *cfg.Container)
	}
}

// getColor возвращает объект *color.Color для указанного цвета
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
	// Безопасные значения контейнера
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

	// Вычисляем maxLen
	maxLen := 0
	for _, module := range modules {
		labelLen := 0
		if module.Label != nil {
			labelLen = len(*module.Label)
		}

		iconMargin := 0
		if module.IconMarginRight != nil {
			iconMargin = *module.IconMarginRight
		}

		length := labelLen
		if module.Icon != nil {
			length += utf8.RuneCountInString(*module.Icon) + iconMargin
		}

		if length > maxLen {
			maxLen = length
		}
	}

	boxWidth := maxLen + paddingLeft + paddingRight

	boxBorderTop := strings.Repeat(" ", marginLeft) + "╭" + strings.Repeat("─", boxWidth) + "╮"
	boxBorderBottom := strings.Repeat(" ", marginLeft) + "╰" + strings.Repeat("─", boxWidth) + "╯"

	fmt.Println(boxBorderTop)

	// Print modules
	for i, module := range modules {
		// Разделитель между модулями
		if i > 0 {
			fmt.Println(strings.Repeat(" ", marginLeft) + "├" + strings.Repeat("─", boxWidth) + "┤")
		}

		fmt.Print(strings.Repeat(" ", marginLeft) + "│" + strings.Repeat(" ", paddingLeft))

		// Icon
		if module.Icon != nil && module.IconColor != nil {
			getColor(*module.IconColor).Print(*module.Icon)
			iconMargin := 0
			if module.IconMarginRight != nil {
				iconMargin = *module.IconMarginRight
			}
			fmt.Print(strings.Repeat(" ", iconMargin))
		}

		// Label
		if module.Label != nil && module.LabelColor != nil {
			getColor(*module.LabelColor).Print(*module.Label)
		} else if module.Label != nil {
			fmt.Print(*module.Label)
		}

		// Заполнение до края строки
		labelLen := 0
		if module.Label != nil {
			labelLen = len(*module.Label)
		}
		iconMargin := 0
		if module.Icon != nil && module.IconMarginRight != nil {
			iconMargin = *module.IconMarginRight + utf8.RuneCountInString(*module.Icon)
		}

		spaces := maxLen - labelLen - iconMargin + paddingRight
		fmt.Print(strings.Repeat(" ", spaces))

		fmt.Print("│" + strings.Repeat(" ", marginRight))

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
