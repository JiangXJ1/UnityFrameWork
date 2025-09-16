REM 声明采用UTF-8编码
chcp 65001
svn add --force Notice/TEST
svn commit Notice/TEST -m "[LH]备份SDK后台配置"
aws s3 cp Notice/TEST s3://ipretobucket/orig-roa/message/TEST/ --recursive
pause