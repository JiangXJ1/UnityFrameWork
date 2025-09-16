#!/bin/bash
set -e

InitIOS()
{
    # NativeProjectPath=${CLIENT_BUILD}/One_iOS
    TmpProjectPath=${TMP_ROOT}/TMP_${PROJECT_NAME}_iOS
    AppFileName=${PROJECT_SHORT_NAME}_$(date "+%m%d")_V${ROA_G_Version}_${ROA_G_BundleVersionCode}_${ROA_G_SVNRevision}.ipa
	echo ${AppFileName}
    echo "ios"
	#清空临时文件夹
	rm -rf ${TmpProjectPath}
	echo ${TmpProjectPath}
}
export -f InitIOS

ExecuteBuildIOS()
{
    #发布前根据SDK删除不需要的SDK目录
	# if [ "${SDK}" == "1" ]
	# then
	# 	#TODO 删除匿名SDK的各目录
	# 	rm -rf ${CLIENT}/Assets/PappusSDK
	# elif [ "${SDK}" == "2" ]
	# then
	# 	#TODO 删除IGG SDK 各目录(保留AF和Firebase以及Firebase需要的库)
	# 	rm -rf ${CLIENT}/Assets/USDK/ADSDK/Facebook-Vendor-Impl
	# 	rm -rf ${CLIENT}/Assets/USDK/FacebookSDK
	# 	rm -rf ${CLIENT}/Assets/USDK/GPC
	# 	rm -rf ${CLIENT}/Assets/USDK/Resources
	# 	rm -rf ${CLIENT}/Assets/USDK/Script
    #     #使用Pappus的google-service.json和GoogleService-Info.plist替换正式版的
	# 	cp -rf ${CLIENT}/Assets/google-services.json.pappus ${CLIENT}/Assets/google-services.json
	# 	cp -rf ${CLIENT}/Assets/GoogleService-Info.plist.pappus ${CLIENT}/Assets/GoogleService-Info.plist
	# else
	# 	#TODO 删除IGG SDK和匿名SDK的各目录
	# 	rm -rf ${CLIENT}/Assets/USDK
	# 	rm -rf ${CLIENT}/Assets/PappusSDK
	# fi

	# #不同版本证书不同,匿名SDK证书也不同
	if [ "${SDK}" == "0" ]
	then
		ExportPotionsPlist=${RootDir}/Tools/build/ios_export_options/dev.plist
	elif [ "${SDK}" == "1" ]
	then
		ExportPotionsPlist=${RootDir}/Tools/build/ios_export_options/dev_gpc.plist
	elif [ "${SDK}" == "2" ]
	then
		ExportPotionsPlist=${RootDir}/Tools/build/ios_export_options/dev_pappus.plist
	else
		ExportPotionsPlist=${RootDir}/Tools/build/ios_export_options/dev.plist
	fi

    CMD_BUILD_SWITCH_PLATFORM="${CMD_BUILD_SWITCH_PLATFORM} -ios"
    CMD_BUILD_RES="${CMD_BUILD_RES} -ios"
    CMD_BUILD_APP="${CMD_BUILD_APP} -ios"
    #切换平台
	echo ${CMD_BUILD_SWITCH_PLATFORM}
    eval ${CMD_BUILD_SWITCH_PLATFORM}
    #先打包资源
    MsgPush ${MSG_TITLE} "构建资源"
    echo ${CMD_BUILD_RES}
    eval ${CMD_BUILD_RES}
    if [[ $? != 0 ]]; then
        MsgPush ${MSG_TITLE} "构建资源失败"
        exit
    fi
    #再打包APP
    MsgPush ${MSG_TITLE} "构建Unity"
    echo ${CMD_BUILD_APP}
    eval ${CMD_BUILD_APP}
    if [[ $? != 0 ]]; then
        MsgPush ${MSG_TITLE} "构建Unity失败"
        exit
    fi
	# #Pappus SDK的Scheme名不同
	if [ "${SDK}" == "0" ]
	then
		IOS_SCHEME=ProjectROA
	elif [ "${SDK}" == "1" ]
	then
		IOS_SCHEME=ProjectROA
	elif [ "${SDK}" == "2" ]
	then
		IOS_SCHEME=ProjectROA
	else
		IOS_SCHEME=ProjectROA
	fi

    # IOS_CONFIG=Release
	# IOS_WORKSPACE=Unity-iPhone.xcworkspace
	# IOS_XARCHIVE=${ArchiveRoot}/$(date "+%Y-%m-%d")/${PROJECT_NAME}_${ROA_G_Version}_$(date "+%m%d")_${ROA_G_SVNRevision}.xcarchive

    # cd ${TmpProjectPath}

	# if [ "${SDK}" == "1" ] || [ "${SDK}" == "2" ]
	# then
	# 	#清理
	# 	xcodebuild clean -workspace ${IOS_WORKSPACE} -scheme ${IOS_SCHEME} -configuration ${IOS_CONFIG}
	# 	#构建
	# 	xcodebuild archive -archivePath ${IOS_XARCHIVE} -workspace ${IOS_WORKSPACE} -sdk iphoneos -scheme ${IOS_SCHEME} -configuration Release >> ${TmpProjectPath}/archive.log
	# 	#打包
	# 	xcodebuild -exportArchive -archivePath ${IOS_XARCHIVE} -exportPath ${TmpProjectPath}/Build -exportOptionsPlist ${ExportPotionsPlist} >> ${TmpProjectPath}/export.log
	# 	#上传符号表
	# 	./Pods/FirebaseCrashlytics/upload-symbols -gsp ./GoogleService-Info.plist -p ios "${IOS_XARCHIVE}/dSYMs"
	# else
	# 	#清理
	# 	xcodebuild clean -scheme ${IOS_SCHEME} -configuration ${IOS_CONFIG}
	# 	#构建
	# 	xcodebuild archive -archivePath ${IOS_XARCHIVE} -sdk iphoneos -scheme ${IOS_SCHEME} -configuration Release >> ${TmpProjectPath}/archive.log
	# 	#打包
	# 	xcodebuild -exportArchive -archivePath ${IOS_XARCHIVE} -exportPath ${TmpProjectPath}/Build -exportOptionsPlist ${ExportPotionsPlist} >> ${TmpProjectPath}/export.log
	# fi

    # if [[ $? != 0 ]]; then
    #     MsgPush ${MSG_TITLE} "Xcode项目编译失败"
    #     exit
    # fi

    # #创建导出目录
    # mkdir -p ${TMP_ROOT}/IPA
    # cp ${TmpProjectPath}/Build/${PROJECT_NAME}.ipa ${TMP_ROOT}/IPA/${AppFileName}
    # MsgPush ${MSG_TITLE} "发包完毕 ${AppFileName}"

	# if [ ${SHARE} == 1 ]
	# then
	# 	cp -rvf ${TMP_ROOT}/IPA/${AppFileName} /Volumes/cdc01/03.QA/04.临时文件/Debug/IPA
	# fi
}
export -f ExecuteBuildIOS