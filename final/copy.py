import os

# 排除的文件和文件夹
exclude_files = [
    "copy.py", "go.mod", "go.sum", "output.txt", "README.md",
    "proto/external-api/external-master_grpc.pb.go", "proto/external-api/external-master.pb.go",
    "proto/external-api/external-tablet_grpc.pb.go", "proto/external-api/external-tablet.pb.go",
    "proto/internal-api/internal-master_grpc.pb.go", "proto/internal-api/internal-master.pb.go",
    "proto/internal-api/internal-tablet_grpc.pb.go", "proto/internal-api/internal-tablet.pb.go"
]

exclude_dirs = [
    "test"
]

def write_file_content_to_output(file_path, output_file):
    with open(file_path, 'r', encoding='utf-8', errors='ignore') as file:
        output_file.write(f"--- 文件: {file_path} ---\n")
        output_file.write(file.read())
        output_file.write("\n\n")

def traverse_directory_and_copy_content_to_txt(directory, output_txt):
    with open(output_txt, 'w', encoding='utf-8') as output_file:
        for root, dirs, files in os.walk(directory):
            # 排除指定的文件夹
            dirs[:] = [d for d in dirs if d not in exclude_dirs]

            for file in files:
                file_path = os.path.join(root, file)

                # 确保文件路径完全匹配并排除
                if (file.endswith(('.go', '.proto')) and 
                    not any(file_path.endswith(exclude) for exclude in exclude_files) and 
                    file != os.path.basename(__file__)):

                    write_file_content_to_output(file_path, output_file)
                    print(f"正在处理文件: {file_path}")

if __name__ == "__main__":
    current_directory = os.getcwd()  # 获取当前文件夹路径
    output_txt_file = "output.txt"   # 输出的txt文件名
    traverse_directory_and_copy_content_to_txt(current_directory, output_txt_file)
    print(f"所有文件内容已复制到 {output_txt_file}")
