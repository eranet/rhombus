
# RhoMBus

[![Build Status](https://travis-ci.org/eranet/rhombus.svg?branch=master)](https://travis-ci.org/eranet/rhombus)

<b>This is ROS competitor written in 100 lines. Robotic MicroService Bus.</b> 

The main idea - to use [NATS](https://nats.io/) message queue as a 
master or server.
Instead of nodes - microservices. The rest is almost the same. 
You can use any language in microservices: NATS supports a lot of 
[clients](https://nats.io/download/)

 ### Run
Just start NATS server 
    
    ./run_server.sh
and then run any microservice
    
    go run rhomgo/example/pubsub/subscriber_ms/subscriber.go  
    go run rhomgo/example/pubsub/publisher_ms/publisher.go
 (python or golang, see examples) without annoying CMakeLists.txt, 
 package.xml, catkin, etc. 
 
 ### Full App example (Gazebo + RhoMGo)
 https://github.com/eranet/rhombus_example
 