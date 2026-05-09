# llm-local
A framework for local agentic llm usage.

Notes:

sudo apt update && sudo apt upgrade -y
sudo apt install -y build-essential cmake git libcurl4-openssl-dev libopenblas-dev

wget https://developer.download.nvidia.com/compute/cuda/repos/wsl-ubuntu/x86_64/cuda-keyring_1.1-1_all.deb
sudo dpkg -i cuda-keyring_1.1-1_all.deb
sudo apt-get update
sudo apt-get -y install cuda-toolkit-13-1 //currently 13.2 is not working with qwen