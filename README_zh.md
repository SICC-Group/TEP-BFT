使用的环境是Ubuntu 22.04.
# 编译安装
需要安装的依赖以及编译官方源码的过程可以参考FISCO官方编译文档[5. 节点源码编译 — FISCO BCOS 3.0 v3.6.0 文档](https://fisco-bcos-doc.readthedocs.io/zh-cn/latest/docs/tutorial/compile_binary.html)
（建议安装v3.2.0版本）
## 克隆代码
克隆我们的优化代码：
```bash
# 建议创建另一个编译目录，区别于官方源码
mkdir -p ~/TEPBFT && cd ~/TEPBFT

# 克隆我们的优化代码
https://github.com/SICC-Group/TEP-BFT.gitt

# 切换到源码目录
cd FISCO-BCOS
```
## 编译
```bash
# 进入源码目录
cd ~/TEPBFT/FISCO-BCOS

# 创建编译目录
mkdir -p build && cd build
cmake .. -DCMAKE_BUILD_TYPE=Release -DBUILD_STATIC=ON

# 编译源码(高性能机器可添加-j4使用4核加速编译)
make -j4
# 如果上面这个命令卡的时间太长，可以make clean后改用下面这个命令：
# make

# 生成tgz包
rm -rf fisco-bcos-tars-service/*.tgz && make tar
```
编译完的二进制文件位于
**FISCO-BCOS/build/fisco-bcos-air/**
路径下。
# 部署
## 安装SGX驱动程序
```bash
# 更新apt源并安装软件
sudo apt update
sudo apt-get install -y gcc git make vim curl

# 克隆linux-sgx-driver仓库并编译
git clone https://github.com/intel/linux-sgx-driver
cd linux-sgx-driver/
make

# 安装驱动
sudo mkdir -p "/lib/modules/$(uname -r)/kernel/drivers/intel/sgx"
sudo cp isgx.ko "/lib/modules/$(uname -r)/kernel/drivers/intel/sgx"
sudo sh -c "grep -qxF 'isgx' /etc/modules || echo 'isgx' >> /etc/modules"
sudo depmod
sudo modprobe isgx
lsmod | grep isgx
```
## 安装SGX PSW
```bash
cd ~
# 添加SGX仓库并安装依赖
echo 'deb [arch=amd64] https://download.01.org/intel-sgx/sgx_repo/ubuntu jammy main' | sudo tee /etc/apt/sources.list.d/intel-sgx.list
wget "https://download.01.org/intel-sgx/sgx_repo/ubuntu/intel-sgx-deb.key"
sudo apt-key add intel-sgx-deb.key
sudo apt update
sudo apt-get install -y libssl-dev libcurl4-openssl-dev libprotobuf-dev

# 安装SGX相关库
sudo apt-get install -y libsgx-launch libsgx-urts libsgx-epid libsgx-quote-ex libsgx-dcap-ql
```
## 安装Gramine包
```bash
# 添加gramine仓库并安装
sudo curl -fsSLo /usr/share/keyrings/gramine-keyring.gpg https://packages.gramineproject.io/gramine-keyring.gpg
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/gramine-keyring.gpg] https://packages.gramineproject.io/ $(lsb_release -sc) main" | sudo tee /etc/apt/sources.list.d/gramine.list
sudo curl -fsSLo /usr/share/keyrings/intel-sgx-deb.asc https://download.01.org/intel-sgx/sgx_repo/ubuntu/intel-sgx-deb.key
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/intel-sgx-deb.asc] https://download.01.org/intel-sgx/sgx_repo/ubuntu jammy main" | sudo tee /etc/apt/sources.list.d/intel-sgx.list
sudo apt-get update
sudo apt-get install -y gramine

# 生成SGX私钥
sudo gramine-sgx-gen-private-key

# 安装32位库
sudo apt-get install -y lib32stdc++6 lib32z1
```
## 搭建网络
```bash
cd TEPBFT
curl -#LO https://osp-1257653870.cos.ap-guangzhou.myqcloud.com/FISCO-BCOS/FISCO-BCOS/releases/v3.2.0/build_chain.sh

# 这代码作用是通过build脚本搭建多机节点,ip根据服务器私网ip修改
# 只在一台服务器中运行这个命令就行
# 已经搭建过需要先删除nodes文件夹，再搭建
bash build_chain.sh -p 30300,20200 -l [私网IP]:1,[私网IP]:1,[私网IP]:1,[私网IP]:1


# 创建Makefile fisco-bcos.manifest.template，填充内容
# 位置应该在每个服务器的nodes/服务器私网IP/node0/下
cd ~/fisco/nodes/私网IP/node0/
vim Makefile
# 填充文件内容如最后
vim fisco-bcos.manifest.template
# 填充文件内容如最后

# 将编译得到的fisco-bcos二进制可执行程序、Makefile和fisco-bcos.manifest.template复制到每个服务器或物理机的节点文件夹中
cp fisco-bcos Makefile fisco-bcos.manifest.template ~/fisco/nodes/[私网IP]/node0/

# 启动所有的服务器或物理机，创建对应文件夹
mkdir -p ~/fisco/nodes

# 使用公网IP把文件夹分别复制到各自节点的服务器上
# 注意公私网IP地址要对应好，公网IP会更换
scp -r ~/fisco/nodes/[私网IP]/ root@[公网IP]:~/fisco/nodes


# 修改openssl文件
sudo vim /etc/ssl/openssl.cnf
# 注释掉
# providers = provider_sect

# 编译
cd ~/fisco/nodes/[私网IP]/node0/
make SGX=1
```
# 运行
```bash
# 在~/fisco/nodes/[私网IP]/node0/下
# 每台机器这样启动，直接命令行输入
sudo gramine-sgx fisco-bcos -c config.ini -g config.genesis
```
# 压测
参考FISCO官方的压测文件[9. 压力测试指南](https://fisco-bcos-doc.readthedocs.io/zh-cn/latest/docs/operation_and_maintenance/stress_testing.html)
# 配置文件
## Makefile
```bash
# none/error/warning/debug/trace/all
GRAMINE_LOG_LEVEL = error
# directory with arch-specific libraries, used by Redis
# the below path works for Debian/Ubuntu; for CentOS/RHEL/Fedora, you should
# overwrite this default like this: `ARCH_LIBDIR=/lib64 make`
ARCH_LIBDIR ?= /lib/$(shell $(CC) -dumpmachine)

.PHONY: all
all: fisco-bcos.manifest.sgx fisco-bcos.sig

fisco-bcos.manifest: fisco-bcos.manifest.template
		gramine-manifest \
				-Dlog_level=$(GRAMINE_LOG_LEVEL) \
				-Darch_libdir=$(ARCH_LIBDIR) \
				$< >$@

fisco-bcos.sig fisco-bcos.manifest.sgx: sgx_sign
		@:

.INTERMEDIATE: sgx_sign
sgx_sign: fisco-bcos.manifest fisco-bcos
		gramine-sgx-sign \
				--manifest $< \
				--output $<.sgx

.PHONY: clean
clean:
		$(RM) *.token *.sig *.manifest.sgx *.manifest

.PHONY: distclean
distclean: clean
```
## fisco-bcos.manifest.template 
```bash
# Redis manifest file example

################################## GRAMINE ####################################

loader.entrypoint = "file:{{ gramine.libos }}"

libos.entrypoint = "/fisco-bcos"

loader.log_level = "{{ log_level }}"

#loader.pal_internal_mem_size = "1G"

loader.insecure__use_cmdline_argv = true

################################# ENV VARS ####################################

loader.env.LD_LIBRARY_PATH = "/lib:{{ arch_libdir }}:/usr/{{ arch_libdir }}"

################################## SIGNALS ####################################

sys.enable_sigterm_injection = true

################################# MOUNT FS ###################################

fs.mounts = [
	{ path = "/lib", uri = "file:{{ gramine.runtimedir() }}" },

	{ path = "{{ arch_libdir }}", uri = "file:{{ arch_libdir }}" },
  
	{ path = "/usr/{{ arch_libdir }}", uri = "file:/usr/{{ arch_libdir }}" },
  
	{ path = "/etc", uri = "file:/etc" },
  
	{ path = "/fisco-bcos", uri = "file:fisco-bcos" },
]

############################### SGX: GENERAL ##################################

sgx.debug = true

sgx.enclave_size = "8192M"

sgx.max_threads = 1024

#sgx.rpc_thread_num = 64

#sgx.enable_status = true

#sgx.profile.with_stack = true

#sgx.profile.enable = "all"

#sgx.profile.mode = "ocall_outer"

################################## Thread Size
####################################

sys.stack.size = "218K"

sgx.file_check_policy = "allow_all_but_log"

sgx.nonpie_binary = true

############################# SGX: TRUSTED FILES ###############################

sgx.trusted_files = [
	"file:{{ gramine.libos }}",
	"file:fisco-bcos",
	"file:{{ gramine.runtimedir() }}/",
	"file:{{ arch_libdir }}/",
	"file:/usr/{{ arch_libdir }}/",
]

############################# SGX: ALLOWED FILES ###############################

sgx.allowed_files = [
	# Name Service Switch (NSS) files. Glibc reads these files as part of name-
	# service information gathering. For more info, see 'man nsswitch.conf'.
	"file:/etc/nsswitch.conf",
	"file:/etc/ethers",
	"file:/etc/hosts",
	"file:/etc/group",
	"file:/etc/passwd",
	# getaddrinfo(3) configuration file. Glibc reads this file to correctly find
	# network addresses. For more info, see 'man gai.conf'.
	"file:/etc/gai.conf",
]
```
sgx.enclave_size和sys.stack.size这两个参数可以根据机器性能适当调试。