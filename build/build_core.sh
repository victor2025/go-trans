MAKE_FILE=Makefile_core
export SRC_PATH=$(pwd)/../core_main.go
export TARGET_PATH=$(pwd)/../bin/core
export TARGET_NAME=libgotrans
export TARGET_SUFFIX=.so

# build for linux
export OS=linux
# amd64
export ARCH=amd64
export BUILD_CC=gcc
make -f $MAKE_FILE
# arm64
export ARCH=arm64
export BUILD_CC=aarch64-linux-gnu-gcc
make -f $MAKE_FILE


# build for android
export SYSROOT=$ANDROID_NDK_HOME/toolchains/llvm/prebuilt/linux-x86_64/sysroot
export NDK_CC_ROOT=$ANDROID_NDK_HOME/toolchains/llvm/prebuilt/linux-x86_64/bin
export CFLAGS="--sysroot=$SYSROOT -I$SYSROOT/usr/include"
export LDFLAGS="--sysroot=$SYSROOT -L$SYSROOT/usr/lib"
export OS=android
# arm64
export ARCH=arm64
export BUILD_CC=$NDK_CC_ROOT/aarch64-linux-android35-clang
make -f $MAKE_FILE
# arm
export ARCH=arm
export BUILD_CC=$NDK_CC_ROOT/armv7a-linux-androideabi35-clang
make -f $MAKE_FILE

# build for windows
export OS=windows
export TARGET_SUFFIX=.dll
export ARCH=amd64
export BUILD_CC=x86_64-w64-mingw32-gcc
make -f $MAKE_FILE

# build for macos
#export OS=darwin
#export ARCH=amd64
#export BUILD_CC=clang
#make -f $MAKE_FILE

