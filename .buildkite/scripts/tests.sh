#!/bin/bash

mise trust -y
mise install

sleep 5

mise x -- gotestsum --format testname
