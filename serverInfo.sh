#!/bin/bash
for host in `cat servers`
do
        for port in `cat ports`
        do
            ./cmd -s "${host}" -p "${port}" -ping
        done
done
