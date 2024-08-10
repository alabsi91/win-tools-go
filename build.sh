#!/bin/bash

go build -o build/win-tools.exe main.go

"C:/Program Files (x86)/NSIS/makensis.exe" installer/installer.nsi