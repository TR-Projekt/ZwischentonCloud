#!/bin/bash
#
# install.sh 1.0.0
#
# Enables the firewall, installs the newest zwischentoncloud and starts it as a service.
#
# (c)2024 Simon Gaus
#

# Test for web server user
#
WEB_USER="www-data"
id -u "$WEB_USER" &>/dev/null;
if [ $? -ne 0 ]; then
  WEB_USER="www"
  if [ $? -ne 0 ]; then
    echo "Failed to find user to run web server. Exiting."
    exit 1
  fi
fi

# Move to working dir
#
mkdir -p /usr/local/zwischentoncloud/install || { echo "Failed to create working directory. Exiting." ; exit 1; }
cd /usr/local/zwischentoncloud/install || { echo "Failed to access working directory. Exiting." ; exit 1; }
echo "Installing zwischentoncloud using port 443."
sleep 1

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

# Install zwischentoncloud to /usr/local/bin/zwischentoncloud. TODO: Maybe just link to /usr/local/bin?
#
echo "Downloading newest zwischentoncloud binary release..."
curl -L "$file_url" -o zwischentoncloud.tar.gz
tar -xf zwischentoncloud.tar.gz
mv zwischentoncloud /usr/local/bin/zwischentoncloud || { echo "Failed to install zwischentoncloud binary. Exiting." ; exit 1; }
echo "Installed the zwischentoncloud binary to '/usr/local/bin/zwischentoncloud'."
sleep 1

## Install server config file
mv config_template.toml /etc/zwischentoncloud.conf
echo "Moved default zwischentoncloud config to '/etc/zwischentoncloud.conf'."
sleep 1

## Prepare log directory
mkdir /var/log/zwischentoncloud || { echo "Failed to create log directory. Exiting." ; exit 1; }
echo "Create log directory at '/var/log/zwischentoncloud'."

## Prepare server update workflow
mv update.sh /usr/local/zwischentoncloud/update.sh
chmod +x /usr/local/zwischentoncloud/update.sh
cp /etc/sudoers /tmp/sudoers.bak
echo "$WEB_USER ALL = (ALL) NOPASSWD: /usr/local/zwischentoncloud/update.sh" >> /tmp/sudoers.bak
# Check syntax of the backup file to make sure it is correct.
visudo -cf /tmp/sudoers.bak
if [ $? -eq 0 ]; then
  # Replace the sudoers file with the new only if syntax is correct.
  sudo cp /tmp/sudoers.bak /etc/sudoers
else
  echo "Could not modify /etc/sudoers file. Please do this manually." ; exit 1;
fi

# Enable and configure the firewall.
#
if command -v ufw > /dev/null; then

  ufw allow https >/dev/null
  echo "Added zwischentoncloud to ufw using port 443."
  sleep 1

elif ! [ "$(uname -s)" = "Darwin" ]; then
  echo "No firewall detected and not on macOS. Exiting."
  exit 1
fi

# Install systemd service
#
if command -v service > /dev/null; then

  if ! [ -f "/etc/systemd/system/zwischentoncloud.service" ]; then
    mv service_template.service /etc/systemd/system/zwischentoncloud.service
    echo "Created systemd service."
    sleep 1
  fi

  systemctl enable zwischentoncloud > /dev/null
  echo "Enabled systemd service."
  sleep 1

elif ! [ "$(uname -s)" = "Darwin" ]; then
  echo "Systemd is missing and not on macOS. Exiting."
  exit 1
fi

## Set appropriate permissions
##
chown -R "$WEB_USER":"$WEB_USER" /usr/local/zwischentoncloud
chown -R "$WEB_USER":"$WEB_USER" /var/log/zwischentoncloud
chown "$WEB_USER":"$WEB_USER" /etc/zwischentoncloud.conf

# Download Zwischenton Root CA certificate
#--> to /usr/local/zwischentoncloud/ca.crt

# Remving unused files
#
echo "Cleanup..."
cd /usr/local/zwischentoncloud || exit
rm -R /usr/local/zwischentoncloud/install
sleep 1

echo "Done!"
sleep 1

echo "You can start the server manually by running 'sudo systemctl start zwischentoncloud' after you updated the configuration file at '/etc/zwischentoncloud.conf'"
sleep 1
