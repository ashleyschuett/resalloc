# resalloc

Resalloc is a resource allocation and resource allocation management application. It uses a master/slave setup where the master keeps a record of each slave and its location on the network. The master does not lease any of its own resources so at least one slave is required in order to start creating leases. Slave creation requires that docker is setup on the machine and the resources/setup.sh script contains all commands that must be run to properly configure a slave server. After the slave has been properly setup you can add it to the master by issuing the `python client.py machine create <machine_name> <username> <ip>` command.

In order to run any commands against the master node, you must first be registered and logged in which will provide you with an OAuth token on login. Resalloc allows you to create "resources" which define the type of image you can start with the "lease" command. Resources are defined using native [Dockerfile](https://docs.docker.com/reference/builder/) syntax the most basic being "FROM busybox". This allows you to lease a [busybox](https://hub.docker.com/_/busybox/) image. Something more familiar may be ubuntu, which you could create a resource for using "FROM ubuntu \nRUN apt-get update \nRUN apt-get install iputils-ping". Ping is needed since the containers needs a process running or it will stop and running `ping 8.8.8.8` is an easy hack to allow for the container to act closer to a VM that is always running. A more polished system would require people to write Dockerfiles or bash scripts that specify the process they would like to carry out and report back to them when the task is completed; eliminating the need for a user to lease specific machines.

Once you have made the proper resources available for use to your users you can use the `python client.py lease create <resource> <lease_name>` command to start a lease and `python client.py lease delete <lease_name>` to remove it when you are done.

#### Setup

##### Local Development

1. Install docker or docker-machine for mac
  - OSX Setup
    - run the following in the terminal you will be using the docker command from
    - download https://www.docker.com/toolbox and install
    - `docker-machine create default --driver virtualbox`
    - `eval "$(docker-machine env default)"`
  - Ubuntu 14.04 Setup (resources/setup.sh without the last two lines)
    - `sudo apt-key adv --keyserver hkp://pgp.mit.edu:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D`
    - `sudo touch /etc/apt/sources.list.d/docker.list`
    - `sudo bash -c 'echo "deb https://apt.dockerproject.org/repo ubuntu-trusty main" > /etc/apt/sources.list.d/docker.list'`
    - `sudo apt-get update`
    - `sudo apt-get -y install docker-engine`
    - `sudo usermod -aG docker yourusername`
    - `sudo su - yourusername` or log out and back into your user.
2. `pip install docker-compose` or `sudo pip install docker-compose` depending how python and pip are installed on your machine.
3. run `make test` or `sudo make test` if you run into issues with connecting to the docker socket on Ubuntu. This will build your entire environment inside a container and make it available at port 8080. It will hang at the end giving you a way to see any stderr or stdout that occurs from commands being run against the server. You can hit Ctrl+C to stop this and the server will keep running.
4. Create slaves (optional)
  - If you do not want to use the slave machines I have created you can create your own by starting up a ubuntu:14.04 machine that the master has access to.
  - add `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCvvUN6dmF2t1HH3V/cfQyBm7dxVkKfsEwhGeoxT9pTAxLgNO7JcXnjW65w17vprxLpdoiYp+W9FQAEAN8b95F5Z6PPwiksbIVN8fy3G06Tkdi3kb6bFd9/okU/xPIKpIfOl9LuXND5yM6hNPVwQ/SzoZnGzf+tRxNuq/38umjDwLufLPNTq8Ga4NLc+G30puOlu2YB5+CQgcs0WQZsUADtnT0o+Ai64Wme3W5GbtWyoBIMRvEf5mtSjRJOR469HZM7qiqPY5btyXWShallZYSLQBxF/z4tcgOpFxGDaKfYIehpS+zu9rFOpkqjDwgGQLpjw3/pkjnNNiy2TNGihntr resalloc` to the ~/.ssh/autorized_keys file.
  - run the resources/setup.sh script on the machine.
4. Initial Setup
  - from the client directory run the following commands
  - `pip install -r requirements.txt`
  - `python client.py register yourusername yourpassword` Create a user
  - `python client.py login yourusername yourpassword` Login with your user to get a token
  - `python client.py resources create busybox "FROM busybox"` Setups a small linux distro. You can use any valid docker file for the second parameter if you wish [Dockerfiles](https://docs.docker.com/reference/builder/)
  - `python client.py machine create slave ubuntu 52.10.32.82` This uses an EC2 instance that has had the setup.sh script run on it and has an ubuntu user with a public key created from the resalloc.pem private key.
  - `python client.py lease create busybox mylease` Use the resource we created above and give my lease a name.
  - `python client.py ssh mylease` log into the machine that was just created. Type exit to log out of your leased machine.
    - If the permission of the key has changed run `chmod 600 resources/resalloc.pem`

##### Building

1. Install docker or docker-machine if you are using OSX
2. Install docker-compose via pip
3. `docker-compose build`
4. `docker ps` will show you what IP and Port the container is available on. Note when using docker-machine the IP you access it on is actually shown via `docker-machine ip default`.
5. The remote server needs to have a user setup with a user who's public key was generated from the resalloc.pem and has had resources/setup.sh run on it. 52.10.32.82 has had this setup run on it and can be added using `python client.py machine create slave1 ubuntu 52.10.32.82`. Since the master server only manages leases at least one slave server is required.

##### Bad Things

- Some code inside of controller.go is repeated and should be abstracted out into helper functions.
- Every created user right now is by default activated. I have put an account_type and activated column that could allow for setting an admin account that then has access to approve or reject new user requests.
- Error messages are not appropriate for end users. Some errors that bubble up are right from a SQL error such as UNIQUE constraint failed.
- Migrations aren't handled in a great way right now since it's not possible to roll back the database to a previous version.
- "Files" are stored as strings with \n for line breaks making this a unix specific implementation.
- The remote docker servers do not have any kind of security on port 5555.
- Getting access to your lease actually gives you access to everyones lease. A daemon should be running on the remote system that allows users to have access to only machines that they created.
- The remote server needs to have a user setup with a user who's public key was generated from the resalloc.pem and has had resources/setup.sh run on it. 52.10.32.82 has had this setup run on it and can be added using `python client.py machine create slave1 ubuntu 52.10.32.82`. Since the master server only manages leases at least one slave server is required to be attached to start created leases.



##### Helpful

- `docker -H 0.0.0.0:5555 rm -f $(docker -H 0.0.0.0:5555 ps -a -q)` will clean out your remote docker server of all containers. Useful for when you are testing and delete may have not actually removed the container.

##### Test Plan

This test plan is a manual one right now but does a decent job at ensuring any changes that you made did not break something. This could easily be moved into an automated test in the future using the go http client.

Requires Postman (load in Resalloc.json.postman_collection)

- Run POST /register
  - Should return Success: true JSON
- Run POST /login (Grab token returned to you)
  - Should return Success: true JSON
- Run POST /resource (Fill in the token header with token returned to you)
  - Should return Success: true JSON
- Run POST /machine (Fill in the token header with token returned to you) (52.10.32.82 can be used as a testing server)
  - Should return Success: true JSON
- Run POST /lease (Fill in the token header with token returned to you)
  - Should return Success: true JSON
- Run GET /lease (Fill in the token header with token returned to you)
  - Should return Success: true JSON and List the Lease you just added
- Run `ssh -o StrictHostKeyChecking=no -i resources/resalloc.pem ubuntu@{{ ip from POST /machine }} -t 'docker -H 0.0.0.0:5555 exec -it {{ lease name }} sh'`
  - Should place you inside of a remote container
- Run DEL /lease (Fill in the token header with token returned to you)
  - Should return Success: true JSON
- Run GET /lease (Fill in the token header with token returned to you)
  - Should return Success: true JSON and null for Lease field
