import os
import shutil
import socket
import subprocess
import sys
from datetime import datetime

from paramiko import SSHClient, WarningPolicy, SSHException
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


def execute(command):
    process = subprocess.Popen(
        command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)

    # Poll process for new output until finished
    while True:
        next_line = process.stdout.readline()
        if next_line == b'' and process.poll() is not None:
            break
        sys.stdout.write(str(next_line))
        sys.stdout.flush()

    output = process.communicate()[0]
    exit_code = process.returncode

    if exit_code == 0:
        return output
    else:
        pass
        # raise ProcessException(command, exit_code, output)


def push(source_path: str, ssh: SSHClient):
    # Define progress callback that prints the current percentage completed for the file
    def progress(file_name, size, sent):
        sys.stdout.write("%s\'s progress: %.2f%%   \r" %
                         (file_name, float(sent) / float(size) * 100))

    with SCPClient(ssh.get_transport(), progress=progress) as scp:
        scp.put(source_path, source_path)


def execute_command(client, commands):
    """Execute a command on the remote host.Return a tuple containing
    an integer status and a two strings, the first containing stdout
        and the second containing stderr from the command."""
    result_flag = True
    command = ""
    try:
        for cmd in commands:
            command = cmd
            print("Executing command --> {}".format(command))
            _, stdout, stderr = client.exec_command(command)
            ssh_output = stdout.read()
            ssh_error = stderr.read()
            for i in ssh_output.split(b"\n"):
                print(str(i))

            if ssh_error:
                print("Problem occurred while running command:" + command + " The error is " + str(ssh_error))
                result_flag = False
            else:
                print("Command execution completed successfully", command)
                result_flag = True
    except socket.timeout as e:
        print("Command timed out.", command)
        print(e)
        client.close()
        result_flag = False
    except SSHException:
        print("Failed to execute the command!", command)
        client.close()
        result_flag = False
    return result_flag


def main():
    zip_file_path = os.path.join(current_dir, filename)
    shutil.make_archive(filename, 'zip', zip_file_path)

    dbupgrade_path = os.path.join(zip_file_path, "bin", "deploy", "dbupgrade.bat")

    execute("{} 192.168.1.44 markit markit123 {}".format(dbupgrade_path, dbschema))
    ssh_info = SSHInfo("192.168.0.59", password="abc@123")
    ssh = SSHClient()
    ssh.load_system_host_keys()
    ssh.set_missing_host_key_policy(WarningPolicy())
    ssh.connect(*ssh_info)
    target_path = "/home/lp_app/lp_live/{}".format(filename)

    push("{}.zip".format(filename), ssh)
    commands = ['bash -c "unzip {}.zip -d {} && rm {}.zip"'.format(filename, target_path, filename),
                'cd {} && ./deploy.sh {}'.format(target_path, filename)]
    execute_command(ssh, commands)


if __name__ == "__main__":
    start_time = datetime.now()
    if len(sys.argv) < 2:
        print("use the script like: %s code_dir dbschema" %
              sys.argv[0])
        sys.exit()

    current_dir = os.path.abspath(".")
    filename = sys.argv[1]
    dbschema = sys.argv[2]
    main()
    # # asyncio.run(main(loop))

    print("done")
    end_time = datetime.now()
    total = end_time - start_time
    print("used time: " + str(total))
