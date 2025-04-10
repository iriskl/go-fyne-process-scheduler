package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ProcessGUI 表示进程调度器的图形用户界面
type ProcessGUI struct {
	App       fyne.App
	MainWin   fyne.Window
	Scheduler *Scheduler

	// GUI组件
	ProcessListContainer *fyne.Container // 进程列表容器
	StatusLabel          *widget.Label
	RunningProcInfo      *widget.Label // 当前运行进程详细信息
	QueueStatusInfo      *widget.Label // 队列状态信息
	OperationLog         *widget.Label // 操作日志
	TimeSliceLabel       *widget.Label // 时间片倒计时

	// 控制信号通道
	createProcessChan chan struct{}
}

// NewProcessGUI 创建新的GUI实例
func NewProcessGUI(app fyne.App, scheduler *Scheduler) *ProcessGUI {
	gui := &ProcessGUI{
		App:               app,
		MainWin:           app.NewWindow("多级反馈队列进程调度模拟"),
		Scheduler:         scheduler,
		createProcessChan: make(chan struct{}, 10),
	}

	// 设置窗口大小
	gui.MainWin.Resize(fyne.NewSize(900, 700))

	// 初始化界面组件
	gui.setupUI()

	// 设置键盘快捷键
	gui.setupShortcuts()

	return gui
}

// 设置用户界面
func (g *ProcessGUI) setupUI() {
	// 状态标签
	g.StatusLabel = widget.NewLabel("进程调度模拟系统 - 按Ctrl+F创建进程, Ctrl+Q退出")
	g.StatusLabel.TextStyle = fyne.TextStyle{Bold: true}

	// 当前运行进程信息
	g.RunningProcInfo = widget.NewLabel("当前无运行进程")
	g.RunningProcInfo.TextStyle = fyne.TextStyle{Bold: true}

	// 队列状态信息
	g.QueueStatusInfo = widget.NewLabel("就绪队列: 0个进程")

	// 操作日志
	g.OperationLog = widget.NewLabel("系统启动...")
	g.OperationLog.Wrapping = fyne.TextWrapWord

	// 时间片倒计时
	g.TimeSliceLabel = widget.NewLabel("下一次调度: 准备中")

	// 进程列表标题
	allProcessesTitle := canvas.NewText("所有进程PCB列表", theme.ForegroundColor())
	allProcessesTitle.TextStyle = fyne.TextStyle{Bold: true}

	// 创建进程列表容器
	g.ProcessListContainer = container.NewVBox()

	// 使进程列表可滚动
	processListScroll := container.NewScroll(g.ProcessListContainer)
	processListScroll.SetMinSize(fyne.NewSize(300, 400))

	// 更新进程列表显示
	g.updateProcessList()

	// 创建按钮
	createBtn := widget.NewButtonWithIcon("创建进程 (Ctrl+F)", theme.ContentAddIcon(), func() {
		g.createProcessChan <- struct{}{}
	})

	// 退出按钮
	quitBtn := widget.NewButtonWithIcon("退出 (Ctrl+Q)", theme.CancelIcon(), func() {
		g.Scheduler.SetQuit(true)
	})

	// 操作按钮布局
	buttons := container.NewHBox(createBtn, quitBtn)

	// 创建分割线
	divider := widget.NewSeparator()

	// 创建信息面板
	infoBox := container.NewVBox(
		widget.NewLabelWithStyle("运行状态", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		g.RunningProcInfo,
		widget.NewSeparator(),
		g.TimeSliceLabel,
		widget.NewSeparator(),
		g.QueueStatusInfo,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("操作日志", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		g.OperationLog,
	)

	// 右侧面板（进程列表）
	rightPanel := container.NewVBox(
		allProcessesTitle,
		processListScroll,
	)

	// 整体布局 - 使用左右分割
	splitContainer := container.NewHSplit(
		container.NewVBox(g.StatusLabel, buttons, divider, infoBox),
		rightPanel,
	)
	// 调整分割比例，确保右侧有足够空间
	splitContainer.Offset = 0.3

	// 设置主窗口内容
	g.MainWin.SetContent(splitContainer)
}

// 更新进程列表显示
func (g *ProcessGUI) updateProcessList() {
	// 获取所有进程
	processes := g.Scheduler.GetAllProcesses()

	// 清空列表容器
	g.ProcessListContainer.RemoveAll()

	// 没有进程时显示提示
	if len(processes) == 0 {
		g.ProcessListContainer.Add(widget.NewLabel("当前没有进程"))
		return
	}

	// 创建每个进程的显示项
	for _, p := range processes {
		var procLabel *widget.Label

		// 高亮显示当前运行的进程
		if p.Status == Run {
			text := fmt.Sprintf("→ PID: %d | 状态: %s | 优先级: %d | 生命周期: %d",
				p.Pid, p.Status, p.Priority, p.Life)
			procLabel = widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		} else {
			text := fmt.Sprintf("PID: %d | 状态: %s | 优先级: %d | 生命周期: %d",
				p.Pid, p.Status, p.Priority, p.Life)
			procLabel = widget.NewLabel(text)
		}

		// 添加到容器
		g.ProcessListContainer.Add(procLabel)
	}

	// 刷新容器
	g.ProcessListContainer.Refresh()
	fmt.Printf("更新进程列表: 共%d个进程\n", len(processes))
}

// 设置键盘快捷键
func (g *ProcessGUI) setupShortcuts() {
	// Ctrl+F: 创建进程
	g.MainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyF,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		fmt.Println("触发快捷键: Ctrl+F - 创建新进程")
		g.createProcessChan <- struct{}{}
	})

	// Ctrl+Q: 退出程序
	g.MainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		fmt.Println("触发快捷键: Ctrl+Q - 退出程序")
		g.Scheduler.SetQuit(true)
	})
}

