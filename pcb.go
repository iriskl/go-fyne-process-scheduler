package main

import (
	"fmt"
	"math/rand"
)

// 进程状态常量
const (
	Ready = "ready" // 就绪状态
	Run   = "run"   // 运行状态
)

// PCB 进程控制块结构
type PCB struct {
	Pid      int    // 进程标识符(1-100)
	Status   string // 进程状态(ready/run)
	Priority int    // 进程优先级(0-49)
	Life     int    // 进程生命周期(1-5)
	Next     *PCB   // 进程队列指针
}

// 创建新PCB实例
func NewPCB(pid int) *PCB {
	return &PCB{
		Pid:      pid,
		Status:   Ready,
		Priority: rand.Intn(50),    // 0-49的随机整数
		Life:     rand.Intn(5) + 1, // 1-5的随机整数
		Next:     nil,
	}
}

// String 返回PCB的字符串表示
func (p *PCB) String() string {
	return fmt.Sprintf("PID: %d, 状态: %s, 优先级: %d, 生命周期: %d", p.Pid, p.Status, p.Priority, p.Life)
}

// 生成一个在指定范围内的随机整数
func randInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}
