import json

# 读取现有的 JSON 文件
def read_json(file_path):
    try:
        with open(file_path, 'r', encoding='utf-8') as file:
            return json.load(file)
    except FileNotFoundError:
        # 如果文件不存在，返回一个空字典
        return {}


# 向 JSON 文件中的数组节点添加多个内容
def add_multiple_to_json_array(file_path, array_key, new_items):
    # 读取现有的数据
    data = read_json(file_path)

    # 确保指定的键是一个数组
    if array_key not in data:
        data[array_key] = []  # 如果键不存在，则初始化为空数组
    elif not isinstance(data[array_key], list):
        raise ValueError(f"The value associated with '{array_key}' is not a list.")

    # 当前数组中的内容
    current_items = data[array_key]
    
    # 添加新的内容之前先检查是否已存在
    for item in new_items:
        if item not in current_items:
            current_items.append(item)  # 如果元素不在数组中则添加

    # 更新字典中的数组
    data[array_key] = current_items

    # 将修改后的数据写回文件
    with open(file_path, 'w', encoding='utf-8') as file:
        json.dump(data, file, ensure_ascii=False, indent=4)

    print(f"added to array in {file_path}")


# 向 JSON 文件的字典节点动态添加多个键值对，并检查是否已存在
def add_unique_to_json_dict(file_path, dict_key, new_items):
    # 读取现有的数据
    data = read_json(file_path)
    
    # 确保指定的键是一个字典
    if dict_key not in data:
        data[dict_key] = {}  # 如果键不存在，则初始化为空字典
    elif not isinstance(data[dict_key], dict):
        raise ValueError(f"The value associated with '{dict_key}' is not a dictionary.")
    
    # 当前字典中的内容
    current_dict = data[dict_key]
    
    # 添加新的键值对之前先检查是否已存在
    for key, value in new_items.items():
        if key not in current_dict:
            current_dict[key] = value  # 如果键不在字典中则添加
    
    # 更新字典中的内容
    data[dict_key] = current_dict

    # 将修改后的数据写回文件
    with open(file_path, 'w', encoding='utf-8') as file:
        json.dump(data, file, ensure_ascii=False, indent=4)
    
    print(f"added to dic in {file_path}")
