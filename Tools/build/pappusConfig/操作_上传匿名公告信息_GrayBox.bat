REM 声明采用UTF-8编码
chcp 65001
svn add --force Notice/GrayBox
svn commit Notice/GrayBox -m "[LH]Backup Config"
aws s3 cp Notice/GrayBox s3://ipretobucket/orig-roa/message/GrayBox/ --recursive
pause