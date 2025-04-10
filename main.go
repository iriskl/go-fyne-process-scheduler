package main

import (
	"fmt"
	"math/rand"
	"time"

	"fyne.io/fyne/v2/app"
)

func main() {
	fmt.Println("程序开始运行...")

	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())
	fmt.Println("随机数种子已初始化")

	// 创建调度器
	scheduler := NewScheduler()
	fmt.Println("调度器已创建")

	// 创建Fyne应用
	fmt.Println("正在创建Fyne应用...")
	fyneApp := app.New()

	// 设置中文字体主题
	fyneApp.Settings().SetTheme(NewChineseTheme())
	fmt.Println("已设置中文字体支持")

	fmt.Println("Fyne应用已创建")

	// 创建GUI
	fmt.Println("正在创建GUI...")
	gui := NewProcessGUI(fyneApp, scheduler)
	fmt.Println("GUI已创建")

	// 初始化时先创建几个进程
	fmt.Println("初始化进程调度系统...")
	for i := 0; i < 5; i++ {
		p := scheduler.CreateProcess()
		fmt.Printf("创建初始进程: %s\n", p)
	}

	// 运行GUI
	fmt.Println("启动进程调度系统...")
	fmt.Println("按Ctrl+F创建新进程，按Ctrl+Q退出程序")

	// 启动GUI并进入主循环
	fmt.Println("启动GUI...")
	gui.Run()

	// 运行Fyne应用(会阻塞直到应用关闭)
	fmt.Println("进入Fyne主循环...")
	fyneApp.Run()

	fmt.Println("进程调度系统已关闭")
}
