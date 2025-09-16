
#!/bin/bash

#UNITY程序的路径
export UNITY=/Applications/2022.3.33f1/Unity.app/Contents/MacOS/Unity

#临时构建目录
export TMP_ROOT=/Users/cdc01/Documents/export

#编译日志目录
export LOG_ROOT=$TMP_ROOT/BuildLog

#日志路径
export SWITCH_LOG_PATH=${LOG_ROOT}/switchlog_$(date "+%Y%m%d%H%M%S").log
export LOG_PATH=${LOG_ROOT}/log_$(date "+%Y%m%d%H%M%S").log
export RES_LOG_PATH=${LOG_ROOT}/reslog_$(date "+%Y%m%d%H%M%S").log

#XCode编译后xcarchive文档根目录
export ArchiveRoot=/Users/cdc01/Library/Developer/Xcode/Archives

#android aab 分包install_pack
export Install_Pack=UnityDataAssetPack
#android aab 分包fast_pack_1
export Fast_Pack_1=CustomFastFollo
#android aab 分包fast_pack_2
export Fast_Pack_2=SecondFastFollow
#android aab 分包fast_pack_3
export Fast_Pack_3=ThirdFastFollow

#公共路径
export CLIENT_VER_RECORD_PATH=/Users/cdc01/Documents/export/ver.txt
export CLIENT_BUILDCODE_RECORD_PATH=/Users/cdc01/Documents/export/buildcode.txt

export RootDir=""
export PathchDir=""
export CLIENT=/Users/cdc01/Documents/IGG/UnityProject
export PathchDir=/Users/cdc01/Documents/export/RemoteResource
SetBuildPath()
{
    echo "-----------------------------------------------------------设置工程目录"
    #默认项目根目录
    if [[ "${ROA_G_Branch}" == *branches* ]]
	then
        echo "build:使用线上工程"
        RootDir=/Users/cdc01/Documents/IGG_Online
	elif [[ "${ROA_G_Branch}" == *tags* ]]
    then
		echo "build:使用历史版本工程"
        RootDir=/Users/cdc01/Documents/IGG_History
        echo "build:工程目录${RootDir}"
    elif [[ "${ROA_G_Branch}" == *trunkqa* ]]
    then
		echo "build:使用QA工程"
        RootDir=/Users/cdc01/Documents/IGG_QA
        echo "build:工程目录${RootDir}"
    else
		echo "build:使用默认工程"
        RootDir=/Users/cdc01/Documents/IGG
        echo "build:工程目录${RootDir}"
	fi
    PathchDir=/Users/cdc01/Documents/export/RemoteResource
    echo "build:资源生成目录${PathchDir}"
    #Unity工程目录
    CLIENT=$RootDir/UnityProject
    echo "-----------------------------------------------------------"
}
export -f SetBuildPath


GitLog()
{
    last_log=`git log --pretty=oneline -5`
    binfo=`git symbolic-ref HEAD`

    MsgPush AutoScp "$1代码拉取完成"
    MsgPush AutoScp $binfo

    echo "$binfo \n"

    MsgPush AutoScp $last_log

    echo "$last_log \n"
}
export -f GitLog


#解析参数模板
ReadParamsFromTemplate()
{
    # 目前文件参数
    # echo "PROJECT_NAME参数:"${PROJECT_NAME}
    # echo "VERSION_ADD_TYPE参数:"${VERSION_ADD_TYPE}
    # echo "SDK参数:"${SDK}
    # echo "ARMV7参数:"${ARMV7}
    # echo "SHARE参数:"${SHARE}
    # echo "USE_UWA参数:"${USE_UWA}
    # echo "PATCH参数:"${PATCH}
    # echo "GPC_TEST参数:"${GPC_TEST}
    # echo "JUST_UPLOAD_PATCH参数:"${JUST_UPLOAD_PATCH}
    # echo "JUST_UPLOAD_ALL参数:"${JUST_UPLOAD_ALL}
    echo "------------------------------------模板参数"
    TplName=${ROA_G_TplFileOption}".tpl"
	echo "解析参数模板文件:"${TplName}

    if [ -e  ${TplName} ]
    then
        for line in `cat ${TplName}`
        do
            line=$(echo ${line} | tr '\n\r' '=')
            line=$(echo ${line} | tr '\n' '=')
            array=(${line//=/ })
            eval export ${array[0]}=${array[1]}
            echo ${array[0]}=${array[1]}
        done
    else
        echo "找不到模板参数文件:"${TplName}
        exit 1
    fi
    echo "------------------------------------"
}
export -f ReadParamsFromTemplate


MsgPush()
{
    echo "$1 $2"
    # curl 'http://im-api.skyunion.net/msg'   \
    #     -H 'Content-Type: application/json'     \
    #     -d "{\"token\": \"fcf0f94a12059dfad8ec807ebe924c16\", \"target\": \"group\",\"room\":\"10035669\",\"title\":\"$1\",\"content\": \"$2\",\"content_type\":\"1\"}"
}
export -f MsgPush


CloseUnity()
{
    # 方法1:  如果当前unity有多个进程会有问题 $UnityPID 会把所有pid连起来无法遍历
    UnityPID=`ps -ef | grep '/Applications/2022.3.33f1/Unity.app/Contents/MacOS/Unity' | grep -v grep | awk '{printf $2}'`
    echo "UnityPID   $UnityPID"
    for id in $UnityPID
    do
        echo "UnityPID--  $id"
        # kill -9 $id
        killall -9 Unity
        echo "找到打开的Unity3D，已强制关闭！"
    done
}
export -f CloseUnity