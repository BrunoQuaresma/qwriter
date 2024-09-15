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

# Instructions
echo "Before you can start using QWriter, you'll need to set your OpenAI API key."
echo "This key allows QWriter to connect to OpenAI's powerful language models."
echo ""
echo "1. $(tput bold)Generate an API Key:$(tput sgr0)"
echo "   Follow the instructions in the OpenAI quickstart guide:"
echo "   $(tput smul)https://platform.openai.com/docs/quickstart/create-and-export-an-api-key$(tput sgr0)"
echo ""
echo "2. $(tput bold)Export the API Key:$(tput sgr0)"
echo "   Set the OPENAI_API_KEY environment variable with your API key:"
echo "   export OPENAI_API_KEY=your-key-here"
echo ""