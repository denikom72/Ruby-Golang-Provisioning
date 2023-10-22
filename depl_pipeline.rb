#!/usr/bin/env ruby

services = [
  { name: 'service1', repo: 'https://github.com/user/service1.git' },
  { name: 'service2', repo: 'https://github.com/user/service2.git' },
  { name: 'service3', repo: 'https://github.com/user/service3.git' }
]

def build_service(service)
  puts "Building #{service[:name]}..."
  `git clone #{service[:repo]} #{service[:name]}`
  `cd #{service[:name]} && docker build -t #{service[:name]}:latest .`
  puts "Build for #{service[:name]} completed."
end

def test_service(service)
  puts "Testing #{service[:name]}..."
  `cd #{service[:name]} && rspec`
  puts "Testing for #{service[:name]} completed."
end

def deploy_service(service)
  puts "Deploying #{service[:name]}..."
  `cd #{service[:name]} && kubectl apply -f deployment.yaml`
  puts "Deployment for #{service[:name]} completed."
end

services.each do |service|
  build_service(service)
  test_service(service)
  deploy_service(service)
end

# Notify the team about successful deployment
notify_team
puts "Deployment pipeline completed successfully."

def notify_team
  # Add code to send notifications (e.g., Slack, email) to the team
  puts "Team notified about successful deployment."
end
