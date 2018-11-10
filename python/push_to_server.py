import asyncio
import os
import subprocess
import sys
from datetime import datetime
from zipfile import ZipFile

from paramiko import SSHClient
from scp import SCPClient


class SSHInfo(object):

    def __init__(self, hostname="localhost", port=22, username="root", password="password"):
        self.hostname = hostname
        self.port = port
        self.username = username
        self.password = password

    def __iter__(self):
        return (i for i in (self.hostname, self.port, self.username, self.password))

    def __str__(self):
        return "SSHInfo(hostname={!r}, port={!r}, username={!r}, password={!r})".format(*self)


def get_all_file_paths(directory):
    file_paths = []

    for root, directories, files in os.walk(directory):
        for filename in files:
            filepath = os.path.join(root, filename)
            if "target" not in filepath or ".svn" not in filepath:
                file_paths.append(filepath)

    return file_paths


def create_zip(zip_file_path: str, directory: str):
    file_paths = get_all_file_paths(directory)

    print('Following files will be zipped:')

    with ZipFile(zip_file_path, 'w') as zip:
        for file in file_paths:
            print(os.path.basename(file))
            zip.write(file)

    print('All files zipped successfully!')


def execute(command):
    process = subprocess.Popen(
        command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)

    # Poll process for new output until finished
    while True:
        nextline = process.stdout.readline()
        if nextline == '' and process.poll() is not None:
            break
        sys.stdout.write(nextline)
        sys.stdout.flush()

    output = process.communicate()[0]
    exit_code = process.returncode

    if (exit_code == 0):
        return output
    else:
        pass
        # raise ProcessException(command, exit_code, output)


def build(source_dir: str):
    cmd = 'cd {} && mvn clean install "-Dmaven.test.skip=true"'.format(
        source_dir)
    execute(cmd)


async def push(source_path: str, ssh: SSHClient):
    # Define progress callback that prints the current percentage completed for the file
    def progress(filename, size, sent):
        sys.stdout.write("%s\'s progress: %.2f%%   \r" %
                         (filename, float(sent) / float(size) * 100))

    with SCPClient(ssh.get_transport(), progress=progress) as scp:
        await scp.put("code.zip", source_path)


async def remote_unzip(ssh: SSHClient, target_dir: str):
    cmd = "unzip code.zip -d {} && rm code.zip".format(target_dir)
    await ssh.exec_command(cmd)


async def remote_build(ssh: SSHClient, target_dir: str):
    cmd = "cd {} && mvn clean install -Dmaven.test.skip=true".format(target_dir)
    await ssh.exec_command(cmd)


async def remote_deploy(ssh: SSHClient, target_dir: str):
    cmd = "cd {} && mvn clean install -Dmaven.test.skip=true".format(target_dir)
    await ssh.exec_command(cmd)


def main():
    create_zip("", ".")
    loop = asyncio.get_event_loop()

    ssh_info = SSHInfo()
    ssh = SSHClient()
    ssh.connect(*ssh_info)
    remote_tasks = [push("", ssh), remote_unzip(ssh, "")]
    futures = [asyncio.ensure_future(t, loop=loop) for t in remote_tasks]
    gathered = asyncio.gather(*futures, loop=loop, return_exceptions=True)
    loop.run_until_complete(gathered)

    for fut in futures:
        try:
            yield fut.result()
        except Exception as e:
            yield repr(e)


if __name__ == "__main__":
    start_time = datetime.now()
    if len(sys.argv) < 2:
        print("use the script like: %s code_dir remote_host target_dir" %
              sys.argv[0])
        sys.exit()

    source_dir = sys.argv[1]
    main()
    print("done")
    end_time = datetime.now()
    total = end_time - start_time
    print("used time: " + str(total))
