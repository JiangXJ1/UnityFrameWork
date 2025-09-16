import jsonUtil

# 文件路径
json_file_path = '../../../UnityProject/Assets/Scripts/Runtime/Static/Custom.Static.asmdef'

new_elements = [
    "PeapodSDK.Module.Base",
    "PeapodSDK.Module.FirebaseVendor",
    "PeapodSDK.Module.AppsflyerVendor"
]

jsonUtil.add_multiple_to_json_array(json_file_path, "references", new_elements)