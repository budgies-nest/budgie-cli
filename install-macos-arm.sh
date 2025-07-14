#!/bin/bash
set -o allexport; source release.env; set +o allexport

curl -L -o budgie-macos-arm64 https://github.com/budgies-nest/budgie-cli/releases/download/${TAG}/budgie-macos-arm64

chmod +x budgie-macos-arm64
sudo mv -f budgie-macos-arm64 /usr/local/bin/budgie


