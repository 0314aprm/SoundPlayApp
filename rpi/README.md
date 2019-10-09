

* service
* player


===========================================================
service

# services
一個でいい？


# chars
曲データ
  フォルダ数
  フォルダ毎ファイル数
プレイヤー
  status:  playing? mp3 state
  current file: readCurrentFileNumber
  volume (1 octet)
  playback mode (1 octet)
  EQ


# 操作
pause
start
next
previous
playFile
increase volume
decrease volume
set volume
set loop
enable loop
disable loop

#
error response

===========================================================
player

# DFPlayer 操作
UART で指示
音楽ファイルリスト？


===========================================================

# docker
https://hub.docker.com/_/debian/


# サービス化, デーモン化
https://www.write-ahead-log.net/entry/2017/07/18/230634
systemd
  unit
他のライブラリ


# API
BLE
Nordic UART
  UART/Serial Port Emulation over BLE




