# resalloc
resource allocation with docker

#### Setup

##### Local Development

1. Install docker or docker-machine for mac
  - OSX Setup
    - run the following in the terminal you will be using the docker command from
    - download https://www.docker.com/toolbox and install
    - `docker-machine create default --driver virtualbox`
    - `eval "$(docker-machine env default)"`
  - Ubuntu 14.04 Setup
    - `sudo apt-key adv --keyserver hkp://pgp.mit.edu:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D`
    - `sudo touch /etc/apt/sources.list.d/docker.list`
    - `sudo bash -c 'echo "deb https://apt.dockerproject.org/repo ubuntu-trusty main" > /etc/apt/sources.list.d/docker.list'`
    - `sudo apt-get update`
    - `sudo apt-get -y install docker-engine`
    - `sudo usermod -aG docker yourusername`
    - `sudo su - yourusername` or log out and back into your user.
2. `pip install docker-compose` or `sudo pip install docker-compose` depending how python and pip are installed on your machine.
3. run `make test` or `sudo make test` if you run into issues with connecting to the docker socket on Ubuntu. This will build your entire environment inside a container and make it available at port 8080.
4. Configure the client.
  - OSX
    - `docker-machine ip default`
      - this is the IP that you should place on line 11 for client.py to use.
  - Ubuntu
    - run `sudo docker ps` and use the IP and PORT that is listed under PORTS
      - this is the IP that you should place on line 11 for client.py to use.
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

- Every created user right now is by default activated. I have put an account_type and activated column that could allow for setting an admin account that then has access to approve or reject new user requests.
- Error messages are not appropriate for end users. Some errors that bubble up are right from a SQL error such as UNIQUE constraint failed.
- Migrations aren't handled in a great way right now since it's not possible to roll back the database to a previous version.
- "Files" are stored as strings with \n for line breaks making this a unix specific implementation.
- The remote docker servers do not have any kind of security on port 5555.
- Getting access to your lease actually gives you access to everyones lease. A daemon should be running on the remote system that allows users to have access to only machines that they created.
- The remote server needs to have a user setup with a user who's public key was generated from the resalloc.pem and has had resources/setup.sh run on it. 52.10.32.82 has had this setup run on it and can be added using `python client.py machine create slave1 ubuntu 52.10.32.82`. Since the master server only manages leases at least one slave server is required to be attached to start created leases.



#### Helpful

- `docker -H 0.0.0.0:5555 rm -f $(docker -H 0.0.0.0:5555 ps -a -q)` will clean out your remote docker server of all containers. Useful for when you are testing and delete may have not actually removed the container.

#### Test Plan

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
