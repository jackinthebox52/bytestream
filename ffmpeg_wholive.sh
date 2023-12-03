ffmpeg -user_agent 'Mozilla/5.0 (Windows NT 10.0; rv:108.0) Gecko/20100101 Firefox/108.0' \
    -headers 'Referer: https://www.niaomea.me/\r\nOrigin: https://www.niaomea.me\r\n' \
    -i "https://ed-c003.edgking.me/plyvivo/cipotimu1aluximet0xi/chunklist.m3u8"\
    -codec: copy -start_number 0 -hls_time 10 -hls_list_size 0 -f hls ./hls/buffstream.m3u8\
    -reconnect 1 -reconnect_at_eof 1 -reconnect_streamed 1
