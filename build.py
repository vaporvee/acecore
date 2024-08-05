import os
import subprocess
import shutil

app_name = 'acecore'

src_dir = './plugin_src'
output_dir = './build'

subprocess.run('go build .', shell=True, check=True)
shutil.copy(app_name, output_dir)
print(f'Built main executable in {output_dir}/{app_name}')

if not os.path.exists(output_dir):
    os.makedirs(output_dir)

for folder_name in os.listdir(src_dir):
    folder_path = os.path.join(src_dir, folder_name)
    if os.path.isdir(folder_path):
        command = f'go build -buildmode=plugin -o {output_dir}/plugins/{folder_name}.so {folder_path}'
        subprocess.run(command, shell=True, check=True)
        print(f'Built plugin: {folder_name}.so')

