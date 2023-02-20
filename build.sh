# Read version
version=$(< VERSION)

# Output folder
mkdir -p build
cd build

# AMD64 Windows
GOOS=windows GOARCH=amd64 go build -o bctext.exe -ldflags="-s -w" ../
zip -m -9 -q bctext-${version}-win64.zip ./bctext.exe

# AMD64 Linux
GOOS=linux GOARCH=amd64 go build -o bctext -ldflags="-s -w" ../
upx --brute bctext
tar -zcvf bctext-${version}-x86_64-linux.tar.gz ./bctext
rm bctext

# ARM64 Linux
GOOS=linux GOARCH=arm64 go build -o bctext -ldflags="-s -w" ../
upx --brute bctext
tar -zcvf bctext-${version}-arm64-linux.tar.gz ./bctext
rm bctext

# Get SHA256 hash of files
shasum * > SHA256SUMS

# Back to parent
cd ..
