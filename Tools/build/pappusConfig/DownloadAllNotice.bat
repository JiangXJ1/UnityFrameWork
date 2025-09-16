@echo off
REM 声明采用UTF-8编码
chcp 65001
REM bitsadmin /transfer "download" /download /priority normal "http://10.23.27.145/message"  %cd%/notice/
REM s3cmd get s3://ipretobucket/orig-roa/message/ ./Notice —recursive
aws s3 sync "s3://ipretobucket/orig-roa/message" ./Notice 
echo 下载完成 Export Success



