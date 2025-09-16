REM 声明采用UTF-8编码
chcp 65001
svn add --force Notice/DH
svn commit Notice/DH -m "[DH] backup"
aws s3 cp Notice/DH s3://ipretobucket/orig-roa/message/DH/ --recursive
pause