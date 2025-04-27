#!/bin/bash
#
# install.sh - ZwischentonCloud Update Script
#
# (c)2020-2025 Simon Gaus
#

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
# 📦 Install ZwischentonCloud
# ─────────────────────────────────────────────────────────────────────────────
mv zwischentoncloud /usr/local/bin/zwischentoncloud || {
    echo -e "\n🚨  ERROR: Failed to install Zwischenton Cloud binary. Exiting.\n"
    exit 1
}
echo -e "✅  Updated ZwischentonCloud at \e[1;34m/usr/local/bin/zwischentoncloud\e[0m."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🎉 Restart ZwischentonCloud
# ─────────────────────────────────────────────────────────────────────────────
systemctl restart zwischentoncloud
echo -e "✅  Restarted ZwischentonCloud\e[0m."

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
echo -e "\033[1;32m✅  UPDATE COMPLETE! 🚀\033[0m"
echo -e "\033[1;32m══════════════════════════════════════════════════════════════════════════\033[0m"