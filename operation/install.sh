#!/bin/bash
#
# install.sh - ZwischentonCloud Installer Script
#
# (c)2020-2025 Simon Gaus
#

# ─────────────────────────────────────────────────────────────────────────────
# 🛑 Check if all parameters are supplied
# ─────────────────────────────────────────────────────────────────────────────
if [ $# -ne 3 ]; then
    echo -e "\n\033[1;31m🚨  ERROR: Missing parameters!\033[0m"
    echo -e "\033[1;34m🔹  USAGE:\033[0m sudo ./install.sh \033[1;32m<mysql_root_pw> <mysql_backup_pw> <database_pw>\033[0m"
    echo -e "\033[1;31m❌  Exiting.\033[0m\n"
    exit 1
fi

# ─────────────────────────────────────────────────────────────────────────────
# 🎯 Store parameters in variables
# ─────────────────────────────────────────────────────────────────────────────
root_password="$1"
backup_password="$2"
database_password="$3"

# ─────────────────────────────────────────────────────────────────────────────
# 🔍 Detect Web Server User
# ─────────────────────────────────────────────────────────────────────────────
WEB_USER="www-data"
if ! id -u "$WEB_USER" &>/dev/null; then
    WEB_USER="www"
    if ! id -u "$WEB_USER" &>/dev/null; then
        echo -e "\n\033[1;31m❌  ERROR: Web server user not found! Exiting.\033[0m\n"
        exit 1
    fi
fi

# ─────────────────────────────────────────────────────────────────────────────
# 🔍 Set Database System User
# ─────────────────────────────────────────────────────────────────────────────
database_user="mysql"

# ─────────────────────────────────────────────────────────────────────────────
# 📁 Setup Working Directory
# ─────────────────────────────────────────────────────────────────────────────
WORK_DIR="/usr/local/zwischentoncloud/install"
mkdir -p "$WORK_DIR" && cd "$WORK_DIR" || { echo -e "\n\033[1;31m❌  ERROR: Failed to create/access working directory!\033[0m\n"; exit 1; }
echo -e "\n📂  Working directory set to \e[1;34m$WORK_DIR\e[0m"
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🖥  Detect System OS and Architecture
# ─────────────────────────────────────────────────────────────────────────────
if [ "$(uname -s)" = "Darwin" ]; then
    os="darwin"
elif [ "$(uname -s)" = "Linux" ]; then
    os="linux"
else
    echo -e "\n🚨  ERROR: Unsupported OS. Exiting.\n"
    exit 1
fi
if [ "$(uname -m)" = "x86_64" ]; then
    arch="amd64"
elif [ "$(uname -m)" = "arm64" ]; then
    arch="arm64"
else
    echo -e "\n🚨  ERROR: Unsupported CPU architecture. Exiting.\n"
    exit 1
fi

# ─────────────────────────────────────────────────────────────────────────────
# 📦 Download latest release
# ─────────────────────────────────────────────────────────────────────────────
file_url="https://github.com/TR-Projekt/zwischentoncloud/releases/latest/download/zwischentoncloud-$os-$arch.tar.gz"
echo -e "\n📥  Downloading latest ZwischentonCloud release..."
curl --progress-bar -L "$file_url" -o zwischentoncloud.tar.gz
echo -e "📦  Extracting archive..."
tar -xf zwischentoncloud.tar.gz

# ─────────────────────────────────────────────────────────────────────────────
# 📦 Install & Enable & Start MySQL Server
# ─────────────────────────────────────────────────────────────────────────────
echo -e "\n🗂️  Installing MySQL server..."
apt-get install mysql-server -y > /dev/null 2>&1
systemctl enable mysql &>/dev/null && systemctl start mysql &>/dev/null
echo -e "✅  MySQL service is up and running."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🔐 Install MySQL Backup Credential File
# ─────────────────────────────────────────────────────────────────────────────
credentialsFile=/usr/local/zwischentoncloud/mysql.conf
cat << EOF > $credentialsFile
# zwischentoncloud configuration file v1.0
# TOML 1.0.0-rc.2+

[client]
user = 'zwischentoncloud.backup'
password = '$backup_password'
host = 'localhost'
EOF
if [ -f "$credentialsFile" ]; then
    echo -e "✅  MySQL backup credential file successfully created at \e[1;34m$credentialsFile\e[0m"
else
    echo -e "🚨  ERROR: Failed to create MySQL credential file. Exiting.\n"
    exit 1
fi
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🔑 Secure MySQL
# ─────────────────────────────────────────────────────────────────────────────
chmod +x secure-mysql.sh
./secure-mysql.sh "$root_password"
echo -e "✅  MySQL security script executed."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🗄️  Setup Databases & Users
# ─────────────────────────────────────────────────────────────────────────────
mysql -e "source $WORK_DIR/create_identity_database.sql"
mysql -e "CREATE USER 'zwischentoncloud.identity.writer'@'localhost' IDENTIFIED BY '$database_password';"
mysql -e "GRANT SELECT, INSERT, UPDATE, DELETE ON zwischenton_identity_database.* TO 'zwischentoncloud.identity.writer'@'localhost';"

mysql -e "source $WORK_DIR/create_zwischenton_database.sql"
mysql -e "CREATE USER 'zwischentoncloud.api.writer'@'localhost' IDENTIFIED BY '$database_password';"
mysql -e "GRANT SELECT, INSERT, UPDATE, DELETE ON zwischenton_cloud_database.* TO 'zwischentoncloud.api.writer'@'localhost';"

mysql -e "CREATE USER 'zwischentoncloud.backup'@'localhost' IDENTIFIED BY '$backup_password';"
mysql -e "GRANT ALL PRIVILEGES ON zwischenton_identity_database.* TO 'zwischentoncloud.backup'@'localhost';"
mysql -e "GRANT ALL PRIVILEGES ON zwischenton_cloud_database.* TO 'zwischentoncloud.backup'@'localhost';"
mysql -e "FLUSH PRIVILEGES;"
echo -e "✅  Database and users created."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 📂 Setup Database Backup Directory
# ─────────────────────────────────────────────────────────────────────────────
mkdir -p /srv/zwischentoncloud/backups
mv backup.sh /srv/zwischentoncloud/backups/backup.sh
chmod +x /srv/zwischentoncloud/backups/backup.sh
echo -e "✅  Database backup directory and script configured."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# ⏳ Install Cronjob for Daily Backup at 3 AM
# ─────────────────────────────────────────────────────────────────────────────
echo -e "0 3 * * * $WEB_USER /srv/zwischentoncloud/backups/backup.sh" | tee -a /etc/cron.d/zwischentoncloud_backup > /dev/null
echo -e "✅  Cronjob installed! Backup will run daily at \e[1;34m3 AM\e[0m"
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 📦 Install ZwischentonCloud
# ─────────────────────────────────────────────────────────────────────────────
echo -e "\n📥  Installing latest ZwischentonCloud binary..."
mv zwischentoncloud /usr/local/bin/zwischentoncloud || {
    echo -e "\n🚨  ERROR: Failed to install Zwischenton Cloud binary. Exiting.\n"
    exit 1
}
echo -e "✅  Installed ZwischentonCloud to \e[1;34m/usr/local/bin/zwischentoncloud\e[0m."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🛠  Install Server Configuration File
# ─────────────────────────────────────────────────────────────────────────────
mv config_template.toml /etc/zwischentoncloud.conf
if [ -f "/etc/zwischentoncloud.conf" ]; then
    echo -e "✅  Configuration file moved to \e[1;34m/etc/zwischentoncloud.conf\e[0m."
else
    echo -e "\n🚨  ERROR: Failed to move configuration file. Exiting.\n"
    exit 1
fi
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 📂  Prepare Log Directory
# ─────────────────────────────────────────────────────────────────────────────
mkdir -p /var/log/zwischentoncloud || {
    echo -e "\n🚨  ERROR: Failed to create log directory. Exiting.\n"
    exit 1
}
echo -e "✅  Log directory created at \e[1;34m/var/log/zwischentoncloud\e[0m."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🔄 Prepare Remote Update Workflow
# ─────────────────────────────────────────────────────────────────────────────
mv update.sh /usr/local/zwischentoncloud/update.sh
chmod +x /usr/local/zwischentoncloud/update.sh
cp /etc/sudoers /tmp/sudoers.bak
echo "$WEB_USER ALL = (ALL) NOPASSWD: /usr/local/zwischentoncloud/update.sh" >> /tmp/sudoers.bak
# Validate and replace sudoers file if syntax is correct
if visudo -cf /tmp/sudoers.bak &>/dev/null; then
    sudo cp /tmp/sudoers.bak /etc/sudoers
    echo -e "✅  Prepared remote update workflow."
else
    echo -e "\n🚨  ERROR: Could not modify /etc/sudoers file. Please do this manually. Exiting.\n"
    exit 1
fi
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🔥 Enable and Configure Firewall
# ─────────────────────────────────────────────────────────────────────────────
if command -v ufw > /dev/null; then
    echo -e "\n🔥  Configuring UFW firewall..."
    mv ufw_app_profile /etc/ufw/applications.d/zwischentoncloud
    ufw allow zwischentoncloud > /dev/null
    echo -e "✅  Added zwischentoncloud to UFW with port 2340."
    sleep 1
elif ! [ "$(uname -s)" = "Darwin" ]; then
    echo -e "\n🚨  ERROR: No firewall detected and not on macOS. Exiting.\n"
    exit 1
fi

# ─────────────────────────────────────────────────────────────────────────────
# ⚙️  Install Systemd Service
# ─────────────────────────────────────────────────────────────────────────────
if command -v service > /dev/null; then
    echo -e "\n🚀  Configuring systemd service..."
    if ! [ -f "/etc/systemd/system/fzwischentoncloud.service" ]; then
        mv service_template.service /etc/systemd/system/zwischentoncloud.service
        echo -e "✅  Created systemd service configuration."
        sleep 1
    fi
    systemctl enable zwischentoncloud > /dev/null
    echo -e "✅  Enabled systemd service for ZwischentonCloud."
    sleep 1
elif ! [ "$(uname -s)" = "Darwin" ]; then
    echo -e "\n🚨  ERROR: Systemd is missing and not on macOS. Exiting.\n"
    exit 1
fi

# ─────────────────────────────────────────────────────────────────────────────
# 🔑 Set Appropriate Permissions
# ─────────────────────────────────────────────────────────────────────────────
chown -R "$WEB_USER":"$WEB_USER" /usr/local/zwischentoncloud
chown -R "$WEB_USER":"$WEB_USER" /var/log/zwischentoncloud
chown -R "$WEB_USER":"$WEB_USER" /srv/zwischentoncloud
chown "$WEB_USER":"$WEB_USER" /etc/zwischentoncloud.conf
echo -e "\n🔐  Set Appropriate Permissions."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🧹 Cleanup Installation Files
# ─────────────────────────────────────────────────────────────────────────────
echo -e "🧹  Cleaning up installation files..."
cd /usr/local/zwischentoncloud || exit
rm -rf /usr/local/zwischentoncloud/install
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🎉 COMPLETE Message
# ─────────────────────────────────────────────────────────────────────────────
echo -e "\n\033[1;32m══════════════════════════════════════════════════════════════════════════\033[0m"
echo -e "\033[1;32m✅  INSTALLATION COMPLETE! 🚀\033[0m"
echo -e "\033[1;32m══════════════════════════════════════════════════════════════════════════\033[0m"
echo -e "\n📂 \033[1;34mBefore starting, you need to:\033[0m"
echo -e "\n   \033[1;34m1. Configure the mTLS certificates.\033[0m"
echo -e "   \033[1;34m2. Configure the JWT signing keys.\033[0m"
echo -e "   \033[1;34m3. Update the configuration file at:\033[0m"
echo -e "\n   \033[1;32m    /etc/zwischentoncloud.conf\033[0m"
echo -e "\n🔹 \033[1;34mThen start the server manually:\033[0m"
echo -e "\n   \033[1;32m    sudo systemctl start zwischentoncloud\033[0m"
echo -e "\n\033[1;32m══════════════════════════════════════════════════════════════════════════\033[0m\n"
sleep 1
