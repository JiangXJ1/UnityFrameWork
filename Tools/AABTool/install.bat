echo build local-testing apks
echo off
del output.apks
java -jar bundletool-all-1.8.1.jar build-apks --bundle=%1 --output=output.apks --local-testing
echo on
echo install apks
echo off
java -jar bundletool-all-1.8.1.jar install-apks --apks=output.apks
echo on
echo install apks finished
pause