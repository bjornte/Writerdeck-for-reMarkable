cmake_minimum_required(VERSION 3.16)
# Spike wrapper: build Writerdeck-keywriter from a sibling checkout.
project(Writerdeck_fork_probe_stub LANGUAGES NONE)

if(NOT FORK_DIR)
  set(FORK_DIR "${CMAKE_CURRENT_SOURCE_DIR}/../.cache/writerdeck-keywriter")
endif()

if(NOT EXISTS "${FORK_DIR}/CMakeLists.txt")
  message(FATAL_ERROR "FORK_DIR missing CMakeLists.txt: ${FORK_DIR}")
endif()

set(WRITERDECK_RM2_QT6_PROBE ON CACHE BOOL "Spike verify screenshot hook" FORCE)
add_subdirectory("${FORK_DIR}" fork_build)
