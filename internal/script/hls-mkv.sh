#!/bin/sh
#Format of the command: ./hls-mkv.sh UUID Format  or './hls-mkv.sh FAIOEF mp4'
ffmpeg -i ./streams/hls/$1/chunklist.m3u8 -c copy -bsf:a aac_adtstoasc ./streams/video/$1.mkv