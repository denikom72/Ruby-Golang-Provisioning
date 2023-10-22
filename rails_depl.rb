#!/usr/bin/env ruby

require 'net/ssh'

server = 'your_server'
user = 'deploy'
deploy_to = '/var/www/your_app'

Net::SSH.start(server, user) do |ssh|
  # Pull the latest code from your repository
  ssh.exec!("cd #{deploy_to} && git pull origin master")

  # Install/update gems
  ssh.exec!("cd #{deploy_to} && bundle install")

  # Precompile assets
  ssh.exec!("cd #{deploy_to} && RAILS_ENV=production bundle exec rake assets:precompile")

  # Run database migrations
  ssh.exec!("cd #{deploy_to} && RAILS_ENV=production bundle exec rake db:migrate")

  # Restart the application server
  ssh.exec!("sudo systemctl restart your_app.service")
end
