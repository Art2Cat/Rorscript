import subprocess
import sys
import os
from datetime import datetime
from paramiko import SSHClient
from scp import SCPClient


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
    exitCode = process.returncode

    if (exitCode == 0):
        return output
    else:
        pass
        # raise ProcessException(command, exitCode, output)


def build(source_dir: str):
    cmd = 'cd {} && mvn clean install "-Dmaven.test.skip=true"'.format(
        source_dir)
    execute(cmd)


def push(source_dir: str):
    ssh = SSHClient()
    ssh.load_system_host_keys()
    ssh.connect('user@server:path')

    # Define progress callback that prints the current percentage completed for the file
    def progress(filename, size, sent):
        sys.stdout.write("%s\'s progress: %.2f%%   \r" %
                         (filename, float(sent)/float(size)*100))

    with SCPClient(ssh.get_transport(), progress=progress) as scp:
        scp.put('my_file.txt', 'my_file.txt')

    for root, _, files in os.walk(source_dir):
        for file in files:
            new_dir = os.path.join(root, os.path.basename(file))
            # print(file)
            if "target" in new_dir and ".jar" in new_dir:
                print(new_dir)


def main():
    pass


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
