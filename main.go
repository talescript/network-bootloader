package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {

	url := "http://releases.ubuntu.com/20.04/ubuntu-20.04-desktop-amd64.iso"
	fmt.Println("Downloading Ubuntu ISO...")

	cmd := exec.Command("wget", "-O", "ubuntu.iso", url)

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error downloading ISO:", err)
		return
	}

	fmt.Println("Creating bootable USB drive...")
	cmd = exec.Command("dd", "if=ubuntu.iso", "of=/dev/sdb", "bs=4M", "status=progress")

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error creating bootable USB:", err)
		return
	}

	fmt.Println("Rebooting system and booting from USB...")
	cmd = exec.Command("reboot")

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error rebooting system:", err)
		return
	}

	fmt.Println("Waiting for installation to complete...")
	fmt.Println("Press Enter when installation is complete")

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	fmt.Println("Rebooting system and booting into installed operating system...")
	cmd = exec.Command("reboot")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error rebooting system:", err)
		return
	}
}
