#!/usr/bin/env bash

set -eo pipefail

export NIXPKGS_ALLOW_UNFREE=1
export unameOut="$(uname -s)"
export os_type=Unsupported

get_os_type() {
  case "${unameOut}" in
    Linux*) os_type=Linux ;;
    Darwin*) os_type=Mac ;;
    *) os_type="UNKNOWN:${unameOut}" ;;
  esac
}

download_nixpkg_manager_install_script() {
  rm -f install-nix
  curl -o install-nix https://releases.nixos.org/nix/nix-2.8.1/install
  chmod +x ./install-nix
}

configure_nix_flakes() {
  if [ ! -f ${HOME}/.config/nix/nix/nix.conf ]; then
    mkdir -p ${HOME}/.config/nix
    touch ${HOME}/.config/nix/nix.conf
  fi

  if ! grep -Fxq "experimental-features = nix-command flakes" ${HOME}/.config/nix/nix.conf; then
    echo 'experimental-features = nix-command flakes' >> ${HOME}/.config/nix/nix.conf
  fi
}

ensure_nix_profile() {
  if [ ! -f ${HOME}/.nix-profile/etc/profile.d/nix.sh ]; then
    user=$(whoami)
    ln -sf /nix/var/nix/profiles/per-user/${user}/profile ${HOME}/.nix-profile
  fi

  [ -f ${HOME}/.nix-profile/etc/profile.d/nix.sh ] && echo "Please start a new shell & execute '. ~/nix-profile/etc/profile.d/nix.sh'"
}

install_nixpkg_manager() {
  get_os_type
  if ! command -v nix >/dev/null 2>&1; then
    download_nixpkg_manager_install_script
    if [ "${os_type}" = "Mac" ]; then
      ./install-nix --darwin-use-unencrypted-nix-store-volume
    elif [ "${os_type}" = "Linux" ]; then
      ./install-nix --no-daemon
    else
      echo "Unsupported OS Platform for Nix development enviroment. Exiting!!!"
      exit 1
    fi
    rm -f ./install-nix
  fi
  configure_nix_flakes
}

install_nixpkg_manager

ensure_nix_profile

echo "------------------------------------------ CONGRATS !!! ----------------------------------------------------"
echo "  doko basic nix development environment tooling setup complete."
echo "  Please exit this shell terminal and start a new one in doughnut root directory then execute 'nix develop'."
echo "------------------------------------------    END       ----------------------------------------------------"
