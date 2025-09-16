REM 声明采用UTF-8编码
chcp 65001
svn add --force Notice/LHB
svn commit Notice/LHB -m "[LHB]备份SDK后台配置"
aws s3 cp Notice/LHB s3://ipretobucket/orig-roa/message/LHB/ --recursive
pause