import jsonUtil

# 文件路径
json_file_path = '../../../UnityProject/Packages/manifest.json'

new_elements = {
        "com.gpc.sdk.base": "file:com.gpc.sdk.base@0.21.0-su.1701-beta.2.0",
        "com.gpc.sdk.dependence.huawei": "file:com.gpc.sdk.dependence.huawei@0.21.0-su.17",
        "com.gpc.sdk.dependence.samsung": "file:com.gpc.sdk.dependence.samsung@0.21.0-su.17",
        "com.gpc.sdk.module.aiguide": "file:com.gpc.sdk.module.aiguide@0.1.0",
        "com.gpc.sdk.module.blacklistedword": "file:com.gpc.sdk.module.blacklistedword@0.4.1",
        "com.gpc.sdk.module.calendar": "file:com.gpc.sdk.module.calendar@0.2.0",
        "com.gpc.sdk.module.cpd": "file:com.gpc.sdk.module.cpd@0.2.1",
        "com.gpc.sdk.module.crypto": "file:com.gpc.sdk.module.crypto@0.2.0",
        "com.gpc.sdk.resource.theme.slategrey": "file:com.gpc.sdk.resource.theme.slategrey@0.21.0-su.17",
        "com.gpc.sdk.resource.i18n.slategrey": "file:com.gpc.sdk.resource.i18n.slategrey@0.21.0-su.17",
        "com.thirdparty.icsharpcode.sharpziplib": "file:com.thirdparty.icsharpcode.sharpzipli@0.2.0",
        "com.thirdparty.microsoft.bcl.async-interfaces": "file:com.thirdparty.microsoft.bcl.async-interfaces@0.1.0",
        "com.thirdparty.newtonsoft-json": "file:com.thirdparty.newtonsoft-json@0.2.0",
        "com.thirdparty.notch-solution": "file:com.thirdparty.notch-solution@0.2.0",
        "com.gpc.sdk.module.operations": "file:com.gpc.sdk.module.operations@0.2.0",
        "com.gpc.sdk.module.push": "file:com.gpc.sdk.module.push@0.1.2",
        "com.gpc.sdk.bridge": "file:com.gpc.sdk.bridge@0.5.0",
        "com.gpc.sdk.module.payment": "file:com.gpc.sdk.module.payment@0.2.0",
        "com.gpc.sdk.module.peapodsdkproxy": "file:com.gpc.sdk.module.peapodsdkproxy@0.1.0",
        "com.gpc.sdk.peapod.base": "file:com.gpc.sdk.peapod.base@1.2.0",
        "com.gpc.sdk.peapod.aws.vendor": "file:com.gpc.sdk.peapod.aws.vendor@1.2.0",
        "com.gpc.sdk.peapod.appsflyer.vendor": "file:com.gpc.sdk.peapod.appsflyer.vendor@1.2.0",
        "com.gpc.sdk.peapod.firebase.vendor": "file:com.gpc.sdk.peapod.firebase.vendor@1.2.0",
        "com.thirdparty.system.runtime.compilerservices.unsafe": "file:com.thirdparty.system.runtime.compilerservices.unsafe@0.1.0",
        "com.thirdparty.system.threading.tasks.extensions": "file:com.thirdparty.system.threading.tasks.extensions@0.1.0",
        "com.gpc.sdk.diting": "file:com.gpc.sdk.diting@0.21.0-su.17"
}

jsonUtil.add_unique_to_json_dict(json_file_path, "dependencies", new_elements)