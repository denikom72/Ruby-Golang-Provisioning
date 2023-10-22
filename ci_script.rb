#!/usr/bin/env ruby

require 'rake'

task :test do
  # Run unit and integration tests
  sh 'rspec'
end

task :build do
  # Build the application
  sh 'docker build -t my_app .'
end

task :deploy => [:test, :build] do
  # Deploy to a staging environment
  sh 'kubectl apply -f deployment.yaml'
end

task :cleanup do
  # Clean up temporary files, test artifacts, etc.
  sh 'rm -rf tmp/*'
end

Rake::Task[:deploy].invoke
Rake::Task[:cleanup].invoke

puts "Continuous integration pipeline completed successfully."
