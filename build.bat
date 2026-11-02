@echo off
set GOARCH=amd64
set GOOS=windows
set CGO_ENABLED=0
del .\\deploy\\win64\\liansu-server.exe
%GOROOT%\\bin\\go.exe build -ldflags="-H windowsgui" -o .\\deploy\\win64\\liansu-server.exe .\\src
@echo on