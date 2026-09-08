cmake_minimum_required(VERSION 3.16)
project(socket_spike VERSION 0.1 LANGUAGES CXX)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(Qt6 REQUIRED COMPONENTS Quick Network)
qt_standard_project_setup(REQUIRES 6.5)

qt_add_executable(socket_spike main.cpp socket_feed.cpp socket_feed.h)
qt_add_qml_module(socket_spike
    URI remarkable_socket_spike
    VERSION 1.0
    QML_FILES Main.qml)

target_link_libraries(socket_spike PRIVATE Qt6::Quick Qt6::Network)

add_executable(socket_inject inject.cpp)
