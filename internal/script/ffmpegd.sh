#!/bin/sh
#Format of the command: ./ffmpegd.sh StreamURL UUID Referrer
origin="${3%/}"
ffmpeg -user_agent 'Mozilla/5.0 (Windows NT 10.0; rv:108.0) Gecko/20100101 Firefox/108.0' \
    -headers "Referer: $3\r\nOrigin: $origin\r\n" \
    -reconnect 1 -reconnect_at_eof 1 -reconnect_streamed 1 -reconnect_delay_max 10 -nostdin \
    -i "$1"\
    -codec: copy -start_number 0 -hls_time 10 -hls_list_size 0 -f hls ./streams/hls/$2/chunklist.m3u8 \
    & echo $! > ./streams/pids/$2.pid
