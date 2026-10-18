#!/bin/sh

##############
## Graphics ##
##############
xrdb -merge ~/.Xresources
brightnessctl set 50%
# dernier fond d'écran choisi (feh l'enregistre dans ~/.fehbg)
~/.fehbg &
# couleur de l'écran selon l'heure
redshift &

# for graphical effects: blur, round corner
picom &

############
## Locker ##
############
export locker="i3lock -e -i /home/alex/Images/Wallpapers/linux/thinkpad.png"
xset s 300
xss-lock -- $locker &

###################
## Notifications ##
###################
dunst &

#########
# Sound #
#########
mpd

###############
## Dashboard ##
###############
# overlay de widgets façon Mac OS X, démarré masqué ; Win Escape l'affiche (voir config.h)
~/.local/bin/osx-dashboard --hidden &

#########
## bar ##
#########
# define the status part (right-hand side) of the dwm bar
# go build -o dwm_bar bar.go
# to know the name of the network interface: ls /sys/class/net
~/softs/chadwm/scripts/dwm_bar wlp0s20f3 &

#######################
# finally, launch dwm #
#######################
while type chadwm >/dev/null; do chadwm && continue || break; done
