#!/bin/sh

# remove
if systemctl status reverse-ssh.service | grep "active (running)" -c > 0; then
  echo "ssh tunnel daemon active"
  systemctl restart reverse-ssh.service
else
  echo "ssh tunnel daemon inactive"
  systemctl start reverse-ssh.service
fi
