#! /usr/bin/env bash
/opt/bastion/bin/admin/setup-first-admin-account.sh poweruser auto < /mount/idthebastion.pub
/opt/bastion/docker/entrypoint.sh --sandbox