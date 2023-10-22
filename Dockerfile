# Use the official Ruby image as a base image
FROM ruby:2.7

# Set the working directory in the container
WORKDIR /app

# Copy the Gemfile and Gemfile.lock into the container
COPY Gemfile Gemfile.lock ./

# Install project dependencies
RUN gem install bundler && bundle install

# Copy the rest of the application code into the container
COPY . .

# Install Golang
RUN apt-get update && apt-get install -y golang

# Copy your Golang provisioner script into the container
COPY provisioner.go /app/provisioner.go

# Expose port 3000 for the Rails application
EXPOSE 3000

# Set the provisioner script as the entry point
CMD ["go", "run", "/app/provisioner.go"]
