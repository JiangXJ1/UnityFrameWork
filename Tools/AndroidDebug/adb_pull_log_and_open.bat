cd .
adb pull /sdcard/Android/data/com.tribalrise.android/files/Log ./

for /f "delims=" %%x in ('dir /b /od "Log\*"') do set recent=%%x
start "" "Log\%recent%"

#pause