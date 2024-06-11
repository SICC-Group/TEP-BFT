使用说明：
1. 按照官方使用ubuntu22.04为环境下载依赖
sudo apt update
sudo apt install -y wget python3-dev git curl zip unzip tar
sudo apt install -y --no-install-recommends clang make build-essential cmake libssl-dev zlib1g-dev ca-certificates libgmp-dev flex bison patch libzstd-dev ninja-build pkg-config


# 安装rust
curl https://sh.rustup.rs -sSf | bash -s -- -y
source $HOME/.cargo/env

2.克隆代码
# 创建源码编译目录
mkdir -p ~/fisco && cd ~/fisco

# 克隆代码
git clone 后续添加

# 若因为网络问题导致长时间无法执行上面的命令，请尝试下面的命令：
git clone 后续添加仓库

# 切换到源码目录
cd FISCO-BCOS


3. 编译
# 进入源码目录
cd ~/fisco/FISCO-BCOS

# 创建编译目录
mkdir -p build && cd build
cmake .. -DCMAKE_BUILD_TYPE=Release -DBUILD_STATIC=ON(使用这一项)

# 要挂个vpn

# 编译源码(高性能机器可添加-j4使用4核加速编译)
make -j4

4. 使用
在./build/fisco-bcos-air/fisco-bcos,使用它。拷贝到运行区块链的127.0.0.1目录里面，进行替换。
