package types

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
	"tinyfetch/internal/config"
	"tinyfetch/internal/utils/logger"

	"github.com/fatih/color"
)

func GetTypeInfo(module config.Module) string {
	switch *module.Type {
		case "user": return UserGetInfo()
		case "hostname": return HostnameGetInfo()
		case "os": return OSGetInfo()
		case "kernel": return KernelGetInfo()
		case "uptime": return UptimeGetInfo()
		case "shell": return ShellGetInfo() 
		case "packages": return PackagesGetInfo()
		case "memory": return MemoryGetInfo()
		case "colors": return ColorsGetInfo()
		case "command": return CommandGetInfo(module)
		default: return ""
	}
}


func CommandGetInfo(module config.Module) string{
    cmd := exec.Command("sh", "-c", *module.Script)

    out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Fatal("failed to execute script: %s", *module.Script)
	}
    return strings.TrimSuffix(string(out), "\n")
}


func UserGetInfo() string {
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}
	return user
}

func HostnameGetInfo() string {
	host, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return host
}

func OSGetInfo() string {
	out, err := exec.Command("cat", "/etc/os-release").Output()
	if err != nil {
		return runtime.GOOS
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(line[len("PRETTY_NAME="):], `"`)
		}
	}
	return runtime.GOOS
}

func KernelGetInfo() string {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func UptimeGetInfo() string {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "unknown"
	}
	seconds, err := strconv.ParseFloat(strings.Fields(string(data))[0], 64)
	if err != nil {
		return "unknown"
	}
	duration := time.Duration(seconds) * time.Second
	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 || days > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	parts = append(parts, fmt.Sprintf("%dm", minutes))

	return strings.Join(parts, " ")
}

func ShellGetInfo() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "unknown"
	}
	parts := strings.Split(shell, "/")
	return parts[len(parts)-1]
}

func PackagesGetInfo() string {
	out, err := exec.Command("bash", "-c", "pacman -Q | wc -l").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func MemoryGetInfo() string {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return "unknown"
	}
	var memTotal, memFree, buffers, cached int
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		value, _ := strconv.Atoi(fields[1])
		switch fields[0] {
		case "MemTotal:":
			memTotal = value
		case "MemFree:":
			memFree = value
		case "Buffers:":
			buffers = value
		case "Cached:":
			cached = value
		}
	}
	used := memTotal - (memFree + buffers + cached)
	return fmt.Sprintf("%d | %d Mib", used/1024, memTotal/1024)
}

func ColorsGetInfo() string {
	colors := []*color.Color{
		color.New(color.FgWhite),
		color.New(color.FgRed),
		color.New(color.FgYellow),
		color.New(color.FgGreen),
		color.New(color.FgCyan),
		color.New(color.FgBlue),
		color.New(color.FgMagenta),
		color.New(color.FgBlack),
	}

	result := ""
	for _, c := range colors {
		result += c.Sprint(" ")
	}
	return result
}