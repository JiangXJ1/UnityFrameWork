import os
import shutil
import stat

def copy_folder_contents(src_folder, dest_folder):
    # 检查源文件夹是否存在
    if not os.path.exists(src_folder):
        print(f"源文件夹 {src_folder} 不存在!")
        return

    # 如果目标文件夹不存在，则创建它
    if not os.path.exists(dest_folder):
        os.makedirs(dest_folder)

    # 遍历源文件夹中的所有文件和子目录
    for item in os.listdir(src_folder):
        src_item = os.path.join(src_folder, item)
        dest_item = os.path.join(dest_folder, item)

        # 如果是文件，则复制
        if os.path.isfile(src_item):
            shutil.copy2(src_item, dest_item)
            os.chmod(dest_item, stat.S_IRWXU)  # 用户可读、写、执行权限
        # 如果是目录，则递归复制
        elif os.path.isdir(src_item):
            shutil.copytree(src_item, dest_item, dirs_exist_ok=True)

    print(f"已成功将 {src_folder} 下的所有内容拷贝到 {dest_folder}.")