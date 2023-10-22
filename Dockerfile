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

# Expose port 3000 for the Rails application
EXPOSE 3000

# Command to start the Rails application
CMD ["bundle", "exec", "rails", "server", "-b", "0.0.0.0"]
