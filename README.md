# netlinx-language-server

<div align="center">
    <img align="center" src="./assets/img/NetLinx1.png" alt="netlinx-logo" width="150"/>
</div>

---

[![CI][ci]](https://github.com/Norgate-AV/netlinx-language-server/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/Norgate-AV/netlinx-language-server/branch/develop/graph/badge.svg)](https://codecov.io/gh/Norgate-AV/netlinx-language-server)
[![GitHub Release](https://img.shields.io/github/v/release/Norgate-AV/netlinx-language-server)](https://github.com/Norgate-AV/netlinx-language-server/releases)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)
[![GitHub contributors](https://img.shields.io/github/contributors/Norgate-AV/netlinx-language-server)](https://github.com/Norgate-AV/tree-sitter-netlinx/graphs/contributors)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

The `netlinx-language-server` is a Language Server Protocol (LSP) implementation for the NetLinx programming language. This project is currently under development.

[ci]: https://img.shields.io/github/actions/workflow/status/Norgate-AV/netlinx-language-server/ci.yml?logo=github&label=CI

## Contents :book:

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Installation :zap:](#installation-zap)
  - [macOS](#macos)
  - [Nix](#nix)
  - [Linux](#linux)
    - [Debian/Ubuntu](#debianubuntu)
    - [Fedora/RHEL/CentOS](#fedorarhelcentos)
    - [Alpine](#alpine)
    - [Arch (AUR)](#arch-aur)
    - [Snap](#snap)
  - [Windows](#windows)
    - [Scoop](#scoop)
    - [Chocolatey](#chocolatey)
    - [Winget](#winget)
  - [Direct Download](#direct-download)
  - [Build from Source](#build-from-source)
- [Team :soccer:](#team-soccer)
- [LICENSE :balance_scale:](#license-balance_scale)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation :zap:

### macOS

```sh
brew install norgate-av/tap/netlinx-language-server
```

### Nix

```sh
nix-env -iA nixpkgs.norgate-av.netlinx-language-server
```

### Linux

#### Debian/Ubuntu

```sh
# Add GPG key
curl -fsSL https://apt.norgate-av.com/public.key | sudo gpg --dearmor -o /usr/share/keyrings/norgate-av.gpg

# Add repository
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/norgate-av.gpg] https://apt.norgate-av.com stable main" | sudo tee /etc/apt/sources.list.d/norgate-av.list > /dev/null

# Update and install
sudo apt update
sudo apt install netlinx-language-server
```

#### Fedora/RHEL/CentOS

```sh
# Download the RPM package
sudo rpm -i netlinx-language-server_[version]_[arch].rpm
```

#### Alpine

```sh
# Download the APK package
sudo apk add --allow-untrusted ./netlinx-language-server_[version]_[arch].apk
```

#### Arch (AUR)

```sh
paru -S netlinx-language-server-bin
# or
yay -S netlinx-language-server-bin
```

#### Snap

```sh
snap install netlinx-language-server
```

### Windows

#### Scoop

```sh
scoop bucket add norgate-av https://github.com/Norgate-AV/scoop-bucket.git
scoop install netlinx-language-server
```

#### Chocolatey

```sh
choco install netlinx-language-server
```

#### Winget

```sh
winget install norgate-av.netlinx-language-server
```

### Direct Download

You can download the latest release directly from GitHub Releases.

1. Download the appropriate file for your operating system and architecture
2. Extract the archive
3. Move the binary to a location in your PATH

```sh
# Example for Linux/macOS
tar -xzf netlinx-language-server_Linux_x86_64.tar.gz
sudo mv netlinx-language-server /usr/local/bin/
```

### Build from Source

```sh
git clone https://github.com/Norgate-AV/netlinx-language-server.git
cd netlinx-language-server
make clean build
```

## Team :soccer:

This project is maintained by the following person(s) and a bunch of [awesome contributors](https://github.com/Norgate-AV/netlinx-language-server/graphs/contributors).

<table>
  <tr>
    <td align="center"><a href="https://github.com/damienbutt"><img src="https://avatars.githubusercontent.com/damienbutt?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Damien Butt</b></sub></a><br /></td>
  </tr>
</table>

## LICENSE :balance_scale:

[MIT](./LICENSE)
