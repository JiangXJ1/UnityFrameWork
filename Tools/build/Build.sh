#!/bin/bash
echo "+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ Build.sh"
source ./Build_Android.sh
source ./Build_IOS.sh
set -e

echo "Jenkins全局参数"
echo "工作目录:"${WORKSPACE}
echo  "版本号:"$ROA_G_Version
echo  "分支名:"$ROA_G_Branch
echo  "构建号:"$ROA_G_BundleVersionCode
echo  "模板名:"$ROA_G_TplFileOption
echo  "平台:"$ROA_G_Platform
echo  "渠道:"$ROA_G_ChannelName
echo  "SVN Revision:"$ROA_G_SVNRevision
echo  "Debug:"$ROA_G_Debug

export TmpProjectPath=""
export AppFileName=""

if [[ "${ROA_G_Branch}" == *branches* ]]
then
    echo "当前是正式包"
fi

if [ $ROA_G_Debug == "true" ]
then
    MSG_TITLE="构建项目:${PROJECT_NAME}(${ROA_G_Platform},${ROA_G_Version},分支:${ROA_G_Branch},svn:${ROA_G_SVNRevision},debug)"
else
    MSG_TITLE="构建项目:${PROJECT_NAME}(${ROA_G_Platform},${ROA_G_Version},分支:${ROA_G_Branch},svn:${ROA_G_SVNRevision},release)"
fi

#缩写
if [ ${PROJECT_NAME} == "LastHome" ]
then
	PROJECT_SHORT_NAME="LH"
fi
#检测版本号
if [ ! ${ROA_G_Version} ]
then
    MsgPush ${MSG_TITLE} "版本参数错误，必须为1.0.0格式，当前值:${ROA_G_Version}"
    exit
fi
#构建号
if [ ! ${ROA_G_BundleVersionCode} ]
then
    MsgPush ${MSG_TITLE} "构建号错误，当前值:${ROA_G_BundleVersionCode}"
    exit
fi
#初始化工程信息
if [ "${ROA_G_Platform}" == "ios" ]
then
	InitIOS
elif [ "${ROA_G_Platform}" == "android" ]
then
	InitAndroid
else
    MsgPush  ${MSG_TITLE} "平台参数错误,必须为ios/android,当前值:${ROA_G_Platform}"
    exit
fi
#通知开始
MsgPush ${MSG_TITLE} "开始自动构建"

#检测Unity是否开启,如果已经开启,直接杀掉
CloseUnity
# unity命令打包
# -quit 			在其他命令执行完毕后退出 Unity Editor。这可能导致错误消息被隐藏（但是，它们仍会出现在 Editor.log 文件中）。
# -batchmode		以批处理模式运行 Unity。在批处理模式下，Unity 无需人工交互即可运行命令行参数。它还会禁止需要人工交互的弹出窗口（例如 Save Scene 窗口）；但是，Unity 编辑器本身会像往常一样打开
# -executeMethod	Unity 打开项目后以及可选的资源服务器更新完成后，立即执行静态方法。
# -projectPath		在指定路径下打开项目。如果路径名包含空格，请将其用引号引起来。

#切换平台
CMD_BUILD_SWITCH_PLATFORM="${UNITY} \
-projectPath ${CLIENT} \
-logFile ${SWITCH_LOG_PATH} \
-executeMethod BuildTools.SwitchPlatform \
-batchmode \
-quit"

#构建资源
CMD_BUILD_RES="${UNITY} \
-projectPath ${CLIENT} \
-logFile ${RES_LOG_PATH} \
-executeMethod BuildTools.BeforeBuildResource \
-ver ${ROA_G_Version} \
-useAAB ${ROA_G_AAB} \
-updatePack ${ROA_UPDATE_PACK} \
-branch ${ROA_G_Branch} \
-batchmode \
-quit"

#构建APP
CMD_BUILD_APP="${UNITY} \
-projectPath ${CLIENT} \
-logFile ${LOG_PATH} \
-executeMethod BuildTools.CommandBuild \
-export ${TmpProjectPath}/${AppFileName} \
-updatePack ${ROA_UPDATE_PACK} \
-projectname ${PROJECT_NAME} \
-fileName ${AppFileName} -ver ${ROA_G_Version} \
-buildcode ${ROA_G_BundleVersionCode} \
-useAAB ${ROA_G_AAB} \
-channelName ${ROA_G_ChannelName} \
-ip ${ROA_G_IP} \
-batchmode \
-quit"

#检测是否是DEBUG
if [ $ROA_G_Debug == "true" ]
then
    CMD_BUILD_SWITCH_PLATFORM="${CMD_BUILD_SWITCH_PLATFORM} -debug"
fi

#是否接入UWA
if [ ${USE_UWA} == 1 ]
then
    CMD_BUILD_SWITCH_PLATFORM="${CMD_BUILD_SWITCH_PLATFORM} -uwa"
fi

#是否接入UWA
if [ ${GPC_TEST} == 1 ]
then
    CMD_BUILD_SWITCH_PLATFORM="${CMD_BUILD_SWITCH_PLATFORM} -gpctest"
fi

if [ "${ARMV7}" == "1" ]
then
	SYMBOLS="${SYMBOLS};ARMV7"
fi

#检测是否是AB测试
if [ ${ABTEST} == 1 ]
then
    SYMBOLS="${SYMBOLS};ABTEST"
fi

#切换平台
CMD_BUILD_SWITCH_PLATFORM="${CMD_BUILD_SWITCH_PLATFORM} -symbols \"${SYMBOLS}\""

#执行打包
if [ "${ROA_G_Platform}" == "ios" ]
then
	ExecuteBuildIOS
    echo "ios 暂时忽略"
elif [ "${ROA_G_Platform}" == "android" ]
then
    echo "执行打包Android"
    ExecuteBuildAndroid
fi


# echo "准备上传资源"
# MsgPush ${MSG_TITLE} "准备上传资源！"

# #上传资源
# DoUploadAll ${MSG_TITLE}

# MsgPush ${MSG_TITLE} "开始自动拷贝包到共享网盘！"
# #自动拷贝包
# if [ "$ROA_G_Platform" == "ios" ]
# then
# 	echo "ios"
# elif [ "${ROA_G_Platform}" == "android" ]
# then
# 	OnBuildFinishAndroid
# fi
# MsgPush ${MSG_TITLE} "出包完成！"
