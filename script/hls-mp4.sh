#!/bin/sh

ffmpeg -i ./hls/IASDF.m3u8 -c copy -bsf:a aac_adtstoasc ./streams/video/IASDF.mp4  