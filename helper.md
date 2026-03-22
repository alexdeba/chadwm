
## Features added to chadwm
* scratchpads (using st): terminal, calc (with qalc), music (with ncmcpp)
* locker with i3lock
* Caps lock = escape
* dwm status bar in go

## Install
* install eww, qalc, st, alacritty
* install libs...
* compile dwm_bar
* update paths to chadwm

## Dwm 
workspace = tag = 1,2,3..

switch to tag : Win <tag>
move a window to another tag: Shift Win <tag>

Layout tile: Win t 
Layout maximize: Win Shift f 
Fullscreen: Win f 
Previous/next layout: Win Ctrl ,/. 

### layout tile

Shortcut: Win t 
Index: 0

    +------+----------------------------------+--------+
    | tags | title                            | status |
    +------+---------------------+------------+--------+
    |                            |                     |
    |                            |                     |
    |                            |                     |
    |                            |                     |
    |          master            |        stack        |
    |                            |                     |
    |                            |                     |
    |                            |                     |
    |                            |                     |
    +----------------------------+---------------------+

A new window appears in the master area.
Existing windows are pushed to the stack area.

Focus previous/next window: Win j/k

Increase amount of windows width area: Win h/l

Toggle window to master: Win Shift Enter 

### layout monocle [M] 

Shortcut: Win Shift f
Index: 1

All windows on top of each other.
Switch to previous: Win j/k

### layout bstack TTT

Shortcut: Win i
Index: 5

Tile layout but vertically

### layout centered master |M| 

Shortcut: Win Shift m 
Index: 11

