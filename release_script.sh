
VERSION="0.1.0"


#=====================BUILDING(STATIC LINKING)=====================
echo "[INFO] building..."
# GNU/Linux build
CGO_ENABLED=1 CC=gcc GOOS=linux GOARCH=amd64 go build -tags static -ldflags "-s -w" -v -o ./release/tanks-$VERSION-GNU-Linux-amd64
# Windows build
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -tags static -ldflags "-s -w" -v -o ./release/tanks-$VERSION-windows-amd64.exe


#=====================ADDING RESOURCES, AND COMPRESSING INTO AN ARCHIEVE=====================
echo "[INFO] copying resources and creating archieve..."
cp -r resources ./release

echo "[INFO] Entering release dir, and creating archieve..."
cd ./release
zip -r tanks-$VERSION-GNU-Linux-amd64.zip resources tanks-$VERSION-GNU-Linux-amd64
zip -r tanks-$VERSION-windows-amd64.zip resources tanks-$VERSION-windows-amd64.exe


#=====================CLEANUP=====================
echo "[INFO] cleanup..."
rm -r resources
rm tanks-$VERSION-GNU-Linux-amd64
rm tanks-$VERSION-windows-amd64.exe

echo "[INFO] done"
