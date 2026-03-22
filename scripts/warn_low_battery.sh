#!/bin/bash

# Seuil critique
THRESHOLD=35
BATTERY_PATH="/sys/class/power_supply/BAT0"
LEVEL=$(cat "$BATTERY_PATH/capacity")
STATUS=$(cat "$BATTERY_PATH/status")

if [ "$LEVEL" -le "$THRESHOLD" ] && [ "$STATUS" = "Discharging" ]; then
    # Envoi de la notification via dunstify (ou notify-send)
    # L'urgence "critical" permet de garder le message à l'écran selon votre config dunst
    dunstify -u critical "Batterie Critique" "Niveau actuel : ${LEVEL}%"
fi
