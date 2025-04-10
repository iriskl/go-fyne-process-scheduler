package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/flopp/go-findfont"
)

// 自定义主题，支持中文字体
type ChineseTheme struct {
	fontData []byte
}

func NewChineseTheme() *ChineseTheme {
	return &ChineseTheme{
		fontData: loadFontData(),
	}
}

// 加载字体数据
func loadFontData() []byte {
	// 只尝试查找TTF格式的中文字体，避免TTC格式
	fontPaths := []string{
		"msyh.ttf",    // 微软雅黑
		"simhei.ttf",  // 黑体
		"simkai.ttf",  // 楷体
		"simfang.ttf", // 仿宋
	}

	var fontPath string
	var err error

	// 尝试找到一个中文TTF字体
	for _, font := range fontPaths {
		fontPath, err = findfont.Find(font)
		if err == nil && strings.HasSuffix(strings.ToLower(fontPath), ".ttf") {
			break
		}
	}

	// 如果找不到中文字体，尝试使用系统默认TTF字体
	if err != nil || !strings.HasSuffix(strings.ToLower(fontPath), ".ttf") {
		fmt.Println("无法找到TTF中文字体，尝试使用系统TTF字体...")

		// 常见的中文TTF字体路径（排除TTC格式）
		systemFontPaths := []string{
			"C:\\Windows\\Fonts\\msyh.ttf",     // 微软雅黑
			"C:\\Windows\\Fonts\\simhei.ttf",   // 黑体
			"C:\\Windows\\Fonts\\simkai.ttf",   // 楷体
			"C:\\Windows\\Fonts\\simfang.ttf",  // 仿宋
			"C:\\Windows\\Fonts\\STKAITI.TTF",  // 华文楷体
			"C:\\Windows\\Fonts\\STFANGSO.TTF", // 华文仿宋
			"C:\\Windows\\Fonts\\STHUPO.TTF",   // 华文琥珀
		}

		fontPath = ""
		for _, path := range systemFontPaths {
			if _, err := os.Stat(path); err == nil {
				fontPath = path
				break
			}
		}
	}

	// 如果找到了字体，加载它
	if fontPath != "" && strings.HasSuffix(strings.ToLower(fontPath), ".ttf") {
		fmt.Printf("使用字体: %s\n", fontPath)
		fontData, err := os.ReadFile(fontPath)
		if err == nil {
			return fontData
		} else {
			fmt.Printf("读取字体文件失败: %v\n", err)
		}
	}

	// 尝试使用内嵌的默认字体（Arial或其他拉丁字体）
	fmt.Println("警告: 无法加载中文TTF字体，将使用系统默认字体（中文可能显示为乱码）")
	return nil
}

// 在当前目录创建一个字体文件的副本
func createLocalFontCopy(fontData []byte) string {
	tempFontFile := filepath.Join(".", "chinese_font.ttf")
	err := os.WriteFile(tempFontFile, fontData, 0644)
	if err != nil {
		fmt.Printf("创建字体副本失败: %v\n", err)
		return ""
	}
	return tempFontFile
}

// 实现fyne.Theme接口
func (m *ChineseTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (m *ChineseTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *ChineseTheme) Font(style fyne.TextStyle) fyne.Resource {
	if len(m.fontData) > 0 {
		// 使用加载的中文字体
		res := &fyne.StaticResource{
			StaticName:    "chinese_font.ttf",
			StaticContent: m.fontData,
		}
		return res
	}
	return theme.DefaultTheme().Font(style)
}

func (m *ChineseTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
