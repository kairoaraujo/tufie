#!/bin/bash

set -eou pipefail

REPO_BASE_URL=
# Create a temporary directory for downloaded files
TEMPDIR=$(mktemp -d)
# trap syscalls and delete the temporary directory
trap "rm -rf $TEMPDIR" EXIT

#install directory should be the first argument or default to /usr/local/bin
INSTALL_DIR=${1:-"/usr/local/bin"}

# resolve symlinks
INSTALL_DIR=$(readlink -f $INSTALL_DIR)

# check for bad install directory
if [ ! -d $INSTALL_DIR ]; then
    echo "Install directory $INSTALL_DIR does not exist"
    exit 1
fi

# Get the latest version of tufie
VERSION=$(curl -L -s https://api.github.com/repos/kairoaraujo/tufie/releases/latest | grep -o '"tag_name": *"[^"]*"' | sed 's/"//g' | sed 's/tag_name: *//')

#remove the 'v' from the version
readonly VERSION=${VERSION:1}

# Determine the architecture of the system
ARCH=$(uname -m)
if [ "$ARCH" == "x86_64" ]; then
    readonly ARCH="amd64"
fi
if [ "$ARCH" == "aarch64" ]; then
    readonly ARCH="arm64"
fi

### Determine if Arch is supported
if [ $ARCH != "amd64" ] && [ $ARCH != "arm64" ]; then
    echo "Unsupported architecture"
    exit 1
fi

OS=$(uname -s)

## change to lowercase
OS=$(echo "$OS" | tr '[:upper:]' '[:lower:]')
FILENAME="tufie_${VERSION}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/kairoaraujo/tufie/releases/download/v$VERSION/$FILENAME"
EXPECTED_CHECKSUM="$(curl -L -s  https://github.com/kairoaraujo/tufie/releases/download/v$VERSION/tufie_${VERSION}_checksums.txt | grep -w "tufie_${VERSION}_${OS}_$ARCH.tar.gz$" | awk '{print $1}')"
readonly EXPECTED_CHECKSUM EXPECTED_CHECKSUM

echo "Latest version of tufie is $VERSION"
echo "Downloading for $OS $ARCH from $DOWNLOAD_URL"

echo expected checksum: "$EXPECTED_CHECKSUM"


## Check to see if OS is supported
if [ $OS != "linux" ] && [ $OS != "darwin" ]; then
    echo "Unsupported OS"
    exit 1
fi

# Download the binary to the temporary directory
cd $TEMPDIR
curl -s -LO $DOWNLOAD_URL

# Verify the checksum
FILE_CHECKSUM=$(sha256sum -b tufie_${VERSION}_${OS}_${ARCH}.tar.gz | awk '{print $1}')

echo file checksum: "    $FILE_CHECKSUM"

if [ "$EXPECTED_CHECKSUM" != "$FILE_CHECKSUM" ]; then
    echo "Checksum verification failed, exiting"
    exit 1
fi

# Extract and install the binary
tar -xzf tufie_${VERSION}_${OS}_${ARCH}.tar.gz

# Check if the user has write permission for /usr/local/bin/tufie
if [ -w $INSTALL_DIR ]; then
  # If the user has write permission, move the binary to /usr/local/bin
  mv tufie $INSTALL_DIR
else
  # If the user doesn't have write permission, check if the script is running in an interactive terminal
  if [ -n "$TERM" ]; then
    # If the script is running in an interactive terminal, prompt for sudo and move the binary to /usr/local/bin
    echo "You don't have permission to move the file, if you would like to install tufie to ${INSTALL_DIR}, please enter your password."
    echo "You can also install tufie to a different directory by running this script with the directory as the first argument."
    sudo -v && sudo mv tufie $INSTALL_DIR
  else
    # If the script is not running in an interactive terminal, print an error message and exit
    echo "You don't have permission to move the file, please run with sudo."
    exit 1
  fi
fi

# Return to the original directory and remove the temporary directory silently
cd - > /dev/null
rm -rf $TEMPDIR

# Check if the binary is installed
tufie version
echo "tufie v${VERSION} has been installed at ${INSTALL_DIR}/tufie"