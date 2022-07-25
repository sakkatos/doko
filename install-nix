#!/bin/sh

# This script installs the Nix package manager on your system by
# downloading a binary distribution and running its installer script
# (which in turn creates and populates /nix).

{ # Prevent execution if this script was only partially downloaded
oops() {
    echo "$0:" "$@" >&2
    exit 1
}

umask 0022

tmpDir="$(mktemp -d -t nix-binary-tarball-unpack.XXXXXXXXXX || \
          oops "Can't create temporary directory for downloading the Nix binary tarball")"
cleanup() {
    rm -rf "$tmpDir"
}
trap cleanup EXIT INT QUIT TERM

require_util() {
    command -v "$1" > /dev/null 2>&1 ||
        oops "you do not have '$1' installed, which I need to $2"
}

case "$(uname -s).$(uname -m)" in
    Linux.x86_64)
        hash=f163a846d6e4a340705b1642998817b72be8c1ffe785fb24629291ce4495588f
        path=q0mbh317fb7si22fbj0qhia0fc74wzjb/nix-2.8.1-x86_64-linux.tar.xz
        system=x86_64-linux
        ;;
    Linux.i?86)
        hash=b35d4e6d57d79a1c4e462bc65eb706bbc87719c701bbd09f6417a4a2f2e42450
        path=llyg5hc8vxz9ga4pkj1qp1cq5ms79ak8/nix-2.8.1-i686-linux.tar.xz
        system=i686-linux
        ;;
    Linux.aarch64)
        hash=d837ccfa9e9601a283e0cbe8e34a315469797be83c59ea26c318775f52dd1c7b
        path=dx42baipnv6pspx9yfxrx0khsqi3q6mi/nix-2.8.1-aarch64-linux.tar.xz
        system=aarch64-linux
        ;;
    Linux.armv6l_linux)
        hash=91e06f21a341e2e7d94d2f773a1b778f5ea605232856f581049fb4da952410d9
        path=fq5y4bh7m4xhcbqd0fbfqmad7083pn83/nix-2.8.1-armv6l-linux.tar.xz
        system=armv6l-linux
        ;;
    Linux.armv7l_linux)
        hash=e7052ce80e7af6d64b61a8bb2ed68dc0767aad430b1d3f1617700702060c9c61
        path=7kw7f7c24ixwlsah4j9a96bcffcpnxll/nix-2.8.1-armv7l-linux.tar.xz
        system=armv7l-linux
        ;;
    Darwin.x86_64)
        hash=7663a752fc9496ae976751ca57fbf6f81500be3a26b008dd13d0ef9b9c102ecc
        path=wjb3p7yacyk8sba7dbhs37sscwhcw4fq/nix-2.8.1-x86_64-darwin.tar.xz
        system=x86_64-darwin
        ;;
    Darwin.arm64|Darwin.aarch64)
        hash=482abb6bbe8840b69965acc00543c0bf5e413c66e6672ac4b2240dc325b95ade
        path=6mpz3hbjbjw2zrgak9hvg5j0pwrmgqpn/nix-2.8.1-aarch64-darwin.tar.xz
        system=aarch64-darwin
        ;;
    *) oops "sorry, there is no binary distribution of Nix for your platform";;
esac

# Use this command-line option to fetch the tarballs using nar-serve or Cachix
if [ "${1:-}" = "--tarball-url-prefix" ]; then
    if [ -z "${2:-}" ]; then
        oops "missing argument for --tarball-url-prefix"
    fi
    url=${2}/${path}
    shift 2
else
    url=https://releases.nixos.org/nix/nix-2.8.1/nix-2.8.1-$system.tar.xz
fi

tarball=$tmpDir/nix-2.8.1-$system.tar.xz

require_util tar "unpack the binary tarball"
if [ "$(uname -s)" != "Darwin" ]; then
    require_util xz "unpack the binary tarball"
fi

if command -v curl > /dev/null 2>&1; then
    fetch() { curl --fail -L "$1" -o "$2"; }
elif command -v wget > /dev/null 2>&1; then
    fetch() { wget "$1" -O "$2"; }
else
    oops "you don't have wget or curl installed, which I need to download the binary tarball"
fi

echo "downloading Nix 2.8.1 binary tarball for $system from '$url' to '$tmpDir'..."
fetch "$url" "$tarball" || oops "failed to download '$url'"

if command -v sha256sum > /dev/null 2>&1; then
    hash2="$(sha256sum -b "$tarball" | cut -c1-64)"
elif command -v shasum > /dev/null 2>&1; then
    hash2="$(shasum -a 256 -b "$tarball" | cut -c1-64)"
elif command -v openssl > /dev/null 2>&1; then
    hash2="$(openssl dgst -r -sha256 "$tarball" | cut -c1-64)"
else
    oops "cannot verify the SHA-256 hash of '$url'; you need one of 'shasum', 'sha256sum', or 'openssl'"
fi

if [ "$hash" != "$hash2" ]; then
    oops "SHA-256 hash mismatch in '$url'; expected $hash, got $hash2"
fi

unpack=$tmpDir/unpack
mkdir -p "$unpack"
tar -xJf "$tarball" -C "$unpack" || oops "failed to unpack '$url'"

script=$(echo "$unpack"/*/install)

[ -e "$script" ] || oops "installation script is missing from the binary tarball!"
export INVOKED_FROM_INSTALL_IN=1
"$script" "$@"

} # End of wrapping
