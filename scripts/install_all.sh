#!/bin/bash

# Mettre à jour la liste des paquets
sudo apt-get update

# Installation groupée (triée alphabétique)
sudo apt-get install -y --ignore-missing \
    acpi acpid alsa-utils asciidoctor atool autoconf automake autofs \
    bluez-tools boost-all-dev build-essential \
    cifs-utils clang cmake compton console-data cppman curl \
    default-jre dos2unix dosbox duf \
    elinks evince exiftool \
    fbreader feh fonts-font-awesome frotz fzf fvwm-icons \
    gdb gnome-keyring \
    htop \
    imagemagick \
    jmtpfs jq \
    keepass2 \
    libev-dev libmtp-dev libnotify-bin libsdl2-mixer-2.0-0 libspdlog-dev libx11-xcb-dev \
    man-db mcomix mpd mtp-tools mwm \
    ncdu ncmpcpp nethack-console nethack-x11 neofetch ncurses-dev ntfs-3g \
    p7zip pavucontrol pavumeter pmount powerline procps pulseaudio python3 python3-pip \
    ranger redshift rsync rxvt-unicode-256color \
    scummvm sdcv software-properties-common \
    tig tlp tmux \
    ufw unrar unrar-free universal-ctags \
    vlc \
    wine \
    x11-apps xboard xclip xcolorsel xpdf xterm yank \
    chromium

# Installation de Homebrew et ses outils
if ! command -v brew &> /dev/null; then
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    # Assure-toi que brew est dans ton PATH après l'installation
    eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
fi

brew install glow
brew install jstkdng/programs/ueberzugpp

# Installation Go
if [ ! -d "$HOME/softs" ]; then mkdir -p "$HOME/softs"; fi
cd ~/softs
if [ ! -d "update-golang" ]; then
    git clone https://github.com/udhos/update-golang
    cd update-golang
    sudo ./update-golang.sh
fi

# Installation Rust
if ! command -v rustup &> /dev/null; then
    curl https://sh.rustup.rs -sSf | sh -s -- -y
fi

# Permissions Docker
sudo usermod -aG docker $USER

echo "Installation terminée avec succès !"
