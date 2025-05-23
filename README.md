# 进程调度模拟器

基于Go语言和Fyne GUI框架开发的多级反馈队列进程调度算法模拟器。此应用程序模拟操作系统如何管理和调度进程，展示进程状态变化、优先级调整和资源分配。

## 运行效果
![image](https://github.com/user-attachments/assets/7ea1dfdf-0c1f-4014-b1a6-9b42f7ae70dc)


## 项目结构

```
process-scheduling/
├── main.go              # 程序入口点
├── pcb.go               # 进程控制块(PCB)定义
├── scheduler.go         # 调度器实现
├── gui.go               # 图形用户界面
├── font.go              # 中文字体支持
├── go.mod               # Go模块文件
├── ProcessScheduler.exe # Windows可执行文件
└── README.md            # 项目说明文档
```

## 功能介绍

本模拟器实现了操作系统中的多级反馈队列进程调度算法，主要功能包括：

1. **进程创建与管理**
   - 动态创建具有随机优先级和生命周期的进程
   - 使用进程控制块(PCB)记录进程的属性和状态

2. **多级反馈队列调度**
   - 50个优先级队列(0-49)，数字越大优先级越高
   - 高优先级进程优先获取CPU执行
   - 进程执行一个时间片后，优先级减半

3. **进程生命周期管理**
   - 进程执行一个时间片后，生命周期减1
   - 生命周期为0时进程结束，释放PCB
   - 生命周期非0时进程返回就绪队列

4. **可视化界面**
   - 实时显示当前运行进程和就绪队列状态
   - 展示所有进程的PCB信息
   - 使用快捷键快速创建进程(Ctrl+F)和退出程序(Ctrl+Q)
   - 操作日志记录系统事件

## 编译和运行

### 前提条件

- Go 1.16或更高版本
- GCC编译器（Windows环境中可使用MinGW或MSYS2）
- Git

### 获取代码

```bash
git clone https://github.com/iriskl/process-scheduling.git
cd process-scheduling
```

### 安装依赖

```bash
go mod tidy
```

### 运行程序

```bash
go run .
```

### 编译可执行文件

#### 基本编译

```bash
go build -o ProcessScheduler.exe
```

#### 使用Fyne打包工具

首先安装Fyne命令行工具：

```bash
go install fyne.io/fyne/v2/cmd/fyne@latest
```

打包Windows可执行文件：

```bash
fyne package -os windows
```

## 操作说明

- 启动程序后，系统会自动创建5个初始进程
- 按下Ctrl+F可以创建新进程
- 按下Ctrl+Q可以退出程序
- 应用界面左侧显示运行状态、队列信息和操作日志
- 右侧显示所有进程的PCB信息
- 正在运行的进程会在列表中高亮显示

## 技术实现

- 使用Go协程(goroutine)并发处理进程调度
- 使用Go通道(channel)进行事件通信
- 使用Fyne框架构建跨平台GUI
- 采用邻接表表示多级反馈队列
- 使用互斥锁保护共享数据

## 常见问题

### Windows环境中字体显示乱码

问题可能是由于缺少中文字体支持导致的。程序会自动尝试加载系统中的中文字体，如果无法找到合适的字体，可以：

1. 确保系统中安装了常用中文字体（如微软雅黑、宋体等）
2. 重新启动程序，查看控制台输出了哪些字体加载信息

### 程序启动后没有窗口显示

这可能是由于Fyne框架对图形驱动的要求导致的。请确保：

1. 系统已安装最新的图形驱动
2. 如果使用的是远程连接，请确保支持图形转发
3. 在Windows环境中，可能需要安装MinGW工具链
