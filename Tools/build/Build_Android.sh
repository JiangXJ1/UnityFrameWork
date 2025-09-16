#!/bin/bash
set -e

InitAndroid()
{
    # NativeProjectPath=${CLIENT_BUILD}/One_Android
	# NativeProjectPathCopy=${CLIENT_BUILD}/One_Android_Copy
    TmpProjectPath=${TMP_ROOT}/${PROJECT_NAME}_AOS_${ROA_G_Branch}_${ROA_G_SVNRevision}_$(date "+%m%d")_V${ROA_G_Version}

    AppFileName=${PROJECT_SHORT_NAME}_$(date "+%m%d")_V${ROA_G_Version}_${ROA_G_BundleVersionCode}_${ROA_G_SVNRevision}
	echo ${AppFileName}
}
export -f InitAndroid

ExecuteBuildAndroid()
{
    CMD_BUILD_SWITCH_PLATFORM="${CMD_BUILD_SWITCH_PLATFORM} -aos"
    CMD_BUILD_RES="${CMD_BUILD_RES} -aos"
    CMD_BUILD_APP="${CMD_BUILD_APP} -aos"
    #切换平台
	echo ${CMD_BUILD_SWITCH_PLATFORM}
    eval ${CMD_BUILD_SWITCH_PLATFORM}
	if [[ $? != 0 ]]; then
        MsgPush ${MSG_TITLE} "切换平台失败"
        exit
    fi
	echo ">>>>>>>>等待30秒(防止下次启动unity完成后卡死无法关闭)<<<<<<<<"
	sleep 30

    #先刷新资源
    # MsgPush ${MSG_TITLE} "刷新资源（AAS、Shader搜集变体...）"
    # echo ${CMD_BUILD_RES}
    # eval ${CMD_BUILD_RES}
    # if [[ $? != 0 ]]; then
    #     MsgPush ${MSG_TITLE} "刷新资源资源失败"
    #     exit
    # fi
	# echo ">>>>>>>>等待30秒(防止下次启动unity完成后卡死无法关闭)<<<<<<<<"
	# sleep 30
    #再打包APP
    MsgPush ${MSG_TITLE} "构建Unity"
    echo ${CMD_BUILD_APP}
    eval ${CMD_BUILD_APP}
    if [[ $? != 0 ]]; then
        MsgPush ${MSG_TITLE} "构建Unity失败"
        exit
    fi

	MsgPush ${MSG_TITLE} "构建Unity完成！"
	MsgPush ${MSG_TITLE} "正在检查是否上传资源！"

    #if [ "${SDK}" == "1" || [ "${SDK}" == "2" ]
	#then
	#    echo "上传符号表"
	#	sh gradlew uploadCrashlyticsSymbolFileRelease
	#fi
}
export -f ExecuteBuildAndroid

OnBuildFinishAndroid()
{
    if [ "$SDK" == "1" ]
	then
		STR='市调USDK'
	elif [ "$SDK" == "2" ]
	then
		STR='PappusSDK'
	elif [ "$SDK" == "3" ]
	then
		STR='正式USDK'
	else
		STR='内网'
	fi
	TEMP_NAME=${ROA_G_Branch}_${ROA_G_Version}_${ROA_G_SVNRevision}_$(date "+%Y.%m%d.%H%M")
	DestinationFolder=/Volumes/cdc01/真机测试包/$STR/$TEMP_NAME
	mkdir -p $DestinationFolder
	cp -n -f ${TmpProjectPath}/* $DestinationFolder
	echo "拷贝完成"

	original_string=$STR'\\NAME'
	modified_string=$(echo "$original_string" | sed "s/NAME/$TEMP_NAME/")
	echo $modified_string

	MsgPush '拷贝' '新包地址:\\\\10.23.0.3\\cdc01\\真机测试包\\'$modified_string
}