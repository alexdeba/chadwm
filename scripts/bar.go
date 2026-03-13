package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Définition des couleurs (Thème Tundra)
const (
	Black    = "^c#1e222a^"
	White    = "^c#abb2bf^"
	Grey     = "^b#282c34^"
	Green    = "^b#98c379^"
	Red      = "^b#e06c75^"
	Blue     = "^b#61afef^"
	DarkBlue = "^b#519aba^"
)

var (
	iconsDischarg = []string{"󰁺", "󰁻", "󰁼", "󰁽", "󰁾", "󰁿", "󰂀", "󰂁", "󰂂", "󰁹"}
	iconsCharging = []string{"󰢜", "󰂆", "󰂇", "󰂈", "󰢝", "󰂉", "󰢞", "󰂊", "󰂋", "󰂅"}
)

func getCPU() string {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return ""
	}
	val := strings.Split(string(data), " ")[0]
	return fmt.Sprintf("%s %s  %s %s %s%% ^b#1e222a^", Black, Green, White, Grey, val)
}

func getUpdates() string {
	// Adaptation Debian : compte les paquets pouvant être mis à jour via apt
	cmd := "apt list --upgradable 2>/dev/null | grep -c upgradable"
	out, _ := exec.Command("sh", "-c", cmd).Output()
	count := strings.TrimSpace(string(out))

	if count == "0" || count == "" {
		return fmt.Sprintf("  %s    Fully Updated", Green)
	}
	return fmt.Sprintf("  %s    %s updates", White, count)
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
	if index > 9 { index = 9 }

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
	data, _ := os.ReadFile("/sys/class/net/wlan0/operstate") // Ajuste le nom de l'interface si besoin
	state := strings.TrimSpace(string(data))
	
	if state == "up" {
		return fmt.Sprintf("%s %s 󰤨  ^d^^c#61afef^ Connected", Black, Blue)
	}
	return fmt.Sprintf("%s %s 󰤭  ^d^^c#61afef^ Disconnected", Black, Blue)
}

func getClock() string {
	t := time.Now().Format("02/01/2006 15:04")
	return fmt.Sprintf("%s %s 󱑆 %s%s %s  ", Black, DarkBlue, Black, Blue, t)
}

func main() {
	var updates string
	ticker := time.NewTicker(1 * time.Second)
	counter := 0

	for range ticker.C {
		// Mise à jour des paquets toutes les heures
		if counter%3600 == 0 {
			updates = getUpdates()
		}

		status := fmt.Sprintf("%s %s %s %s %s %s", 
			updates, getCPU(), getBattery(), getMem(), getWlan(), getClock())
		
		exec.Command("xsetroot", "-name", status).Run()
		counter++
	}
}
