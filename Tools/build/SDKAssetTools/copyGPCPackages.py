# import shutil
# import os

# def copy_contents(src, dst):
#     # 确保目标目录存在
#     if not os.path.exists(dst):
#         os.makedirs(dst)

#     # 遍历源目录中的所有文件和子目录
#     for item in os.listdir(src):
#         # 构造源和目标路径
#         src_path = os.path.join(src, item)
#         dst_path = os.path.join(dst, item)

#         # 复制文件或目录
#         if os.path.isdir(src_path):
#             shutil.copytree(src_path, dst_path, dirs_exist_ok=True)
#         else:
#             shutil.copy2(src_path, dst_path)

# # 示例路径
# src_directory = 'path/to/source_directory'
# dst_directory = 'path/to/destination_directory'

# # 调用函数
# copy_contents(src_directory, dst_directory)

