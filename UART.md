

https://qiita.com/tetsu_koba/items/f8afbb8326ee42fd27f5

https://github.com/tarm/serial
https://github.com/goburrow/serial


=======================================================================================--
Raspberry Pi (Zero W)
https://physical-computing-lab.net/raspberry-pi-b/1-2.html

two built-in UARTs
https://www.raspberrypi.org/documentation/configuration/uart.md
* wireless/Bluetooth module がない場合は PL011 が primary
- PL011:  connected to the Bluetooth module
  /dev/ttyAMA0
- mini UART:  used as the primary UART
  /dev/ttyS0

symlinks
  primary UART:  /dev/serial0
  seconary UART:  /dev/serial1

デフォルトでは Linux's use of console UART を無効化する必要がある
  https://www.raspberrypi.org/documentation/configuration/uart.md


# GPIO
transmit and receive pins (GPIO)
  GPIO 14 and GPIO 15 respectively, which are pins 8 and 10 on the GPIO header
interface
  /sys/class/gpio/


=======================================================================================--
DFPlayer Mini

lib: https://github.com/DFRobot/DFRobotDFPlayerMini

working voltage: DC3.2~5.0V; Type :DC4.2V

UART Port
  Standard Serial; TTL Level
  baud rate: 9600 (default)

# command format
serial communication


# storage
file: 1~255
folder: 1~99


# module
1. power on
2. initialization


# device
U-disk
TF card

push-in, pull-out



# command
tracking: 0-2999

# query
status
volume:  0-30
EQ:  Normal/Pop/Rock/Jazz/Classic/Base
playback mode:  (0/1/2/3) Repeat/folder, repeat/single, repeat/random 
total number of files (TF card, U-disk, flash)
current track (TF card, U-disk, flash)


