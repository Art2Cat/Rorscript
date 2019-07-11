import os
import re
import shlex
import subprocess
import sys
import zipfile
from pathlib import Path

import requests


def get_platform():
    platforms = {
        'linux': 'Linux',
        'darwin': 'OS X',
        'win32': 'windows'
    }
    if sys.platform not in platforms:
        return sys.platform

    return platforms[sys.platform]


def is_windows() -> bool:
    return sys.platform == "win32"


def download_file(url: str) -> str:
    local_filename = url.split('/')[-1]
    # NOTE the stream=True parameter below
    with requests.get(url, stream=True) as r:
        r.raise_for_status()
        with open(local_filename, 'wb') as f:
            for chunk in r.iter_content(chunk_size=8192):
                if chunk:
                    # filter out keep-alive new chunks
                    f.write(chunk)
                    # f.flush()
    return local_filename


def get_download_url(url: str) -> str:
    res = requests.get(url=url)

    data = str(res.text)
    if is_windows():
        dd = "https:.*windows_amd64.zip"
    elif os.uname()[4][:3] == "arm":
        dd = "https:.*" + sys.platform + "_" + os.uname()[4][:3] + ".tar.gz"
    else:
        dd = "https:.*" + sys.platform + "_amd64.tar.gz"

    # requests.get("https://github.com/fatedier/frp/releases/download/v0.27.0/frp_0.27.0_windows_amd64.zip")
    download_url = ""
    p = re.compile(dd)
    for s in data.split(","):
        # print(s)
        sd = re.search(p, s)
        if sd is not None:
            download_url = sd[0]

        print(download_url)
    return download_url


def main():
    if len(sys.argv) == 0:
        print("for frpc use like: {0} client\n for frps use like: {0} server".format(sys.argv[0]))

    ver = sys.argv[1]
    url = "https://api.github.com/repos/fatedier/frp/releases/latest"
    file = download_file(get_download_url(url))

    if is_windows():
        with zipfile.ZipFile(file, "r") as zip_ref:
            zip_ref.extractall(os.getcwd())

    else:
        with subprocess.Popen("tar -xzvf {} -C /usr/local/".format(file), stdout=subprocess.PIPE, shell=True) as proc:
            print(proc.stdout.read())
        frp = Path("/usr/local/", file.split(".")[0])
        if frp.exists():
            frp.rename("/usr/local/frp")

        cmds = ["systemctl daemon-reload"]

        service = """
        [Unit]
        Description=frp {0}


        [Service]
        Type=forking
        ExecStart=/usr/local/frp/{1} -c /usr/local/{1}.ini
        Restart=on-failure
        RestartSec=10
        KillMode=process

        [Install]
        WantedBy=multi-user.target
        """
        if ver == "client":
            service = service.format("client", "frpc")
            service_file = Path("/etc/systemd/system/frpc.service")
            cmds.append("systemctl enable frpc.service")
            cmds.append("systemctl start frpc.service")
        elif ver == "server":
            service = service.format("server", "frps")
            service_file = Path("/etc/systemd/system/frps.service")
            cmds.append("systemctl enable frps.service")
            cmds.append("systemctl start frps.service")
        else:
            raise

        with service_file.open("w") as f:
            f.write(service)

        for c in cmds:
            args = shlex.split(c)
            with subprocess.Popen(args, stdout=subprocess.PIPE, shell=True) as proc:
                print(proc.stdout.read())


if __name__ == '__main__':
    main()
