#!/bin/bash
echo "+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ Build_Hotfix.sh"
set -e

if [[ "${ROA_G_Branch}" == *branches* ]]
then
    echo "当前是正式分支"
fi

MSG_TITLE="构建热更:${PROJECT_NAME}(${ROA_G_Platform},${ROA_G_Version},分支:${ROA_G_Branch},svn:${ROA_G_SVNRevision})"

if [ ! ${ROA_G_Version} ]
then
	MsgPush ${MSG_TITLE} "补丁参数错误，当前值:${ROA_G_Version}"
	exit
fi

MsgPush ${MSG_TITLE} "开始构建"

#检测Unity是否开启,如果已经开启,直接杀掉
CloseUnity

CMD_BUILD_SWITCH_PLATFORM="${UNITY} -batchmode -quit -projectPath ${CLIENT} -logFile ${SWITCH_LOG_PATH} -executeMethod BuildTools.SwitchPlatform"
CMD_BUILD_HOTFIX="${UNITY} -batchmode -quit -projectPath ${CLIENT} -logFile ${RES_LOG_PATH} -executeMethod BuildTools.BuildHotfix -patch ${ROA_G_Version} -generateDll ${ROA_G_HOTDll} -branch ${ROA_G_Branch}"

if [ "${ARMV7}" == "1" ]
then
	SYMBOLS="${SYMBOLS};ARMV7"
fi

CMD_BUILD_SWITCH_PLATFORM="${CMD_BUILD_SWITCH_PLATFORM} -symbols \"${SYMBOLS}\""

if [ "$ROA_G_Platform" == "ios" ]
then
	CMD_BUILD_SWITCH_PLATFORM="${CMD_BUILD_SWITCH_PLATFORM} -ios"
	CMD_BUILD_HOTFIX="${CMD_BUILD_HOTFIX} -ios"
elif [ "$ROA_G_Platform" == "android" ]
then
	CMD_BUILD_SWITCH_PLATFORM="${CMD_BUILD_SWITCH_PLATFORM} -aos"
	CMD_BUILD_HOTFIX="${CMD_BUILD_HOTFIX} -aos"
fi

echo 工程路径:${CLIENT}
echo 日志路径:${SWITCH_LOG_PATH}

#切换平台
echo ${CMD_BUILD_SWITCH_PLATFORM}
eval ${CMD_BUILD_SWITCH_PLATFORM}
#制作热更新
MsgPush ${MSG_TITLE} "构建热更资源"
echo ${CMD_BUILD_HOTFIX}
eval ${CMD_BUILD_HOTFIX}
MsgPush ${MSG_TITLE} "构建完毕"


#master分支才自动上传补丁资源
echo "准备上传资源"
MsgPush ${MSG_TITLE} "准备上传补丁"

#上传补丁文件
DoUploadPatch ${MSG_TITLE}

MsgPush ${MSG_TITLE} "出补丁完成！"