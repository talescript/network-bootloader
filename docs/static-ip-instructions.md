smdSetting Up a Static IP Address on Ubuntu Server 22.04 LTS
Before you go any farther, it’s best to configure the PXE boot server with a fixed/static IP address. In this section, I am going to show you how to set up a static/fixed IP address on Ubuntu Server 22.04 LTS.

To configure a fixed/static IP address on Ubuntu Server 22.04 LTS , open the netplan configuration file /etc/netplan/00-installer-config.yaml with the nano text editor as follows:

$ sudo nano /etc/netplan/00-installer-config.yaml

To set a static/fixed IP address 192.168.0.130, change the configuration of the ens33 network interface as follows. Once you’re done, press <Ctrl> + X followed by Y and <Enter> to save the netplan configuration file.


To apply the changes, run the following command:

$ sudo netplan apply

A static/fixed IP address 192.168.0.130 should be set on the eno3 network interface, as you can see in the screenshot below. 

$ ip a


Creating the Required Directory Structure
In this section, I will create all the required directories for PXE booting (using the iPXE firmware) to work.

/pxeboot
config/
firmware/
os-images/

In the /pxeboot/firmware/ directory, we will store all the iPXE boot firmware files.

In the /pxeboot/os-images/ directory, we will create a separate subdirectory for each of the Linux distributions (that I want to PXE boot) and store the contents of the ISO images of these Linux distributions there. For example, for PXE booting Ubuntu Desktop 22.04 LTS, you can create a directory ubuntu-22.04-desktop-amd64/ in the /pxeboot/os-images/ directory and store the contents of the Ubuntu Desktop 22.04 LTS ISO image in that directory.

To create all the required directory structures, run the following command:

$ sudo mkdir -pv /pxeboot/{config,firmware,os-images}

All the required directory structures for PXE booting should be created.


Downloading iPXE Source Code and Compiling iPXE on Ubuntu 22.04 LTS
In this section, I am going to show you how to download the iPXE source code and compile it on Ubuntu 22.04 LTS so that we can use it for PXE booting.

First, update the APT package repository cache with the following command:

$ sudo apt update

To install the required build dependencies for iPXE, run the following command:

$ sudo apt install build-essential liblzma-dev isolinux git

Now, navigate to the ~/Downloads directory as follows:

$ cd ~/Downloads

Clone the iPXE GitHub repository on your Ubuntu 22.04 LTS machine as follows:

$ git clone https://github.com/ipxe/ipxe.git

Navigate to the ipxe/src/ directory as follows:

$ cd ipxe/src

You should see a lot of directories there containing the iPXE source code.

$ ls -lh


To configure iPXE to automatically boot from an iPXE boot script stored in the /pxeboot/config/ directory of your computer, you will need to create an iPXE boot script and embed it with the iPXE firmware when you compile it.

Create an iPXE boot script bootconfig.ipxe and open it with the nano text editor as follows:

Type in the following lines of codes in the bootconfig.ipxe file.




To compile iPXE BIOS and UEFI firmwares and embed the bootconfig.ipxe iPXE boot script in the compiled firmwares, run the following command:

$ make bin/ipxe.pxe bin/undionly.kpxe bin/undionly.kkpxe bin/undionly.kkkpxe bin-x86_64-efi/ipxe.efi EMBED=bootconfig.ipxe

Copying the Compiled iPXE Firmwares to /pxeboot/firmware Directory
Once the iPXE boot firmware files are compiled, copy them to the /pxeboot/firmware directory of your Ubuntu 22.04 LTS PXE boot server so that the PXE client computers can access them via TFTP.

$ sudo cp -v bin/{ipxe.pxe,undionly.kpxe,undionly.kkpxe,undionly.kkkpxe} bin-x86_64-efi/ipxe.efi /pxeboot/firmware/


Here, the iPXE boot firmware files ipxe.pxe, undionly.kpxe, undionly.kkpxe, and undionly.kkkpxe are for PXE booting on BIOS systems. The iPXE boot firmware file ipxe.efi is for PXE booting on UEFI systems.

Installing and Configuring a DHCP and TFTP Server on Ubuntu 22.04 LTS
For PXE boot to work, you will need a working DHCP and TFTP server running on your computer. There are many DHCP and TFTP server software. But, in this article, I will use dnsmasq. dnsmasq is mainly a DNS and DHCP server that can also be configured as a TFTP server.

On Ubuntu 22.04 LTS, dnsmasq is not installed by default. But it is available in the official package repository of Ubuntu 22.04, and you can install it with the APT package manager very easily.

To install dnsmasq on Ubuntu 22.04 LTS, run the following command:
$ sudo apt install dnsmasq -y

