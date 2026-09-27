package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	// foreground colors
	Black = "^c#1e222a^"
	White = "^c#abb2bf^"
	// background colors
	Grey     = "^b#282c34^"
	Green    = "^b#98c379^"
	Red      = "^b#e06c75^"
	Blue     = "^b#61afef^"
	DarkBlue = "^b#519aba^"
	BlackBg  = "^b#1e222a^"
)

const (
	batteryPath      = "/sys/class/power_supply/BAT0"
	batteryThreshold = 10 // percent, below which a notification is sent
	refreshInterval  = 2 * time.Second
	batteryInterval  = 10 * time.Second
)

var (
	iconsDischarg = []string{"󰁺", "󰁻", "󰁼", "󰁽", "󰁾", "󰁿", "󰂀", "󰂁", "󰂂", "󰁹"}
	iconsCharging = []string{"󰢜", "󰂆", "󰂇", "󰂈", "󰢝", "󰂉", "󰢞", "󰂊", "󰂋", "󰂅"}
	networkName   string

	// previous /proc/stat counters, used to compute the CPU usage between two ticks
	// (zero on the first call, which then gives the average usage since boot)
	prevIdle, prevTotal uint64
	// true once the low battery notification has been sent, reset when charging or above threshold
	batteryWarned bool
)

// readCPUStat returns the idle (idle + iowait) and total jiffies of the first line of /proc/stat
func readCPUStat() (idle, total uint64, ok bool) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0, false
	}
	line, _, _ := strings.Cut(string(data), "\n")
	fields := strings.Fields(line)
	if len(fields) < 5 || fields[0] != "cpu" {
		return 0, 0, false
	}
	for i, f := range fields[1:] {
		v, err := strconv.ParseUint(f, 10, 64)
		if err != nil {
			return 0, 0, false
		}
		total += v
		if i == 3 || i == 4 { // idle, iowait
			idle += v
		}
	}
	return idle, total, true
}

func getCPU() string {
	idle, total, ok := readCPUStat()
	if !ok {
		return ""
	}
	usage := 0.0
	if dTotal := total - prevTotal; dTotal > 0 {
		usage = 100 * (1 - float64(idle-prevIdle)/float64(dTotal))
	}
	prevIdle, prevTotal = idle, total
	return fmt.Sprintf("%s %s     %s %s %.0f%% %s", Black, Green, White, Grey, usage, BlackBg)
}

// readBattery returns the capacity (in %) and the status of the battery
func readBattery() (int, string, bool) {
	capData, err := os.ReadFile(batteryPath + "/capacity")
	if err != nil {
		return 0, "", false
	}
	statData, err := os.ReadFile(batteryPath + "/status")
	if err != nil {
		return 0, "", false
	}
	val, err := strconv.Atoi(strings.TrimSpace(string(capData)))
	if err != nil {
		return 0, "", false
	}
	return val, strings.TrimSpace(string(statData)), true
}

func getBattery() string {
	val, status, ok := readBattery()
	if !ok {
		return ""
	}

	index := 0
	if val > 0 {
		index = (val - 1) / 10
	}
	if index > 9 {
		index = 9
	}

	icon := iconsDischarg[index]
	if status == "Charging" || status == "Full" {
		icon = iconsCharging[index]
	}

	return fmt.Sprintf("%s %s  %s %s %s %d%% %s", Black, Red, icon, White, Grey, val, BlackBg)
}

// checkBattery sends a notification once when the battery goes below the threshold while discharging
func checkBattery() {
	val, status, ok := readBattery()
	if !ok {
		return
	}
	if val > batteryThreshold || status != "Discharging" {
		batteryWarned = false
		return
	}
	if batteryWarned {
		return
	}
	batteryWarned = true
	// run in the background so that a hanging notification daemon does not freeze the bar
	go exec.Command("dunstify", "-u", "critical", "Batterie Critique",
		fmt.Sprintf("Niveau actuel : %d%%", val)).Run()
}

func getMem() string {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return ""
	}
	var total, available uint64
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "MemTotal:":
			total, _ = strconv.ParseUint(fields[1], 10, 64)
		case "MemAvailable:":
			available, _ = strconv.ParseUint(fields[1], 10, 64)
		}
	}
	if total == 0 {
		return ""
	}

	// values are in kB
	used := float64(total-available) / 1024
	usage := fmt.Sprintf("%.0fM", used)
	if used >= 1024 {
		usage = fmt.Sprintf("%.1fG", used/1024)
	}
	return fmt.Sprintf("%s %s  %s %s  %s %s", Black, Green, White, Grey, usage, BlackBg)
}

func getWlan() string {
	data, _ := os.ReadFile("/sys/class/net/" + networkName + "/operstate")
	state := strings.TrimSpace(string(data))

	if state == "up" {
		return fmt.Sprintf("%s %s 󰤨  ^d^^c#61afef^ Up", Black, Blue)
	}
	return fmt.Sprintf("%s %s 󰤭  ^d^^c#61afef^ Down", Black, Blue)
}

func getClock() string {
	t := time.Now().Format("02/01/2006 15:04")
	return fmt.Sprintf("%s %s 󱑆 %s%s %s  ", Black, DarkBlue, Black, Blue, t)
}

func updateStatus() {
	status := fmt.Sprintf("%s %s %s %s %s",
		getCPU(), getBattery(), getMem(), getWlan(), getClock())
	if err := exec.Command("xsetroot", "-name", status).Run(); err != nil {
		// the X server is gone (e.g. after logout): stop instead of looping forever
		fmt.Fprintln(os.Stderr, "dwm_bar: xsetroot failed:", err)
		os.Exit(1)
	}
}

func main() {
	// Check if a parameter is present
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./dwm_bar <network_interface>")
		return
	}
	networkName = os.Args[1]

	// first update right away, so that the bar is not empty at startup
	updateStatus()
	checkBattery()

	refreshTicker := time.NewTicker(refreshInterval)
	batteryTicker := time.NewTicker(batteryInterval)

	for {
		select {
		case <-refreshTicker.C:
			updateStatus()
		case <-batteryTicker.C:
			checkBattery()
		}
	}
}
