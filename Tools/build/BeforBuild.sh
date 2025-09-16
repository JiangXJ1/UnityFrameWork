#!/bin/bash
echo "+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ BeforBuild.sh"

#异常后停止运行
set -e
source ./BuildDef.sh
source ./client_auto_upload.sh

#从模板文件读取参数
ReadParamsFromTemplate
#设置工程路径
SetBuildPath

cd ${RootDir}/Tools/build/SDKAssetTools/

echo "SDK:{$SDK}"
#拷贝匿名USDK资源
if [ "$SDK" == "1" ]
then
	python3 copyUSDKSimple.py
fi

cd ${RootDir}/Tools/build/

MSG_TITLE="构建项目:${PROJECT_NAME}(${ROA_G_Platform},${ROA_G_Version},分支:${ROA_G_Branch},svn:${ROA_G_SVNRevision})"
#仅上传资源
if [ ${JUST_UPLOAD_PATCH} == "1" ]
then
	echo "上传补丁"
	MsgPush ${MSG_TITLE} "仅上传补丁！"
	DoUploadPatch
	exit
elif [ ${JUST_UPLOAD_ALL} == "1" ]
then
	echo "上传全部资源"
	MsgPush ${MSG_TITLE} "仅上传全量资源！"
	DoUploadAll
	exit
fi


#判断参数是否存在
if [ ${ROA_G_Version} == "0" ]
then
	if [[ ${PATCH} == "0" && "${ROA_G_Branch}" == *branches* ]]
	then
		echo "当前是Release分支，出整包必须填版本号！"
		exit 1
	fi

	#没有版本号参数
	echo "没有VERSION参数"
	#如果不存在版本标记文件
	if [ ! -s ${CLIENT_VER_RECORD_PATH} ]
	then
		echo "没有版本号缓存文件:"${CLIENT_VER_RECORD_PATH}"使用0.0.1"
		echo "0.0.1">${CLIENT_VER_RECORD_PATH}
	fi

	LastVer=$(cat ${CLIENT_VER_RECORD_PATH})
	LastVerArr=(${LastVer//./ })
	echo "最近一次打包版本号:"${LastVer}
	i=1
	for var in ${LastVerArr[@]}
	do
		if [ $i == 1 ]
		then
			LastVerA=$((var))
		elif [ $i == 2 ]
		then
			LastVerB=$((var))
		elif [ $i == 3 ]
		then
			LastVerC=$((var))
		fi
		i=$((i+1))
	done
	if [ ${VERSION_ADD_TYPE} == "0" ]
	then
		NewVer=${LastVer}
	elif [ ${VERSION_ADD_TYPE} == "a" ]
	then
		NewVer=$((LastVerA+1))".0.0"
	elif [ ${VERSION_ADD_TYPE} == "b" ]
	then
		NewVer=${LastVerA}"."$((LastVerB+1))".0"
	elif [ ${VERSION_ADD_TYPE} == "c" ]
	then
		NewVer=${LastVerA}"."${LastVerB}"."$((LastVerC+1))
	else
		NewVer=${LastVer}
	fi
	ROA_G_Version=${NewVer}
else
	#存在版本号参数则强制赋值
	NewVer=${ROA_G_Version}
fi
echo "当前正在打包版本:"${ROA_G_Version}
echo ${ROA_G_Version}>${CLIENT_VER_RECORD_PATH}

#判断参数是否存在 BuildCode是0则自增长
if [ ${ROA_G_BundleVersionCode} == "0" ]
then
	#没有构建号参数
	echo "参数BUILDCODE为0,需要自增长"
	#如果不存在版本标记文件
	if [ ! -s ${CLIENT_BUILDCODE_RECORD_PATH} ]
	then
		echo "0">${CLIENT_BUILDCODE_RECORD_PATH}
	fi

	LastBuildCode=$(cat ${CLIENT_BUILDCODE_RECORD_PATH})
	ROA_G_BundleVersionCode=$((LastBuildCode+1))
fi

echo "当前打包构建号"${ROA_G_BundleVersionCode}

echo ${ROA_G_BundleVersionCode}>${CLIENT_BUILDCODE_RECORD_PATH}

cd ${RootDir}/Tools/build/


if [ ${PATCH} == "0" ]
then
	echo "制作整包"
	sh Build.sh

else
	if [ ${ROA_G_Version} == "0" ]
	then
        echo "制作热更补丁必须填写版本号！(参数名：version)"
        exit 1
    fi

	echo "制作热更补丁"
	sh Build_Hotfix.sh
fi