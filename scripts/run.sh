#!/bin/sh

##############
## Graphics ##
##############
xrdb merge ~/.Xresources 
brightnessctl set 60%
feh --bg-fill ~/Images/Wallpapers/wall.png &

# for graphical effects: blur, round corner
picom &

############
## Locker ##
############
export locker="i3lock -e -i /home/alex/Images/Wallpapers/linux/thinkpad.png"
xset s 300
xss-lock -- $locker &

##############
## Keyboard ##
##############
# define keyboard auto-repeat (delay and speed)
xset r rate 300 50 &
# shift lock = escape
setxkbmap -option:escape

#########
## bar ##
#########
# define the status part (right-hand side) of the dwm bar
# go build -o dwm_bar bar.go
# to know the name of the network interface: ls /sys/class/net
~/softs/chadwm/scripts/dwm_bar wlp0s20f3 &

#############
## Widgets ##
#############
# eww
eww daemon &

# finally, launch dwm
while type chadwm >/dev/null; do chadwm && continue || break; done