We will create a new dnsmasq configuration file. So, rename the original /etc/dnsmasq.conf file to /etc/dnsmasq.conf.backup as follows:

$ sudo mv -v /etc/dnsmasq.conf /etc/dnsmasq.conf.backup

Create an empty dnsmasq configuration file /etc/dnsmasq.conf with the following command:

$ sudo nano /etc/dnsmasq.conf


For the changes to take effect, restart the dnsmasq server as follows:
$ sudo systemctl restart dnsmasq

Installing and Configuring NFS Server on Ubuntu 22.04 LTS
Ubuntu Desktop 22.04 LTS uses casper to boot into Live installation mode. casper supports PXE boot via the NFS protocol only. Other Linux distributions like Fedora, CentOS/RHEL also support PXE booting via the NFS protocol. So, to boot Ubuntu Desktop 22.04 LTS and many other Linux distributions via PXE, you need to have a fully functional NFS server accessible over the network.

To install the NFS server on Ubuntu 22.04 LTS, run the following command:

$ sudo apt install nfs-kernel-server


Open the NFS server configuration file /etc/exports as follows:

$ sudo nano /etc/exports

To share the /pxeboot directory via NFS, add the following line at the end of the /etc/exports file:

/pxeboot           *(ro,sync,no_wdelay,insecure_locks,no_root_squash,insecure,no_subtree_check)
Once you’re done, press <Ctrl> + X followed by Y and <Enter> to save the NFS configuration file /etc/exports.



To make the new NFS share /pxeboot available, run the following command:

$ sudo exportfs -av

Configuring iPXE to PXE Boot Ubuntu Desktop 22.04 LTS Live Installer
In this section, I am going to show you how to configure iPXE on your Ubuntu 22.04 LTS PXE boot server to PXE boot Ubuntu Desktop 22.04 LTS Live installer on other computers (PXE clients).

NOTE: If you want to configure iPXE on your Ubuntu 22.04 LTS PXE boot server to PXE boot other Linux distributions, you will have to make the necessary changes. This shouldn’t be too hard.

First, navigate to the ~/Downloads directory of your Ubuntu 22.04 LTS PXE boot server as follows:

To download the Ubuntu Desktop 22.04 LTS ISO image from the official website of Ubuntu, run the following command:

$ wget https://releases.ubuntu.com/jammy/ubuntu-22.04-desktop-amd64.iso

Mount the Ubuntu Desktop 22.04 LTS ISO file ubuntu-22.04-desktop-amd64.iso in the /mnt directory as follows:

$ sudo mount -o loop ~/Downloads/ubuntu-22.04-desktop-amd64.iso /mnt

Create a dedicated directory ubuntu-22.04-desktop-amd64/ for storing the contents of the Ubuntu Desktop 22.04 LTS ISO image in the /pxeboot/os-images/ directory as follows:


$ sudo mkdir -pv /pxeboot/os-images/ubuntu-22.04-desktop-amd64

To copy the contents of the Ubuntu Desktop 22.04 LTS ISO image in the /pxeboot/os-images/ubuntu-22.04-desktop-amd64/ directory with rsync, run the following command:

$ sudo rsync -avz /mnt/ /pxeboot/os-images/ubuntu-22.04-desktop-amd64/

NOTE: If you don’t have rsync installed on Ubuntu 22.04 LTS and need any assistance in installing rsync on Ubuntu 22.04 LTS, read the article How to Use rsync Command to Copy Files on Ubuntu.

The contents of the Ubuntu Desktop 22.04 LTS ISO image are being copied to the /pxeboot/os-images/ubuntu-22.04-desktop-amd64/ directory. It will take a while to complete.

Unmount the Ubuntu Desktop 22.04 LTS ISO image from the /mnt directory as follows:

$ sudo umount /mnt

Now, create the default iPXE boot configuration file /pxeboot/config/boot.ipxe and open it with the nano text editor as follows:

$ sudo nano /pxeboot/config/boot.ipxe

Type in the following lines in the iPXE boot configuration file /pxeboot/config/boot.ipxe:



Here, server_ip is the IP address of the Ubuntu 22.04 LTS PXE boot server¹, and root_path is the NFS share path².

ubuntu-22.04-desktop-amd64 is the label for the boot menu entry Install Ubuntu Desktop 22.04 LTS, and the boot codes for PXE booting Ubuntu Desktop 22.04 LTS are also labeled with the same name³.

os_root is the name of the subdirectory in the /pxeboot/os-images/ directory where you’ve copied the contents of the Ubuntu Desktop 22.04 LTS ISO image⁴.