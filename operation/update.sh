#!/bin/bash
#
# update.sh 1.0.0
#
# Updates the zwischentoncloud and restarts it.
#
# (c)2024 Simon Gaus
#

# Move to working dir
#
mkdir /usr/local/zwischentoncloud/install || { echo "Failed to create working directory. Exiting." ; exit 1; }
cd /usr/local/zwischentoncloud/install || { echo "Failed to access working directory. Exiting." ; exit 1; }

# Get system os
#
if [ "$(uname -s)" = "Darwin" ]; then
  os="darwin"
elif [ "$(uname -s)" = "Linux" ]; then
  os="linux"
else
  echo "System is not Darwin or Linux. Exiting."
  exit 1
fi

# Get systems cpu architecture
#
if [ "$(uname -m)" = "x86_64" ]; then
  arch="amd64"
elif [ "$(uname -m)" = "arm64" ]; then
  arch="arm64"
else
  echo "System is not x86_64 or arm64. Exiting."
  exit 1
fi

# Build url to latest binary for the given system
#
file_url="https://github.com/TR-Projekt/zwischentoncloud/releases/latest/download/zwischentoncloud-$os-$arch.tar.gz"
echo "The system is $os on $arch."
sleep 1

# Updating zwischentoncloud to the newest binary release
#
echo "Downloading newest zwischentoncloud binary release..."
curl -L "$file_url" -o zwischentoncloud.tar.gz
tar -xf zwischentoncloud.tar.gz
mv zwischentoncloud /usr/local/bin/zwischentoncloud || { echo "Failed to install zwischentoncloud binary. Exiting." ; exit 1; }
echo "Updated zwischentoncloud binary."
sleep 1

# Removing unused files
#
echo "Cleanup..."
cd /usr/local/zwischentoncloud || { echo "Failed to access server directory. Exiting." ; exit 1; }
rm -r /usr/local/zwischentoncloud/install
sleep 1

# Restart the zwischentoncloud
#
systemctl restart zwischentoncloud
echo "Restarted the zwischentoncloud"
sleep 1

echo "Done!"
sleep 1 