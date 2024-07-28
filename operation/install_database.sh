#!/bin/bash
#
# install.sh 1.0.0
#
# Enables the firewall, installs the newest mysql, starts it as a service,
# configures it to be used as the database for the ZwischentonCloud API and identity database and setup
# the backup routines.
#
# (c)2024 Simon Gaus
#

# Check if all passwords are supplied
#
if [ $# -ne 3 ]; then
    echo "$0: usage: sudo ./install.sh <mysql_root_pw> <mysql_backup_pw> <database_pw>"
    exit 1
fi

# Store passwords in variables
#
root_password=$1
backup_password=$2
database_password=$3
echo "All necessary passwords are provided and valid."
sleep 1

# Store database user in variable
#
database_user="mysql"

# Create and move to project directory
#
echo "Creating project directory"
sleep 1
mkdir -p /usr/local/zwischentoncloud/install || { echo "Failed to create project directory. Exiting." ; exit 1; }
cd /usr/local/zwischentoncloud/install || { echo "Failed to access project directory. Exiting." ; exit 1; }

# Install mysql if needed.
#
echo "Installing mysql-server..."
apt-get install mysql-server -y > /dev/null;

# Launch mysql on startup
#
systemctl enable mysql > /dev/null
systemctl start mysql > /dev/null
echo "Enabled and started mysql systemd service."
sleep 1

# Install mysql credential file
#
echo "Installing mysql credential file"
sleep 1
credentialsFile=/usr/local/zwischentoncloud/mysql.conf
cat << EOF > $credentialsFile
# zwischentoncloud database configuration file v1.0
# TOML 1.0.0-rc.2+

[client]
user = 'zwischentoncloud.backup'
password = '$backup_password'
host = 'localhost'
EOF

# Download and run mysql secure script
#
echo "Downloading database security script"
curl --progress-bar -L -o secure-mysql.sh https://raw.githubusercontent.com/TR-Projekt/zwischentoncloud/main/operation/secure-mysql.sh
chmod +x secure-mysql.sh
./secure-mysql.sh "$root_password"

# Download database creation scripts
#
echo "Downloading identity database creation script..."
curl --progress-bar -L -o create_identity_database.sql https://raw.githubusercontent.com/TR-Projekt/zwischentoncloud/main/database/create_identity_database.sql

echo "Downloading zwischenton cloud database creation script..."
curl --progress-bar -L -o create_zwischenton_database.sql https://raw.githubusercontent.com/TR-Projekt/zwischentoncloud/main/database/create_zwischenton_database.sql

# Run database creation script and configure users
#
echo "Configuring mysql"
sleep 1
echo "Creating databases..."
sleep 1
mysql -e "source /usr/local/zwischentoncloud/install/create_identity_database.sql"
mysql -e "source /usr/local/zwischentoncloud/install/create_zwischenton_database.sql"
echo "Creating local backup user..."
mysql -e "CREATE USER 'zwischentoncloud.backup'@'localhost' IDENTIFIED BY '$backup_password';"
mysql -e "GRANT ALL PRIVILEGES ON zwischenton_identity_database.* TO 'zwischentoncloud.backup'@'localhost';"
mysql -e "GRANT ALL PRIVILEGES ON zwischenton_cloud_database.* TO 'zwischentoncloud.backup'@'localhost';"
sleep 1
echo "Creating zwischenton identity database user..."
mysql -e "CREATE USER 'zwischenton.identity.writer'@'localhost' IDENTIFIED BY '$database_password';"
mysql -e "GRANT SELECT, INSERT, UPDATE, DELETE ON zwischenton_identity_database.* TO 'zwischenton.identity.writer'@'localhost';"
sleep 1
echo "Creating zwischentoncloud database user..."
mysql -e "CREATE USER 'zwischenton.api.writer'@'localhost' IDENTIFIED BY '$database_password';"
mysql -e "GRANT SELECT, INSERT, UPDATE, DELETE ON zwischenton_cloud_database.* TO 'zwischenton.api.writer'@'localhost';"
sleep 1
mysql -e "FLUSH PRIVILEGES;"

# Create the backup directory
#
echo "Create backup directory"
sleep 1
mkdir -p /srv/zwischentoncloud/backups || { echo "Failed to create backup directory. Exiting." ; exit 1; }
cd /srv/zwischentoncloud/backups || { echo "Failed to access backup directory. Exiting." ; exit 1; }

# Download the backup script
#
echo "Downloading database backup script"
curl --progress-bar -L -o backup.sh https://raw.githubusercontent.com/TR-Projekt/zwischentoncloud/main/operation/backup.sh
chmod +x /srv/zwischentoncloud/backups/backup.sh

# Installing a cronjob to run the backup every day at 3 pm.
#
echo "Installing a cronjob to periodically run a backup"
sleep 1
echo "0 3 * * * $database_user /srv/zwischentoncloud/backups/backup.sh" | sudo tee -a /etc/cron.d/zwischentoncloud_database_backup

## Set appropriate permissions
#
chown -R "$database_user":"$database_user" /usr/local/zwischentoncloud
chmod -R 761 /usr/local/zwischentoncloud
chown -R "$database_user":"$database_user" /srv/zwischentoncloud
chmod -R 761 /srv/zwischentoncloud
echo "Seting appropriate permissions..."
sleep 1

# Cleanup
#
echo "Cleanup"
cd /usr/local/zwischentoncloud || exit
rm -R /usr/local/zwischentoncloud/install
sleep 1

echo "Done."
sleep 1