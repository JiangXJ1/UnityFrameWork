REM 声明采用UTF-8编码
chcp 65001
svn add --force Notice/DEV
svn commit Notice/DEV -m "[LH]备份SDK后台配置"
aws s3 cp Notice/DEV s3://ipretobucket/orig-roa/message/DEV/ --recursive
pause