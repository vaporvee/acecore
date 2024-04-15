import os
import subprocess

src_dir = './plugin_src'
output_dir = './plugins'

if not os.path.exists(output_dir):
    os.makedirs(output_dir)

for folder_name in os.listdir(src_dir):
    folder_path = os.path.join(src_dir, folder_name)
    if os.path.isdir(folder_path):
        command = f'go build -buildmode=plugin -o {output_dir}/{folder_name}.so {folder_path}'
        subprocess.run(command, shell=True, check=True)
        print(f'Built plugin: {folder_name}.so')
