# Ruby-Golang-Provisioning
Devop-provisioning scripts in golang and ruby with docker-compose for rails apps and more. 

Bash: docker-compose up -> runs docker-compose.yaml and Dockerfile to build&setup an image
and more service as a rails-dev, staging, production and separate mysql-server ( At least the production server should be run with docker-swarm or kubernetes ). 
The other devop-scripts can be used via "docker exec .... ".