// Run 启动GUI并进入调度循环
func (g *ProcessGUI) Run() {
	// 显示主窗口
	g.MainWin.Show()

	// 启动调度协程
	go g.scheduleLoop()

	// 启动一个协程来处理进程创建请求
	go func() {
		for range g.createProcessChan {
			pcb := g.Scheduler.CreateProcess()
			if pcb != nil {
				g.addOperationLog(fmt.Sprintf("创建新进程: PID=%d, 优先级=%d, 生命周期=%d",
					pcb.Pid, pcb.Priority, pcb.Life))
				fmt.Printf("已创建新进程: %s\n", pcb) // 调试信息
			}
			g.updateQueueStatus()
			g.updateProcessList() // 更新进程列表
		}
	}()
}

// 添加操作日志
func (g *ProcessGUI) addOperationLog(message string) {
	currentText := g.OperationLog.Text
	// 保持日志不超过5行
	currentLines := 1
	for i, c := range currentText {
		if c == '\n' {
			currentLines++
		}
		if currentLines >= 5 {
			currentText = currentText[:i]
			break
		}
	}
	g.OperationLog.SetText(message + "\n" + currentText)
}

// 更新队列状态信息
func (g *ProcessGUI) updateQueueStatus() {
	// 计算各优先级队列的进程数量
	queueCounts := make(map[int]int)
	totalReady := 0

	for i := 0; i < 50; i++ {
		count := 0
		for p := g.Scheduler.ReadyQueues[i]; p != nil; p = p.Next {
			count++
		}
		if count > 0 {
			queueCounts[i] = count
			totalReady += count
		}
	}

	// 构建队列状态文本
	statusText := fmt.Sprintf("就绪队列: 共%d个进程\n", totalReady)

	// 最多显示前5个非空队列
	count := 0
	for i := 49; i >= 0; i-- {
		if queueCounts[i] > 0 {
			statusText += fmt.Sprintf("优先级%d: %d个进程\n", i, queueCounts[i])
			count++
			if count >= 5 {
				break
			}
		}
	}

	g.QueueStatusInfo.SetText(statusText)
}

// 调度循环
func (g *ProcessGUI) scheduleLoop() {
	const timeSlice = 2 * time.Second // 时间片大小

	for !g.Scheduler.ShouldQuit() {
		// 更新队列状态
		g.updateQueueStatus()

		// 执行调度
		pcb := g.Scheduler.ScheduleProcess()

		// 更新进程列表
		g.updateProcessList()

		if pcb != nil {
			// 显示时间片开始信息
			startMsg := fmt.Sprintf("开始执行进程: PID=%d, 优先级=%d, 生命周期=%d",
				pcb.Pid, pcb.Priority, pcb.Life)
			g.addOperationLog(startMsg)

			// 更新状态标签
			g.StatusLabel.SetText("进程调度系统 - 正在运行进程")

			// 更新当前运行进程信息
			g.RunningProcInfo.SetText(fmt.Sprintf(
				"当前运行进程: PID=%d\n优先级=%d\n生命周期=%d\n状态=%s",
				pcb.Pid, pcb.Priority, pcb.Life, pcb.Status))

			// 时间片倒计时
			for i := int(timeSlice / time.Millisecond); i > 0; i -= 100 {
				g.TimeSliceLabel.SetText(fmt.Sprintf("当前时间片剩余: %.1f秒", float64(i)/1000.0))
				time.Sleep(100 * time.Millisecond)
				if g.Scheduler.ShouldQuit() {
					break
				}
			}

			// 处理时间片结束后的操作
			oldPriority := pcb.Priority
			oldLife := pcb.Life
			g.Scheduler.ProcessFinishedTimeSlice()

			// 记录操作日志
			if oldLife > 1 {
				g.addOperationLog(fmt.Sprintf(
					"进程PID=%d时间片结束: 优先级%d→%d, 生命周期%d→%d, 重新进入就绪队列",
					pcb.Pid, oldPriority, oldPriority/2, oldLife, oldLife-1))
			} else {
				g.addOperationLog(fmt.Sprintf("进程PID=%d执行完毕，已撤销PCB", pcb.Pid))
			}

			// 更新当前运行进程信息
			g.RunningProcInfo.SetText("当前无运行进程")
			g.TimeSliceLabel.SetText("等待下一次调度...")

			// 更新进程列表
			g.updateProcessList()
		} else {
			// 没有进程可调度
			g.StatusLabel.SetText("进程调度系统 - 空闲")
			g.RunningProcInfo.SetText("当前无运行进程")
			g.TimeSliceLabel.SetText("等待进程...")

			// 等待一段时间
			time.Sleep(500 * time.Millisecond)
		}
	}

	// 关闭应用
	g.App.Quit()
}
