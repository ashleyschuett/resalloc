# resalloc
resource allocation with docker

#### Setup

##### Local Development

1. Install docker or docker-machine for mac
2. Install docker-compose
3. run `make test`. This will build your entire environment inside a container and launch it for access on port 8080.
4. 52.10.32.82 serves as a testing server.

##### Building

1. Install docker or docker-machine if you are using OSX
2. Install docker-compose via pip
3. `docker-compose build`
4. `docker ps` will show you what IP and Port the container is available on. Note when using docker-machine the IP you access it on is actually shown via `docker-machine ip default`.

##### Bad Things

- Every created user right now is by default activated. I have put an account_type and activated column that could allow for setting an admin account that then has access to aprove or reject new user requests.
- Error messages are not appropriate for end users. Some errors that bubble up are right from a SQL error such as UNIQUE constraint failed.
- Migrations aren't handled in a great way right now since it's not possible to roll back the database to a previous version.
- "Files" are stored as strings with \n for line breaks making this a unix specific implementation.
- The remote docker servers do not have any kind of security.
- Getting access to your lease actually gives you access to everyones lease. A daemon should be running on the remote system that allows users to have access to only machines that they created. The other option would be on lease creation allowing only an bash file to be passed in that would do everything you need done and having a web interface where they could check on the results.


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
