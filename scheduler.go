package main

import (
	"fmt"
	"sync"
)

// Scheduler 进程调度器
type Scheduler struct {
	// 邻接表: 50个就绪队列，每个队列代表一个优先级(0-49)
	ReadyQueues [50]*PCB

	// pidUsage 记录PID使用情况，true表示可用，false表示已被使用
	pidUsage [101]bool

	// 当前运行进程
	RunningProcess *PCB

	// 所有进程列表(用于GUI显示)
	AllProcesses []*PCB

	// 互斥锁用于保护共享数据
	mu sync.Mutex

	// 是否停止调度
	quit bool
}

// NewScheduler 创建新的调度器
func NewScheduler() *Scheduler {
	s := &Scheduler{
		ReadyQueues:    [50]*PCB{},
		pidUsage:       [101]bool{},
		RunningProcess: nil,
		AllProcesses:   []*PCB{},
		quit:           false,
	}

	// 初始化PID使用情况，全部标记为可用
	for i := 1; i <= 100; i++ {
		s.pidUsage[i] = true
	}

	return s
}

// CreateProcess 创建新进程并添加到就绪队列
func (s *Scheduler) CreateProcess() *PCB {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 查找可用的PID
	pid := s.getAvailablePID()
	if pid == -1 {
		fmt.Println("无可用PID，进程创建失败")
		return nil
	}

	// 创建PCB
	pcb := NewPCB(pid)

	// 添加到就绪队列
	s.addToReadyQueue(pcb)

	// 添加到所有进程列表
	s.AllProcesses = append(s.AllProcesses, pcb)

	fmt.Printf("创建进程: %s\n", pcb)
	return pcb
}

// ScheduleProcess 进行一次进程调度
func (s *Scheduler) ScheduleProcess() *PCB {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 找到优先级最高的就绪进程
	highestPriorityPCB := s.getHighestPriorityProcess()

	if highestPriorityPCB == nil {
		s.RunningProcess = nil
		return nil
	}

	// 设置为运行状态
	highestPriorityPCB.Status = Run
	s.RunningProcess = highestPriorityPCB

	fmt.Printf("调度进程: %s\n", highestPriorityPCB)
	return highestPriorityPCB
}

// ProcessFinishedTimeSlice 处理完成一个时间片后的操作
func (s *Scheduler) ProcessFinishedTimeSlice() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.RunningProcess == nil {
		return
	}

	// 执行一个时间片后，优先级减半
	s.RunningProcess.Priority /= 2

	// 生命周期减1
	s.RunningProcess.Life--

	// 检查是否完成执行
	if s.RunningProcess.Life <= 0 {
		// 进程执行完毕，释放PID
		s.pidUsage[s.RunningProcess.Pid] = true

		// 从所有进程列表中移除
		s.removeFromAllProcesses(s.RunningProcess.Pid)

		fmt.Printf("进程完成: PID=%d\n", s.RunningProcess.Pid)
		s.RunningProcess = nil
	} else {
		// 进程未完成，重新插入就绪队列
		s.RunningProcess.Status = Ready
		s.addToReadyQueue(s.RunningProcess)
		s.RunningProcess = nil
	}
}

// SetQuit 设置退出标志
func (s *Scheduler) SetQuit(quit bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.quit = quit
}

// ShouldQuit 获取退出标志
func (s *Scheduler) ShouldQuit() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.quit
}

// 添加进程到就绪队列
func (s *Scheduler) addToReadyQueue(pcb *PCB) {
	priority := pcb.Priority
	if priority < 0 {
		priority = 0
	}
	if priority >= 50 {
		priority = 49
	}

	// 插入队列头部
	pcb.Next = s.ReadyQueues[priority]
	s.ReadyQueues[priority] = pcb
}

// 获取优先级最高的就绪进程
func (s *Scheduler) getHighestPriorityProcess() *PCB {
	// 从高优先级到低优先级遍历就绪队列
	for i := 49; i >= 0; i-- {
		if s.ReadyQueues[i] != nil {
			// 找到后从队列中移除
			pcb := s.ReadyQueues[i]
			s.ReadyQueues[i] = pcb.Next
			pcb.Next = nil
			return pcb
		}
	}
	return nil
}

// 从所有进程列表中移除指定PID的进程
func (s *Scheduler) removeFromAllProcesses(pid int) {
	for i, p := range s.AllProcesses {
		if p.Pid == pid {
			// 移除元素
			s.AllProcesses = append(s.AllProcesses[:i], s.AllProcesses[i+1:]...)
			break
		}
	}
}

// 获取可用的PID
func (s *Scheduler) getAvailablePID() int {
	for i := 1; i <= 100; i++ {
		if s.pidUsage[i] {
			s.pidUsage[i] = false // 标记为已使用
			return i
		}
	}
	return -1 // 无可用PID
}

// GetAllProcesses 获取所有进程的副本(用于GUI显示)
func (s *Scheduler) GetAllProcesses() []*PCB {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 创建副本以避免并发修改问题
	result := make([]*PCB, len(s.AllProcesses))
	copy(result, s.AllProcesses)

	return result
}
