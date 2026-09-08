cmake_minimum_required(VERSION 3.16)
project(textedit_spike VERSION 0.1 LANGUAGES CXX)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(Qt6 REQUIRED COMPONENTS Quick)
qt_standard_project_setup(REQUIRES 6.5)

qt_add_executable(textedit_spike main.cpp)
qt_add_qml_module(textedit_spike
    URI remarkable_textedit_spike
    VERSION 1.0
    QML_FILES Main.qml)

target_link_libraries(textedit_spike PRIVATE Qt6::Quick)
