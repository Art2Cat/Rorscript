import subprocess
import sys
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


def push():
    ssh = SSHClient()
    ssh.load_system_host_keys()
    ssh.connect('user@server:path')
    with SCPClient(ssh.get_transport()) as scp:
        scp.put('my_file.txt', 'my_file.txt')


def main():
    pass


if __name__ == "__main__":
    start_time = datetime.now()
    if len(sys.argv) < 2:
        print("use the script like: %s code_dir remote_host target_idr" %
              sys.argv[0])
        sys.exit()

    source_dir = sys.argv[1]
    main()
    print("done")
    end_time = datetime.now()
    total = end_time - start_time
    print("used time: " + str(total))
