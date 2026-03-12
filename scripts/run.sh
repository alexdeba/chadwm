#!/bin/sh

xrdb merge ~/.Xresources 
xbacklight -set 10 &
feh --bg-fill ~/Images/Wallpapers/wall.png &

## Locker ##
export locker="i3lock -e -i /home/alex/Images/Wallpapers/linux/thinkpad.png"
xset s 300
xss-lock -- $locker &

# define keyboard auto-repeat (delay and speed)
xset r rate 300 50 &
setxkbmap -option:escape

# for graphical effects
picom &

# define the status part (right-hand side) of the dwm bar
dash ~/softs/chadwm/scripts/bar.sh &

# eww
eww daemon &

# finally, launch dwm
while type chadwm >/dev/null; do chadwm && continue || break; done
