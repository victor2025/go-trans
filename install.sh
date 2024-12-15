BIN_PATH="./bin"
TARGET_PATH="/usr/local/lib/go-trans"
USER_BIN_PATH="/usr/local/bin/go-trans"
# build file name
OS_NAME=$(uname | tr '[A-Z]' '[a-z]')
ARCH_NAME=$(arch  | tr '[A-Z]' '[a-z]')
if [ $ARCH_NAME == "x86_64" ]; then
  ARCH_NAME="x64"
fi
echo "Installing go-trans for $OS_NAME $ARCH_NAME ..."
FILE_NAME="go-trans-$OS_NAME-$ARCH_NAME"
FILE_PATH=$TARGET_PATH/$FILE_NAME
# mv executable file
sudo mkdir $TARGET_PATH
sudo cp $BIN_PATH/$FILE_NAME $TARGET_PATH
sudo rm $USER_BIN_PATH
sudo ln $TARGET_PATH/$FILE_NAME $USER_BIN_PATH
echo "Install go-trans success!"