package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSH configuration
const (
	sshHost     = "your_instance_ip"
	sshPort     = "22"
	sshUser     = "your_ssh_username"
	sshPassword = "your_ssh_password"
)

// Define a function to execute SSH commands
func sshExec(host, username, password, command string) {
	client, err := ssh.Dial("tcp", host+":"+sshPort, &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()

	// Execute the command
	output, err := session.CombinedOutput(command)
	if err != nil {
		log.Fatalf("Command failed: %v", err)
	}
	fmt.Printf("[%s] %s", host, strings.TrimSpace(string(output)))
}

func main() {
	var wg sync.WaitGroup

	// List of tasks to be executed on the instance
	tasks := []string{
		"sudo apt-get install -y vim curl git unzip",
		"sudo ufw allow 22/tcp",
		"sudo apt-get install -y mysql-server",
		"git clone https://github.com/rbenv/rbenv.git ~/.rbenv",
		`echo 'export PATH="$HOME/.rbenv/bin:$PATH"' >> ~/.bashrc`,
		`echo 'eval "$(rbenv init -)"' >> ~/.bashrc`,
	}

	// Update task (execute before the goroutines)
	cmd := "sudo apt-get update"
	cmdOutput, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Fatalf("Command failed: %v", err)
	}
	fmt.Printf("%s", cmdOutput)

	// Create a Goroutine for each task
	for _, task := range tasks {
		wg.Add(1)
		go func(task string) {
			defer wg.Done()
			sshExec(sshHost, sshUser, sshPassword, task)
		}(task)
	}

	// Wait for all Goroutines to finish
	wg.Wait()

	fmt.Println("Instance provisioning completed.")
}
