#!/bin/bash
export GAZEBO_PLUGIN_PATH=${GAZEBO_PLUGIN_PATH}:./build
gazebo $1 --verbose
