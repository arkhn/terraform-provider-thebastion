---
page_title: "thebastion Provider Tests Guide"
subcategory: ""
description: |-
  
---
# Test

## Setup

The tests need a fresh TheBastion environment. This environment can be generated using the `docker-compose.test.yaml` file available in this repository. 
Be careful to name your own public key `idthebastion.pub` or to modify `docker-compose.test.yaml` file to use your own public key for the `poweruser` user to be able to access the TheBastion test environment.

### Generate ssh key

Run the following command to generate your ssh key.

```shell
$ ssh-keygen -f idthebastion
```

The following output is print, press enter twice. You must not set a password for the `idthebastion` key.

```shell
Generating public/private rsa key pair.
Enter passphrase (empty for no passphrase):
Enter same passphrase again:
```

### If you have change the name or path of key

Replace <path_key> with the path to the **public key** you have created, in your `docker-compose.test.yaml` file.

```docker-compose
version: "3.7"
services:

  thebastion:
    image: ovhcom/the-bastion:sandbox
    container_name: bastiontest
    ports:
      - 1122:22
    volumes:
      - ./entrypoint.sh:/mount/entrypoint.sh
      - ./<path_key>:/mount/idthebastion.pub
    entrypoint: /mount/entrypoint.sh
```

### Launch TheBastion Setup

Run the following command to launch TheBastion setup.

```shell
$ docker-compose -f docker-compose.test.yaml up
```

## Usage

You must run the following command, if you have changed a function inside the project. 
```shell
make install
```

Run the following command to set your THEBASTION environment variables
```shell
export THEBASTION_HOST=127.0.0.1:1122 \
THEBASTION_USERNAME=poweruser \
THEBASTION_PATH_KNOWN_HOST=$HOME/.ssh/known_hosts \
THEBASTION_PATH_PRIVATE_KEY=$HOME/<path_to_idthebastion_or_private_key>
```

Run the following command to run the tests. 
```shell
make testacc
``` 

## Tests description

TheBastion has two users in the default setting: `poweruser` and `healthcheck`.
`poweruser` is used to call TheBastion, while `healthcheck` is a user created by TheBastion.
Keep this in mind when creating new tests involving the number of users.
This is also why we terminate the testsuite if the default setting is not present at the start.

The `PreCheck` and `CheckDestroy` functions test that the number of users at the beginning and at the end of the tests is 2. This is why the tests need a fresh test environment.

Some custom utils functions for testing are in `thebastion/tests` folder. 

### Resource
#### User

`ingress_keys_base` and `ingress_keys_update` are two public keys generated for the purposes of the tests. You may reuse them or generate new keys.
Be mindful that keys used during tests must be valid public SSH keys, as TheBastion checks the keys' format.

The following tests are realized:
- user creation
- user update

### Data source
#### Users

Use some features of user_resource to populate TheBastion with Users.

The following tests are realized to read function:
- in default settings,
- when single user was created,
- when multiples users were created,
