# Uberstatus

![byzanz-record -x 3065 -y 0 -h 22 -w 775 uberstatus.gif](doc/uberstatus.gif)

Status line generator for i3wm and (eventually) other WMs

Goals:

* integrate most of core functionality directly (checking cpu, network etc. without extra forking
* Support push and pull (polling) type of plugins
* Async (each plugin on its own timer) and sync mode with optional low power/slow update (WIP)

# Installation

Needs go installed in OS

Easy install:

    $ export GOPATH="/tmp/src/go" # skip if you have go already set up
    $ go get github.com/XANi/uberstatus
    $ mkdir -p ~/.config/uberstatus
    $ cp $GOPATH/src/github.com/XANi/uberstatus/cfg/uberstatus.default.conf ~/.config/uberstatus/uberstatus.conf # copy default config

Your binary will be in `$GOPATH/bin/uberstatus`. Copy it to non-temporary path or set GOPATH to something permanent before build

Dev install:

    $ export GOPATH="/tmp/src/go" # skip if you have go already set up
    $ export PATH="$PATH:$GOPATH/bin" # skip if you have go already set up
    $ go get github.com/mattn/gom
    $ make
    downloading github.com/op/go-logging
    downloading gopkg.in/yaml.v1
    downloading github.com/VividCortex/ewma
    gom exec go build -ldflags "-X main.version=0.0.1-0-gfbdadf3" [a-z]*go
    go fmt


# Operation

Each of plugins operates asynchronously and can send update of its status at any time. On each plugin update state is sent upstream (including previous cached plugin state) so it is possible to have each plugin update separetely.

Each plugin have a name (by default plugin name) and instance (in case of plugins that will be used multiple times, like network interfaces).


# Plugins

All plugins have `interval` option defining refresh rate in miliseconds

## CPU

## Clock

Parameters:

* `long_format`/`short_format` - time format, as golang formatting string
* `interval` - update interval in miliseconds. Note that it can be longer than it because of run time so for update every second something like 990 ms is better

## Disk free

```yaml
    - name: disk-root
      instance: df
      plugin: df
      config:
       prefix: "💾"
       mounts:
         - /
         - /var
         - /home
```

`prefix` will be added at beginning of the status bar. name and instance are just to distinguis between different instances

## Memory

Click to get detailed stats

## Network

Left click for interface's IP, right click for secondary IP (usually IPv6), middle click to display all addresses.

Parameters:

* `iface` - interface to use

## I3blocks

i3blocks-compatible input. It will also pass events in compatible way so plugins like volume can be used

Parameters:

* `command` - command to run
* `prefix` - prefix command with text

Example:

```yaml
    - name: volume
      plugin: i3blocks
      config:
        command: /path/to/i3blocks/volume
```

## Weather

Displays temperature, left click for more detailed weather data

To set it up first, get you token [here](https://openweathermap.org/), then get city name or id and set it in location

Example:

```yaml
    - name: weather
      plugin: weather
      config:
        openweather_api_key: 1111cccccccc11111111
        openweather_location: London
```

weather will automatically update every 10 minutes which is way below their free tier ratelimit
