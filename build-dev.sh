#!/bin/bash

watchexec --restart -w ./ --no-process-group -- go run .
