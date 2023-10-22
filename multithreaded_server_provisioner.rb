#!/usr/bin/env ruby

require 'open3'
require 'net/ssh'
require 'thread'

# Define a function to execute system commands
def run_command(command)
  stdout, stderr, status = Open3.capture3(command)
  puts stdout
  raise "Command failed: #{stderr}" unless status.success?
end

# Define a function to SSH into the instance and execute commands
def ssh_exec(host, username, password, commands)
  Net::SSH.start(host, username, password: password) do |ssh|
    commands.each do |cmd|
      ssh.exec!(cmd) do |channel, stream, data|
        puts "[#{host}] #{data}"
      end
    end
  end
end

# Configuration
instance_ip = 'your_instance_ip'
ssh_username = 'your_ssh_username'
ssh_password = 'your_ssh_password'

# List of tasks to be executed on the instance
tasks = [
  "sudo apt-get update",
  "sudo apt-get install -y vim curl git unzip",
  "sudo ufw allow 22/tcp",
  "sudo apt-get install -y mysql-server",
  "git clone https://github.com/rbenv/rbenv.git ~/.rbenv",
  "echo 'export PATH=\"$HOME/.rbenv/bin:$PATH\"' >> ~/.bashrc",
  "echo 'eval \"$(rbenv init -)\"' >> ~/.bashrc"
]

# Create an array of threads to execute tasks concurrently
threads = tasks.map do |task|
  Thread.new do
    ssh_exec(instance_ip, ssh_username, ssh_password, [task])
  end
end

# Start the threads and wait for them to finish
threads.each(&:join)

puts "Instance provisioning completed."
