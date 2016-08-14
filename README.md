# Uberstatus

![](doc/uberstatus.gif)

Status line generator for i3wm and (eventually) other WMs

Goals:

* integrate most of core functionality directly (checking cpu, network etc. without extra forking
* Support push and pull (polling) type of plugins

# Operation

Each of plugins operates asynchronously and can send update of its status at any time. On each plugin update state is sent upstream (including previous cached plugin state) so it is possible to have each plugin update separetely.

Each plugin have a name (by default plugin name) and instance (in case of plugins that will be used multiple times, like network interfaces).


# Plugins

## CPU

## Clock

Parameters:

* `long_format`/`short_format` - time format, as golang formatting string
* `interval` - update interval in miliseconds. Note that it can be longer than it because of run time so for update every second something like 990 ms is better

## Network

Instance is network interface to monitor

## I3blocks

i3blocks-compatible input

Parameters:

* `command` - command to run
* `prefix` - prefix command with text
