#!/bin/bash

echo "install clang 6.0 on debian stretch"
sudo sh -c "echo 'deb http://apt.llvm.org/stretch/ llvm-toolchain-stretch-6.0 main' >> /etc/apt/sources.list"

wget -O - https://apt.llvm.org/llvm-snapshot.gpg.key|sudo apt-key add - # Fingerprint: 6084 F3CF 814B 57C1 CF12 EFD5 15CF 4D18 AF4F 7421

sudo apt update

sudo apt install clang-6.0 clang-tools-6.0 clang-6.0-doc libclang-common-6.0-dev libclang-6.0-dev libclang1-6.0 libclang1-6.0-dbg libllvm6.0 libllvm6.0-dbg lldb-6.0 llvm-6.0 llvm-6.0-dev llvm-6.0-doc llvm-6.0-runtime clang-format-6.0 python-clang-6.0 libfuzzer-6.0-dev -y

sudo update-alternatives --install /usr/bin/clang++ clang++ /usr/bin/clang++-6.0 1000
sudo update-alternatives --install /usr/bin/clang clang /usr/bin/clang-6.0 1000
sudo update-alternatives --config clang
sudo update-alternatives --config clang++

echo "configure clang-format"
sudo cp /usr/bin/clang-format-6.0 /usr/bin/clang-format
