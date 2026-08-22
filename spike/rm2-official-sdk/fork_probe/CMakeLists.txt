cmake_minimum_required(VERSION 3.16)
project(Writerdeck_fork_probe VERSION 0.1 LANGUAGES C CXX)
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

if(NOT FORK_DIR)
  set(FORK_DIR "${CMAKE_CURRENT_SOURCE_DIR}/../.cache/writerdeck-keywriter")
endif()
if(NOT EXISTS "${FORK_DIR}/main.cpp")
  message(FATAL_ERROR "FORK_DIR missing fork checkout: ${FORK_DIR}")
endif()

find_package(Qt6 REQUIRED COMPONENTS Quick Qml)
qt_standard_project_setup(REQUIRES 6.5)

set(SUNDOWN "${FORK_DIR}/sundown")
set(SUNDOWN_SRCS
  "${SUNDOWN}/src/autolink.c"
  "${SUNDOWN}/src/buffer.c"
  "${SUNDOWN}/src/markdown.c"
  "${SUNDOWN}/src/stack.c"
  "${SUNDOWN}/html/houdini_href_e.c"
  "${SUNDOWN}/html/houdini_html_e.c"
  "${SUNDOWN}/html/html.c"
  "${SUNDOWN}/html/html_smartypants.c")

qt_add_executable(Writerdeck_fork_probe
  "${FORK_DIR}/main.cpp"
  "${FORK_DIR}/edit_utils.cpp"
  "${FORK_DIR}/rotation_watcher.cpp"
  "${FORK_DIR}/lobby_bridge.cpp"
  "${FORK_DIR}/lobby_ui_config.cpp"
  "${FORK_DIR}/edit_helper.cpp"
  ${SUNDOWN_SRCS})

set(FORK_MAIN_QML "${FORK_DIR}/main.qml")
set_source_files_properties("${FORK_MAIN_QML}" PROPERTIES QT_RESOURCE_ALIAS "main.qml")

qt_add_resources(Writerdeck_fork_probe "fork_qml"
  PREFIX "/"
  FILES "${FORK_MAIN_QML}")

target_compile_definitions(Writerdeck_fork_probe PRIVATE WRITERDECK_RM2_QT6_PROBE=1)

target_include_directories(Writerdeck_fork_probe PRIVATE
  "${FORK_DIR}"
  "${SUNDOWN}/src"
  "${SUNDOWN}/html"
  "${CMAKE_CURRENT_BINARY_DIR}")

target_link_libraries(Writerdeck_fork_probe PRIVATE Qt6::Quick Qt6::Qml)
target_link_options(Writerdeck_fork_probe PRIVATE -pthread)
target_link_libraries(Writerdeck_fork_probe PRIVATE pthread)
