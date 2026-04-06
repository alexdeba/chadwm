package main

import (
	"fmt"
	"os"
	"os/exec"
  "path/filepath"
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
)

var (
	iconsDischarg = []string{"󰁺", "󰁻", "󰁼", "󰁽", "󰁾", "󰁿", "󰂀", "󰂁", "󰂂", "󰁹"}
	iconsCharging = []string{"󰢜", "󰂆", "󰂇", "󰂈", "󰢝", "󰂉", "󰢞", "󰂊", "󰂋", "󰂅"}
	network_name  string
)

func getCPU() string {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return ""
	}
	val := strings.Split(string(data), " ")[0]
	return fmt.Sprintf("%s %s     %s %s %s%% ^b#1e222a^", Black, Green, White, Grey, val)
}

func getBattery() string {
	capData, _ := os.ReadFile("/sys/class/power_supply/BAT0/capacity")
	statData, _ := os.ReadFile("/sys/class/power_supply/BAT0/status")

	capacity := strings.TrimSpace(string(capData))
	status := strings.TrimSpace(string(statData))

	var val int
	fmt.Sscanf(capacity, "%d", &val)

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

	return fmt.Sprintf("%s %s  %s %s %s %d%% ^b#1e222a^", Black, Red, icon, White, Grey, val)
}

func getMem() string {
	cmd := "free -h | awk '/^Mem/ { print $3 }' | sed s/i//g"
	out, _ := exec.Command("sh", "-c", cmd).Output()
	usage := strings.TrimSpace(string(out))
	return fmt.Sprintf("%s %s  %s %s  %s ^b#1e222a^", Black, Green, White, Grey, usage)
}

func getWlan() string {
	data, _ := os.ReadFile(fmt.Sprintf("/sys/class/net/%s/operstate", network_name))
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

func main() {

	// Check if a parameter is present
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./dwm_bar <network_interface>")
		return
	}

	network_name = os.Args[1]

	// get the path to this executable
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	binDir := filepath.Dir(ex)

	// 3. Construire le chemin vers le script
	scriptPath := filepath.Join(binDir, "./warn_low_battery.sh")
	cmd := exec.Command("sh", scriptPath)

	ticker1s := time.NewTicker(1 * time.Second)
	ticker10s := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ticker1s.C:
			status := fmt.Sprintf("%s %s %s %s %s",
			getCPU(), getBattery(), getMem(), getWlan(), getClock())
			exec.Command("xsetroot", "-name", status).Run()

		case <-ticker10s.C:
			cmd.Run()
		}
	}

}
