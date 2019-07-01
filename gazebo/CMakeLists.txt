cmake_minimum_required(VERSION 2.8 FATAL_ERROR)

find_package(gazebo REQUIRED)

include_directories(src ${GAZEBO_INCLUDE_DIRS})
link_directories(${GAZEBO_LIBRARY_DIRS})
list(APPEND CMAKE_CXX_FLAGS "${GAZEBO_CXX_FLAGS}")

add_library(rhombus_gazebo_plugin SHARED src/joint_state_publisher.cpp)
target_link_libraries(rhombus_gazebo_plugin nats ${GAZEBO_LIBRARIES})
