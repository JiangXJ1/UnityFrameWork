#!/bin/bash
cd ${WORKSPACE}/Tools/build
# 当前工作目录
pwd
# 添加权限
chmod +x BeforBuild.sh

echo "+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ Start.sh"

echo "------------------------------------Jenkins全局参数"
echo "工作目录:"${WORKSPACE}
echo  "版本号:"$ROA_G_Version
echo  "分支名:"$ROA_G_Branch
echo  "构建号:"$ROA_G_BundleVersionCode
echo  "模板名:"$ROA_G_TplFileOption
echo  "平台:"$ROA_G_Platform
echo  "渠道:"$ROA_G_ChannelName
echo  "SVN Revision:"$ROA_G_SVNRevision
echo  "Debug:"$ROA_G_Debug #字符串 “true” 或 “false”
echo  "IP:"$ROA_G_IP #字符串
echo  "USE_AAB参数:"${ROA_G_AAB}
echo "------------------------------------"

#解析自定义IP
if [ -n "$ROA_G_IP" ]
then
    if [[ $ROA_G_IP =~ "-" ]]; then
        IFS='-' read -ra parts <<< "$ROA_G_IP"
        ROA_G_IP="${parts[1]}"
        echo "${ROA_G_IP}"
    fi
else
    echo "ROA_G_IP is empty"
fi

#启动打包
if [ $ROA_G_Platform == "android" ] || [ $ROA_G_Platform == "ios" ]
then
    echo "打包（android）"
    sh BeforBuild.sh
fi