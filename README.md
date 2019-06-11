
# RhoMBus

[![Build Status](https://travis-ci.org/l1va/rhombus.svg?branch=master)](https://travis-ci.org/l1va/rhombus)

<b>This is ROS competitor written in 100 lines. Robotic MicroService Bus.</b> 

The main idea - to use [NATS](https://nats.io/) message queue as a 
master or server.
Instead of nodes - microservices. The rest is almost the same. 
You can use any language in microservices: NATS supports a lot of 
[clients](https://nats.io/download/)

 ### Run
 Just start NATS server (run_server.sh) and you can run any microservice 
 (see examples) without annoying CMakeLists.txt, package.xml, catkin, etc. 
 
 