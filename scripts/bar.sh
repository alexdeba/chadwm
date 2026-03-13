#!/bin/bash

# ^c$var^ = fg color
# ^b$var^ = bg color

interval=0

icons_discharg=("󰁺 " "󰁻 " "󰁼 " "󰁽 " "󰁾 " "󰁿 " "󰂀 " "󰂁 " "󰂂 " "󰁹 ")
icons_charging=("󰢜 " "󰂆 " "󰂇 " "󰂈 " "󰢝 " "󰂉 " "󰢞 " "󰂊 " "󰂋 " "󰂅 ")

# load colors
. ~/softs/chadwm/scripts/bar_themes/tundra

cpu() {
  cpu_val=$(grep -o "^[^ ]*" /proc/loadavg)
  printf "^c$black^ ^b$green^  "
  printf "^c$white^ ^b$grey^ $cpu_val%% ^b$black^"
}

pkg_updates() {
  #updates=$({ timeout 20 doas xbps-install -un 2>/dev/null || true; } | wc -l) # void
  updates=$({ timeout 20 checkupdates 2>/dev/null || true; } | wc -l) # arch
  # updates=$({ timeout 20 aptitude search '~U' 2>/dev/null || true; } | wc -l)  # apt (ubuntu, debian etc)

  if [ -z "$updates" ]; then
    printf "  ^c$green^    Fully Updated"
  else
    printf "  ^c$white^    $updates"" updates"
  fi
}

battery() {
	bat_val="$(cat /sys/class/power_supply/BAT0/capacity)"
	index=$(( $bat_val == 0 ? 0 : ($bat_val - 1) / 10 ))
	case "$(cat /sys/class/power_supply/BAT0/status 2>/dev/null)" in
	"Full")
    bat="${icons_charging[9]}"
		;;
	"Charging")
    bat="${icons_charging[$index]}"
		;;
	"Discharging")
    bat="${icons_discharg[$index]}"
		;;
	esac
		printf "^c$black^ ^b$red^  $bat"
		printf "^c$white^ ^b$grey^ $bat_val%% ^b$black^"
}

brightness() {
  printf "^c$red^   "
  printf "^c$red^%.0f\n" $(cat /sys/class/backlight/*/brightness)
}

mem() {
  printf "^c$black^ ^b$green^  "
  printf "^c$white^ ^b$grey^  $(free -h | awk '/^Mem/ { print $3 }' | sed s/i//g) ^b$black^"
}

wlan() {
	case "$(cat /sys/class/net/wl*/operstate 2>/dev/null)" in
	up)
		printf "^c$black^ ^b$blue^ 󰤨  ^d^%s"
		printf " ^c$blue^Connected"
		;;
	down)
		printf "^c$black^ ^b$blue^ 󰤭  ^d^%s"
		printf " ^c$blue^Disconnected"
		;;
	esac
}

clock() {
	printf "^c$black^ ^b$darkblue^ 󱑆 "
	printf "^c$black^^b$blue^ $(date '+%d/%m/%Y %H:%M')  "
}

while true; do

  [ $interval = 0 ] || [ $(($interval % 3600)) = 0 ] && updates=$(pkg_updates)
  interval=$((interval + 1))

	sleep 1 && xsetroot -name "$(pkg_updates) $(cpu) $(battery) $(mem) $(wlan) $(clock)"
  echo "$(cpu) $(battery) $(mem) $(wlan) $(clock)"
done
