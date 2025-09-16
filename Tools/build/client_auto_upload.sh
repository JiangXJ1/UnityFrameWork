#!/bin/bash

#上传补丁
DoUploadPatch()
{
    MSG_TITLE=$1
    echo "DoUploadPatch():"$ROA_G_Platform ${ROA_G_Branch} #1:平台  2：分支名称
    #master分支才自动上传补丁资源
    echo "准备上传资源"
    if [ "${ROA_G_Platform}" == "ios" ]
    then
        CurPlaytform=iOS
    elif [ "$ROA_G_Platform" == "android" ]
    then
        CurPlaytform=Android
    else
        echo "平台参数错误！"
        exit 1
    fi

    UploadAssetsRoot=${PathchDir}/${CurPlaytform}/${ROA_G_Branch}
    playerVersionFile=${UploadAssetsRoot}/CurPlayerVersion.txt
    playerVersion=`sed -n '1p' ${playerVersionFile}`
    TargetPath=${UploadAssetsRoot}/pack_${CurPlaytform}_${playerVersion}/UploadAssetsAfterPatch/ServerData

    echo "playerVersion:"${playerVersion}
    echo "playerVersionFile:"${playerVersionFile}
    echo "待上传路径:"${TargetPath}

    if [[ "${ROA_G_Branch}" == *branches* ]]
    then
        echo "开始上传资源到正式CDN"
        MsgPush ${MSG_TITLE} "开始上传资源到正式CDN"
        aws s3 cp ${TargetPath} s3://ipretobucket/orig-roa/ --recursive
        echo "上传资源完成！"
    elif [ "${ROA_G_Branch}" == "develop" ]
    then
        echo "develop分支复制补丁资源到本地服务器，安卓机访问地址：http://10.23.26.34/ 苹果机访问地址：http://10.23.26.35/"
        echo "启动服务器命令：sudo apachectl start"
        echo "关闭服务器命令：sudo apachectl stop"
        # 打开目录写入权限命令：sudo chmod -R o+w /Library/WebServer/Documents
        cp -rvf ${TargetPath}/* /Library/WebServer/Documents/
    else
        echo "当前分支不会执行上传操作！"
    fi
    MsgPush ${MSG_TITLE} "资源上传完成！"
}
export -f DoUploadPatch



#上传全部资源
DoUploadAll()
{
    MSG_TITLE=$1
    echo "DoUploadAll():"${ROA_G_Platform} ${ROA_G_Branch} #1:平台  2：分支名称 3：工程目录别名
    echo "准备上传资源"
    if [ "${ROA_G_Platform}" == "ios" ]
    then
        CurPlaytform=iOS
    elif [ "${ROA_G_Platform}" == "android" ]
    then
        CurPlaytform=Android
    fi

    UploadAssetsRoot=${PathchDir}/${CurPlaytform}/${ROA_G_Branch}
    playerVersionFile=${UploadAssetsRoot}/CurPlayerVersion.txt
    playerVersion=`sed -n '1p' ${playerVersionFile}`
    TargetPath=${UploadAssetsRoot}/pack_${CurPlaytform}_${playerVersion}/UploadAssetsAfterPack/ServerData
    TargetAllAessetPath=${UploadAssetsRoot}/pack_${CurPlaytform}_${playerVersion}/UploadAftrPackAllAssets/ServerData
    TargetPatchPath=${UploadAssetsRoot}/pack_${CurPlaytform}_${playerVersion}/UploadAssetsAfterPatch/ServerData

    echo "playerVersionFile:"${playerVersionFile}
    echo "playerVersion:"${playerVersion}
    echo "待上传路径:"${TargetPath}

    # AAB包 且 master版本 上传资源（1:DLC+remote资源  2:补丁资源  3:全量资源）
    if [[ "${ROA_G_Branch}" == *branches* ]]
    then
        echo "当前是master分支，开始自动上传资源"
        echo "开始上传资源(DLC + remote)"
        aws s3 cp ${TargetPath} s3://ipretobucket/orig-roa/ --recursive
        echo "必要资源上传完成！QA已经可以开始测试！"

        if [ -e ${TargetPatchPath} ]
        then
            echo "开始上传整包补丁资源"
            aws s3 cp ${TargetPatchPath} s3://ipretobucket/orig-roa/ --recursive
            echo "补丁资源上传完成！"
        else
            echo "没有整包补丁资源"
        fi

        # echo "开始上传全量资源"
        aws s3 cp ${TargetAllAessetPath} s3://ipretobucket/orig-roa/ --recursive
        # echo "全量资源上传完成！"
    elif [ "${ROA_G_Branch}" == "develop" ]
    then
        echo "develop分支复制资源到本地服务器，安卓机访问地址：http://10.23.26.34/ 苹果机访问地址：http://10.23.26.35/"
        echo "如果无法访问请开启服务器"
        echo "启动apache服务器命令：sudo apachectl start"
        echo "关闭apache服务器命令：sudo apachectl stop"
        # 打开目录写入权限命令：sudo chmod -R o+w /Library/WebServer/Documents

        echo "开始拷贝资源(DLC + remote)到本地服务器"
        cp -rvf ${TargetPath}/* /Library/WebServer/Documents/
        echo "开始拷贝补丁资源到本地服务器"
        cp -rvf ${TargetPatchPath}/* /Library/WebServer/Documents/
        # echo "开始拷贝全量资源到本地服务器"
        # cp -rvf ${TargetAllAessetPath}/* /Library/WebServer/Documents/
    else
        echo "非release分支不会自动上传资源到CDN!"
	    MsgPush ${MSG_TITLE} "非release分支不会自动上传资源到CDN! "
    fi
    MsgPush ${MSG_TITLE} "资源上传完成！"
}
export -f DoUploadAll
