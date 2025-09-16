echo install apks
echo off
java -jar bundletool-all-1.8.1.jar install-apks --apks=%1
echo on
echo install apks finished
pause