<div align="center">

<img src="https://i.kagchi.my.id/nezuko.png" alt="Logo" width="200px" height="200px" style="border-radius:50%"/>

# @nezuchan/image-proxy

**A fast Image proxy service, such as resizing written in go.**

[![GitHub](https://img.shields.io/github/license/nezuchan/cordis-brokers)](https://github.com/nezuchan/cordis-brokers/blob/main/LICENSE)
[![Discord](https://discordapp.com/api/guilds/785715968608567297/embed.png)](https://nezu.my.id)

</div>

## Requirements

-   [libvips](https://github.com/libvips/libvips) 8.10+
-   C compatible compiler such as gcc 4.6+ or clang 3.0+
-   Go 1.19+

## Dependencies for govips

### MacOS

Use [homebrew](https://brew.sh/) to install vips and pkg-config:

```bash
brew install vips pkg-config
```

### Ubuntu

You need a recent libvips to work with govips. New govips functionality is continuously added which takes advantage of new libvips functionality. Groovy (20.10) and Hirsute (21.04) repositories have working versions. However on Focal (20.04), you need to install libvips and dependencies from a backports repository:

```bash
sudo add-apt-repository -y ppa:strukturag/libde265
sudo add-apt-repository -y ppa:strukturag/libheif
sudo add-apt-repository ppa:tonimelisma/ppa
```

Then:

```bash
sudo apt -y install libvips-dev
```

### Windows

The recommended approach on Windows is to use Govips via WSL and Ubuntu.

If you need to run Govips natively on Windows, it's not difficult but will require some effort. We don't have a recommended environment or setup at the moment. Windows is also not in our list of CI/CD targets so Govips is not regularly tested for compatibility. If you would be willing to setup and maintain a robust CI/CD Windows environment, please open a PR, we would be pleased to accept your contribution and support Windows as a platform.

## Installation

```bash
go get -u github.com/davidbyttow/govips/v2/vips
```

### MacOS note

On MacOS, govips may not compile without first setting an environment variable:

```bash
export CGO_CFLAGS_ALLOW="-Xpreprocessor"
```


# Features
- Secure image routing, this service ideal when you are trying to hide origin request. using `aes-256-cbc` encryption for encrypting origin url
- Docker ready
- Production Ready