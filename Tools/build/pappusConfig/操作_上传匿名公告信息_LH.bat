REM 声明采用UTF-8编码
chcp 65001
svn add --force Notice/LH
svn commit Notice/LH -m "[LH]备份SDK后台配置"
aws s3 cp Notice/LH s3://ipretobucket/orig-roa/message/LH/ --recursive
pause