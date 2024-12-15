TARGET_PATH=$(pwd)/bin
# build file name
OS_NAME=$(uname | tr '[A-Z]' '[a-z]')
ARCH_NAME=$(arch  | tr '[A-Z]' '[a-z]')
if [ $ARCH_NAME == "x86_64" ]; then
  ARCH_NAME="x64"
fi
echo "Building go-trans for $OS_NAME $ARCH_NAME ..."
FILE_NAME="go-trans-$OS_NAME-$ARCH_NAME"
FILE_PATH=$TARGET_PATH/$FILE_NAME
# build target
go build -o $FILE_NAME main.go
mv ./$FILE_NAME $FILE_PATH
echo "Build target success!, file path: $FILE_PATH"