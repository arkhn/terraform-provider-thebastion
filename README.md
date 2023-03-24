[![CI](https://github.com/arkhn/terraform-provider-thebastion/workflows/ci/badge.svg)](https://github.com/arkhn/terraform-provider-thebastion/actions)

# Terraform Provider TheBastion

<img src="https://raw.githubusercontent.com/hashicorp/terraform-website/d841a1e5fca574416b5ca24306f85a0f4f41b36d/content/source/assets/images/logo-terraform-main.svg" width="300px">

<img src="https://user-images.githubusercontent.com/218502/96882661-d3b21e80-147f-11eb-8d89-a69e37a5870b.png" width="300px">

## Build provider

Run the following command to build the provider

```shell
$ go build -o terraform-provider-thebastion
```

To communicate with TheBastion you must set environment variables.
```shell
export THEBASTION_HOST=host
THEBASTION_USERNAME=username
THEBASTION_PATH_KNOWN_HOST=$HOME/.ssh/known_hosts
THEBASTION_PATH_PRIVATE_KEY=$HOME/.ssh/id_ed25519
```

## Setup (Required if provider not publish to the Terraform Registry)

You need to complete this file for the terraform-provider:
- .terraformrc, you must replace `<Username>` with the username of your session.

If the provider is not on the terraform registery, you must run the following command

```shell
$ cp .terraformrc ~/.terraformrc 
```

This will specify to terraform to look for the local `terraform-provider-thebastion`.

## Test sample configuration

First, build and install the provider.

```shell
$ make install
```

Then, navigate to the `examples` directory. 

```shell
$ cd examples
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```

## Usage

For more information about the `terraform-provider-thebastion` and its features, visit here.