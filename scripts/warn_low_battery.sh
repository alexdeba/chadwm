#!/bin/bash

THRESHOLD=10
BATTERY_PATH="/sys/class/power_supply/BAT0"
LEVEL=$(cat "$BATTERY_PATH/capacity")
STATUS=$(cat "$BATTERY_PATH/status")

if [ "$LEVEL" -le "$THRESHOLD" ] && [ "$STATUS" = "Discharging" ]; then
    dunstify -u critical "Batterie Critique" "Niveau actuel : ${LEVEL}%"
fi
