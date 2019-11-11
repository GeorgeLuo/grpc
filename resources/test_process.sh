#!/bin/sh
(
sleep 1
echo "I'm the script with pid $$"
for i in $(seq 1 20); do
        sleep 1
        echo "Still running $$"
done
)
