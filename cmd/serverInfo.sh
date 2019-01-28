#!/bin/bash
for host in `cat servers`
do
        for port in `cat ports`
        do
            ./gogasm -s "${host}" -p "${port}" -ping
        done
done
