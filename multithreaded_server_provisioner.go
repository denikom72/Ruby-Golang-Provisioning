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

// Define provisioning tasks
var tasks = []string{
	"sudo apt-get update",
	"sudo apt-get install -y vim curl git unzip",
	"sudo ufw allow 22/tcp",
	"sudo apt-get install -y mysql-server",
	"git clone https://github.com/rbenv/rbenv.git ~/.rbenv",
	`echo 'export PATH="$HOME/.rbenv/bin:$PATH"' >> ~/.bashrc`,
	`echo 'eval "$(rbenv init -)"' >> ~/.bashrc`,
}

func main() {
	var wg sync.WaitGroup

	// Create an SSH client configuration
	config := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{ssh.Password(sshPassword)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For demo purposes; should be used securely in production.
	}

	// Start a Goroutine for each task
	for _, task := range tasks {
		wg.Add(1)
		go func(task string) {
			defer wg.Done()

			client, err := ssh.Dial("tcp", sshHost+":"+sshPort, config)
			if err != nil {
				log.Printf("Failed to dial: %v", err)
				return
			}
			defer client.Close()

			session, err := client.NewSession()
			if err != nil {
				log.Printf("Failed to create session: %v", err)
				return
			}
			defer session.Close()

			// Execute the task
			output, err := session.CombinedOutput(task)
			if err != nil {
				log.Printf("Task failed: %v", err)
			} else {
				log.Printf("Task succeeded: %s", strings.TrimSpace(string(output)))
			}
		}(task)
	}

	// Wait for all Goroutines to finish
	wg.Wait()

	fmt.Println("Instance provisioning completed.")
}
