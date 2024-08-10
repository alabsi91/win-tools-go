#!/bin/bash

go build -o build/win-tools.exe main.go

"C:/Program Files (x86)/NSIS/makensis.exe" installer/installer.nsi

# Variables
INSTALLER_PATH="build/installer.exe"
REPO="alabsi91/win-tools-go"
CHOICES_CREATE="create a new release"
CHOICES_UPLOAD="upload assets to a release"

# Prompt for the release tag
read -p "Enter the release tag [default: v1.0.0]: " VERSION
VERSION=${VERSION:-v1.0.0}

# Prompt for the operation
echo "Choose an operation:"
echo "1) ${CHOICES_CREATE}"
echo "2) ${CHOICES_UPLOAD}"
read -p "Enter choice (1 or 2): " OPERATION_CHOICE

case $OPERATION_CHOICE in
    1)
        # Create a new release
        gh release create "$VERSION" -R "$REPO" --notes "$VERSION" --latest --title "$VERSION" "$INSTALLER_PATH"
    ;;
    2)
        # Upload assets to a release
        gh release upload "$VERSION" -R "$REPO" "$INSTALLER_PATH" --clobber
    ;;
    *)
        echo "Invalid choice. Exiting."
        exit 1
    ;;
esac

exit 0