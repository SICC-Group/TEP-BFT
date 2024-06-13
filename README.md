The Chinese version of README.md can be found at [README_zh.md](https://github.com/SICC-Group/TEP-BFT/blob/main/README_zh.md)

The operating system used is Ubuntu 22.04.
# Compile and install
The dependencies that need to be installed and the process of compiling the official source code can be referred to the FISCO official compilation documentation [5. Node source code compilation - FISCO BCOS 3.0 v3.6.0 documentation](https://fisco-bcos-doc.readthedocs.io/zh-cn/latest/docs/tutorial/compile_binary.html)
（We recommend installing v3.2.0）
## Clone code
Clone our optimized code:
```bash
# It is recommended to create another compilation directory, separate from the official source code
mkdir -p ~/TEPBFT && cd ~/TEPBFT

# Clone our optimized code
https://github.com/SICC-Group/TEP-BFT.gitt

cd FISCO-BCOS
```
## Compile
```bash
cd ~/TEPBFT/FISCO-BCOS

# Create compilation directory
mkdir -p build && cd build
cmake .. -DCMAKE_BUILD_TYPE=Release -DBUILD_STATIC=ON

# Compile source code (high-performance machines can be added -j4 uses 4-core accelerated compilation)
make -j4
# If the above command is stuck for too long, you can use 'make clean' command to delete related files and use the following command instead:
make

rm -rf fisco-bcos-tars-service/*.tgz && make tar
```
The compiled binary is located in the path **FISCO-BCOS/build/fisco-bcos-air/**.
# Deploy
## Install the SGX driver
```bash
# Update the apt source and install the software
sudo apt update
sudo apt-get install -y gcc git make vim curl

# Clone the linux-sgx-driver repository and compile it
git clone https://github.com/intel/linux-sgx-driver
cd linux-sgx-driver/
make

# Install the SGX driver
sudo mkdir -p "/lib/modules/$(uname -r)/kernel/drivers/intel/sgx"
sudo cp isgx.ko "/lib/modules/$(uname -r)/kernel/drivers/intel/sgx"
sudo sh -c "grep -qxF 'isgx' /etc/modules || echo 'isgx' >> /etc/modules"
sudo depmod
sudo modprobe isgx
lsmod | grep isgx
```
## Install SGX PSW
```bash
cd ~
# Add SGX repository and install dependencies
echo 'deb [arch=amd64] https://download.01.org/intel-sgx/sgx_repo/ubuntu jammy main' | sudo tee /etc/apt/sources.list.d/intel-sgx.list
wget "https://download.01.org/intel-sgx/sgx_repo/ubuntu/intel-sgx-deb.key"
sudo apt-key add intel-sgx-deb.key
sudo apt update
sudo apt-get install -y libssl-dev libcurl4-openssl-dev libprotobuf-dev

# Install SGX-related libraries
sudo apt-get install -y libsgx-launch libsgx-urts libsgx-epid libsgx-quote-ex libsgx-dcap-ql
```
## Install the Gramine package
```bash
# Add gramine repository and install
sudo curl -fsSLo /usr/share/keyrings/gramine-keyring.gpg https://packages.gramineproject.io/gramine-keyring.gpg
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/gramine-keyring.gpg] https://packages.gramineproject.io/ $(lsb_release -sc) main" | sudo tee /etc/apt/sources.list.d/gramine.list
sudo curl -fsSLo /usr/share/keyrings/intel-sgx-deb.asc https://download.01.org/intel-sgx/sgx_repo/ubuntu/intel-sgx-deb.key
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/intel-sgx-deb.asc] https://download.01.org/intel-sgx/sgx_repo/ubuntu jammy main" | sudo tee /etc/apt/sources.list.d/intel-sgx.list
sudo apt-get update
sudo apt-get install -y gramine

# Generate the SGX private key
sudo gramine-sgx-gen-private-key

# Install the 32-bit library
sudo apt-get install -y lib32stdc++6 lib32z1
```
## Build a network
```bash
cd TEPBFT
curl -#LO https://osp-1257653870.cos.ap-guangzhou.myqcloud.com/FISCO-BCOS/FISCO-BCOS/releases/v3.2.0/build_chain.sh

# This line of code uses the build script to set up multiple nodes, and the ip address is changed according to the private ip address of your machine
# Just run this command on one machine
# If you want to reestablish a network after it has been established, delete folder nodes first
bash build_chain.sh -p 30300,20200 -l [private network IP]:1,[private network IP]:1,[private network IP]:1,[private network IP]:1


# Create the Makefile ISCO -bcos.manifest.template and fill in the contents
# In the "nodes/ [private network IP]/node0/" path for each machine
cd ~/fisco/nodes/[private network IP]/node0/
vim Makefile
# Fill the contents of the file as last
vim fisco-bcos.manifest.template
# Fill the contents of the file as last

# Copy the compiled fisco-bcos binary executable, Makefile, and fisco-bcos.manifest.template to the node folder of each server or physical machine
cp fisco-bcos Makefile fisco-bcos.manifest.template ~/fisco/nodes/[private network IP]/node0/

# Start all servers or physical machine and create corresponding folders
mkdir -p ~/fisco/nodes

# Copy the folders to the servers of the respective nodes using the public IP address
scp -r ~/fisco/nodes/[private network IP]/ root@[public IP address]:~/fisco/nodes


# Modify the openssl file
sudo vim /etc/ssl/openssl.cnf
# Comment out the following line
# providers = provider_sect

# Compile
cd ~/fisco/nodes/[private network IP]/node0/
make SGX=1
```
# Startup
```bash
# Under ~/fisco/nodes/[private IP]/node0/
# Each machine is started like this, direct command line input
sudo gramine-sgx fisco-bcos -c config.ini -g config.genesis
```
# Stress test
Refer to FISCO's official stress testing documentation [9. Stress Testing Guide](https://fisco-bcos-doc.readthedocs.io/zh-cn/latest/docs/operation_and_maintenance/stress_testing.html)
# Configuration files
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
The two parameters`sgx.enclave_size`and`sys.stack.size`can be properly debugged based on machine performance.

