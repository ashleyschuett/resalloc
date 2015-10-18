#!/bin/bash
# This file needs to be run at the startup of any slave containers
# It ensures that the server can be communicated with properly
# The resalloc program.
sudo apt-key adv --keyserver hkp://pgp.mit.edu:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
sudo touch /etc/apt/sources.list.d/docker.list
sudo bash -c 'echo "deb https://apt.dockerproject.org/repo ubuntu-trusty main" > /etc/apt/sources.list.d/docker.list'
sudo apt-get update
sudo apt-get -y install docker-engine
sudo bash -c 'echo "DOCKER_OPTS=\"-H 0.0.0.0:5555\"" >> /etc/default/docker'
sudo service docker restart
