# TinyGo Arduino Projects

## Building

Requires a Linux or WSL system with `tinygo` installed (see [Getting Started](https://tinygo.org/getting-started/install/linux/)).

``` sh
make
```

For `go mod tidy`, run:

``` sh
make go-mod-tidy
```

## Flashing

Makefile invocation for flashing:

### Linux

``` sh
PORT=/dev/ttyACM0 make flash-countdown-lcd-1602-4bits
```
or, 

``` sh
make PORT=/dev/ttyACM0 flash-countdown-lcd-1602-4bits
```


### Windows

Install `make`.

``` sh
    winget install GnuWin32.Make
```

Add `C:\Program Files (x86)\GnuWin32\bin` to PATH.

Then, to flash:

``` sh
make PORT=COM7 flash-countdown-lcd-1602-4bits
```

### Simulations

| Project                              | Simulation |
|--------------------------------------|-------------------------------------------------------------------- |
| countdown-lcd-1602-4bits             | [wokwi.com](https://wokwi.com/projects/387688265780407297)          |

To flash on `wokwi.com`: `F1` then select `Upload Firmware and Start Simulation`.
