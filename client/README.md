# Resalloc Client

#### Installation
- requires python 2.7
- pip install -r requirements.txt

#### Configuration
The client should automatically configure itself to connect to the master server. It will not work on windows. If you are running on windows you would have to manually specify the IP on line 19.

#### Example Commands

- python client.py register mschuett secret
  - Create a new user account
- python client.py login mschuett secret
  - login to your account so a .token file will be generated
- python client.py resources list
  - list all available resources that are available for leasing
- python client.py resources create busybox "FROM busybox"
  - create a new resource for use with leasing
  - the second param is actually a dockerfile which you can set to use any available container on dockerhub. You can also use \n to create multiline dockerfiles.
- python client.py machine create slave2 ubuntu 52.10.32.82
  - create a machine that will be used to deploy leases to. This machine must have had the setup.sh file run on it as well as have a user of your specification set up to use the resalloc.pem file.
- python client.py lease list
  - list all currently available leases.
- python client.py lease create busybox michael1
  - create a lease for your use
- python client.py lease delete michael1
  - delete a lease that you are done with
- python client.py ssh michael
  - ssh into a lease that is available

#### Bad Things
- This should be made using the python setup tools package so it can be installed via pip.
- The ssh connection just makes a call to the os... Sorry windows users. This should use a native python implementation to make it more portable.
