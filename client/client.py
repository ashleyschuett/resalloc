#!/usr/bin/env python

from docopt import docopt
from inspect import getdoc
import platform
import requests
import json
import os

def main():
    # Change this to whatever IP and PORT you are running
    # the server on. likely localhost for linux.
    # auto detect and try to configure the port to use
    server_ip = "0.0.0.0"
    # If docker-machine is not running this will output
    # host is not running to your terminal
    if platform.system() == "Darwin":
        server_ip = os.popen('docker-machine ip default').read().rstrip()
    server_ip = server_ip + ":8080"
    # Uncomment line to use pre configured external server
    # server_ip = "52.89.15.147:8080"
    command = CompressCommand(server_ip)
    command.run()

class CompressCommand(object):
    """
    Usage:
      client.py register <name> <password>
      client.py login <name> <password>
      client.py resources list
      client.py resources create <name> <dockerfile>
      client.py machine create <machine_name> <login> <ip>
      client.py lease list
      client.py lease create <resource_name> <lease_name>
      client.py lease delete <lease_name>
      client.py ssh <lease_name>
    """

    def __init__(self, server_ip):
        self.server_ip = server_ip
        self.arguments = None


    def run(self):
        doc = getdoc(self)
        self.arguments = docopt(doc)
        # Check what command was run
        register = self.arguments.get('register')
        login = self.arguments.get('login')
        resources = self.arguments.get('resources')
        machine = self.arguments.get('machine')
        lease = self.arguments.get('lease')
        ssh = self.arguments.get('ssh')
        if register:
            print self.register()
        if login:
            print self.login()
        if resources:
            get_list = self.arguments.get('list')
            create = self.arguments.get('create')
            if get_list:
                print self.resources_list()
            if create:
                print self.resources_create()
        if machine:
            create = self.arguments.get('create')
            if create:
                print self.machine_create()
        if lease:
            get_list = self.arguments.get('list')
            create = self.arguments.get('create')
            delete = self.arguments.get('delete')
            if get_list:
                print self.lease_list()
            if create:
                print self.lease_create()
            if delete:
                print self.lease_delete()
        if ssh:
            print self.make_ssh_connection()

    # Calls the POST /register endpoint and
    # trys to create a new user account
    def register(self):
        headers = {'Content-type': 'application/json'}
        structure = {
            'Name': self.arguments.get("<name>"),
            'Password': self.arguments.get("<password>")
        }
        json_encoded = json.dumps(structure)
        res = requests.post(
            'http://'+self.server_ip+"/register",
            headers=headers,
            data=json_encoded)
        return self.print_json(res.json())

    # Calls the POST /login endpoint and
    # saves your token for use in later calls
    def login(self):
        headers = {'Content-type': 'application/json'}
        structure = {
            'Name': self.arguments.get("<name>"),
            'Password': self.arguments.get("<password>")
        }
        json_encoded = json.dumps(structure)
        res = requests.post(
            'http://'+self.server_ip+"/login",
            headers=headers,
            data=json_encoded)
        # Write token to file and close
        token_file = open(".token", 'w')
        token_file.write(res.json().get('Token'))
        token_file.close()
        return self.print_json(res.json())

    # Get a list of all available resources for use
    # on the remote server
    def resources_list(self):
        headers = {'token': self.get_user_token()}
        res = requests.get(
            'http://'+self.server_ip+"/resource",
            headers=headers)
        return self.print_json(res.json())

    # Get a list of all available resources for use
    # on the remote server
    def resources_create(self):
        headers = {'token': self.get_user_token(),
                'Content-type': 'application/json'}
        structure = {
            'Name': self.arguments.get("<name>"),
            'File': self.arguments.get("<dockerfile>")
        }
        json_encoded = json.dumps(structure)
        res = requests.post(
            'http://'+self.server_ip+"/resource",
            headers=headers,
            data=json_encoded)
        return self.print_json(res.json())

    # Add a machine to the master server
    def machine_create(self):
        headers = {'token': self.get_user_token(),
                'Content-type': 'application/json'}
        structure = {
            'Name': self.arguments.get("<machine_name>"),
            'Username': self.arguments.get("<login>"),
            'IP': self.arguments.get("<ip>")
        }
        json_encoded = json.dumps(structure)
        res = requests.post(
            'http://'+self.server_ip+"/machine",
            headers=headers,
            data=json_encoded)
        return self.print_json(res.json())

    # List all the currently available leases
    def lease_list(self):
        headers = {'token': self.get_user_token()}
        res = requests.get(
            'http://'+self.server_ip+"/lease",
            headers=headers)
        return self.print_json(res.json())

    def lease_create(self):
        headers = {'token': self.get_user_token(),
                'Content-type': 'application/json'}
        structure = {
            'ResourceName': self.arguments.get("<resource_name>"),
            'LeaseName': self.arguments.get("<lease_name>")
        }
        json_encoded = json.dumps(structure)
        res = requests.post(
            'http://'+self.server_ip+"/lease",
            headers=headers,
            data=json_encoded)
        return self.print_json(res.json())

    def lease_delete(self):
        headers = {'token': self.get_user_token(),
                'Content-type': 'application/json'}
        structure = {
            'Name': self.arguments.get("<lease_name>")
        }
        json_encoded = json.dumps(structure)
        res = requests.delete(
            'http://'+self.server_ip+"/lease",
            headers=headers,
            data=json_encoded)
        return self.print_json(res.json())

    def make_ssh_connection(self):
        lease_name = self.arguments.get('<lease_name>')
        lease_list_json = self.lease_list()
        lease_list_dict = json.loads(lease_list_json)
        leases = lease_list_dict.get('Leases')
        machine_info = None
        for lease in leases:
            if lease.get('Name') == lease_name:
                machine_info = lease
                continue
        if machine_info == None:
            print "This lease doesn't exist"
            exit()
        os.system("ssh -o StrictHostKeyChecking=no -i ./../resources/resalloc.pem "+machine_info.get('Username')+"@"+machine_info.get('MachineName')+" -t 'docker -H 0.0.0.0:5555 exec -it "+lease_name+" sh'")

    # helper function to make the json
    # more readable to end end user
    def print_json(self, j):
        return json.dumps(j, sort_keys=True,
                  indent=4, separators=(',', ': '))

    # attempt to get a user token if
    # it is currenently available
    def get_user_token(self):
        token_file = open(".token", "r")
        user_token = token_file.readlines()
        token_file.close()
        return user_token[0]


if __name__ == "__main__":
    main()
