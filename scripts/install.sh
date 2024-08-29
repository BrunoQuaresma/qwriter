#!/bin/bash

# Get the operating system type
OS_TYPE=$(uname -s)

# Get the architecture type
ARCH_TYPE=$(uname -m)

# Determine OS
if [[ "$OS_TYPE" == "Linux" ]]; then
    OS="linux"
elif [[ "$OS_TYPE" == "Darwin" ]]; then
    OS="darwin"
else
    OS="unsupported"
fi

# Determine architecture
if [[ "$ARCH_TYPE" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH_TYPE" == "arm64" || "$ARCH_TYPE" == "aarch64" ]]; then
    ARCH="arm64"
else
    ARCH="unsupported"
fi

# Download the qwriter binary
curl -fsSL -o /tmp/qwriter "https://github.com/BrunoQuaresma/qwriter/releases/download/v0.1.2/qwriter-$OS-$ARCH"
chmod +x /tmp/qwriter

# Add the qwriter binary to the PATH
sudo mv /tmp/qwriter /usr/local/bin/qwriter
PATH=$PATH:/usr/local/bin/qwriter