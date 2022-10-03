#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Author: johlanse
Date: 2022/10/3 15:30
Description:: 用于一键安装study_xxqg️
"""
import io
import os
import platform
from functools import partial

try:
    import requests
except ImportError as e:
    print("requests依赖不存在，尝试安装依赖")
    os.system("pip3 install requests")
    import requests

if platform.system().lower() == "windows":
    import zipfile

print = partial(print, flush=True)


def updateDependent() -> str:
    """
    更新依赖的主函数
    """
    system = platform.system().lower()
    PyVersion_ = platform.python_version()
    if system == "windows":
        if platform.architecture()[0] == "64bit":
            fileName = f"study_xxqg_windows_amd64.zip"
            print(f"✅识别本机设备为Windows amd64,Py版本为{PyVersion_}\n")
        else:
            fileName = f"study_xxqg_windows_386.zip"
    elif system == "darwin":
        fileName = f"tlib_darwin_amd64.tar.gz"
        print(f"✅识别本机设备为MacOS x86_64,Py版本为{PyVersion_}\n")

    else:

        framework = os.uname().machine
        if framework == "x86_64":
            fileName = f"study_xxqg_linux_amd64.tar.gz"
            print(f"✅识别本机设备为Linux {framework},Py版本为{PyVersion_}\n")
        elif framework == "aarch64" or framework == "arm64":
            fileName = f"study_xxqg_linux_arm64.tar.gz"
            print(f"✅识别本机设备为Linux {framework},Py版本为{PyVersion_}\n")
        elif framework == "armv7l":
            fileName = f"study_xxqg_linux_386.tar.gz"
            print(f"✅识别本机设备为Linux {framework},Py版本为{PyVersion_}\n")
        else:
            fileName = f"study_xxqg_linux_amd64.tar.gz"
            print(f"⚠️无法识别本机设备操作系统,默认本机设备为Linux x86_64,Py版本为{PyVersion_}\n")
    return fileName


def last_version() -> str:
    return requests.get("https://api.github.com/repos/johlanse/study_xxqg/releases/latest").json().get("tag_name")


def download(github: str, version: str, binaryName: str):
    print("正在下载文件中，请耐心等待！！！")
    content = requests.get(f"{github}/johlanse/study_xxqg/releases/download/{version}/{binaryName}").content
    if platform.system().lower() == "windows":
        with zipfile.ZipFile(io.BytesIO(content)) as zf:
            data = zf.open("study_xxqg.exe")
            with open("study_xxqg.exe", "wb") as f:
                f.write(data.read())
    else:
        with open(binaryName, "wb") as f:
            f.write(content)
        os.system(f"tar xvf {binaryName}")
        os.remove(binaryName)


def checkYesOrNo() -> bool:
    data = input("请输入:").lower()
    if data == "y" or data == "yes":
        return True
    else:
        return False


def addSystemctl():
    with open("/etc/systemd/system/study_xxqg.service", "w", encoding="utf-8") as f:
        f.write(f'''[Unit]
    Description=study_xxqg
    Documentation=study_xxqg
    After=network-online.target
    Wants=network-online.target systemd-networkd-wait-online.service

    [Service]
    Restart=always

    ; User and group the process will run as.
    User=root
    Group=root

    WorkingDirectory={os.getcwd()}
    ExecStart={os.getcwd()}/study_xxqg

    ; Limit the number of file descriptors; see `man systemd.exec` for more limit settings.
    LimitNOFILE=1048576
    ; Unmodified caddy is not expected to use more than that.
    LimitNPROC=512

    [Install]
    WantedBy=multi-user.target''')


def main():
    github = "https://github.com"
    version = last_version()
    binaryName = updateDependent()
    if os.path.exists("study_xxqg.exe") or os.path.exists("study_xxqg"):
        print("检测到study_xxqg文件已经存在，是否跳过下载(y/n)")
        if not checkYesOrNo():
            download(github, version, binaryName)
        print("已跳过下载")
    else:
        download(github, version, binaryName)

    if not platform.system().lower() == "windows":
        print("是否将study_xxqg加入系统启动system命令(y/n)")
        if checkYesOrNo():
            addSystemctl()
            os.system("systemctl enable study_xxqg")
            print("已加入开机自启动，输入 systemctl start study_xxqg即可后台启动")
    print("开始安装浏览器依赖")
    try:
        os.mkdir("config")
    except FileExistsError as e:
        pass
    if platform.system().lower() == "windows":
        os.system("study_xxqg --init")
    else:
        os.system("./study_xxqg --init")
        os.system("./tools/driver/ms-playwright-go/1.20.0-beta-1647057403000/playwright install-deps")
    print("是否启动study_xxqg? (y/n)")
    if checkYesOrNo():
        if platform.system().lower() == "windows":
            os.system("study_xxqg.exe")
        else:
            os.system("chmod -R 777 ./study_xxqg")
            os.system("./study_xxqg")


if __name__ == '__main__':
    main()
