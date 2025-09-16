REM 声明采用UTF-8编码
chcp 65001
svn add --force Notice/LHA
svn commit Notice/LHA -m "[LHA]备份SDK后台配置"
aws s3 cp Notice/LHA s3://ipretobucket/orig-roa/message/LHA/ --recursive
pause