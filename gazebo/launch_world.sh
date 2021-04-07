#!/bin/bash
export GAZEBO_PLUGIN_PATH=${GAZEBO_PLUGIN_PATH}:./build
gazebo gripper.world --verbose
