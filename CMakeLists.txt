cmake_minimum_required(VERSION 3.16)
project(Gast VERSION 1.0.0 LANGUAGES NONE)

# 查找Go编译器
find_program(GO_EXECUTABLE go REQUIRED)
if(NOT GO_EXECUTABLE)
    message(FATAL_ERROR "Go compiler not found. Please install Go and add it to PATH.")
endif()

# 显示Go版本信息
execute_process(
    COMMAND ${GO_EXECUTABLE} version
    OUTPUT_VARIABLE GO_VERSION
    OUTPUT_STRIP_TRAILING_WHITESPACE
)
message(STATUS "Found Go: ${GO_EXECUTABLE} (${GO_VERSION})")

# 设置输出目录
set(OUTPUT_DIR ${CMAKE_BINARY_DIR}/bin)
file(MAKE_DIRECTORY ${OUTPUT_DIR})

# 定义二进制文件名
if(WIN32)
    set(BINARY_NAME gast.exe)
else()
    set(BINARY_NAME gast)
endif()

set(BINARY_PATH ${OUTPUT_DIR}/${BINARY_NAME})

# Go源文件列表
set(GO_SOURCES
    main.go
    cmd_basic.go
    cmd_config.go
    cmd_file.go
    cmd_network.go
    cmd_grep.go
    cmd_interactive.go
    config.go
    utils.go
    go.mod
)

# 构建目标
add_custom_command(
    OUTPUT ${BINARY_PATH}
    COMMAND ${GO_EXECUTABLE} build -o ${BINARY_PATH} .
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    DEPENDS ${GO_SOURCES}
    COMMENT "Building Gast with Go compiler"
    VERBATIM
)

# 创建自定义目标
add_custom_target(gast ALL DEPENDS ${BINARY_PATH})

# 清理目标
add_custom_target(clean-go
    COMMAND ${CMAKE_COMMAND} -E remove_directory ${OUTPUT_DIR}
    COMMAND ${GO_EXECUTABLE} clean
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    COMMENT "Cleaning Go build artifacts"
)

# 测试目标
add_custom_target(test-go
    COMMAND ${GO_EXECUTABLE} test -v ./...
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    COMMENT "Running Go tests"
)

# 格式化目标
add_custom_target(format-go
    COMMAND ${GO_EXECUTABLE} fmt ./...
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    COMMENT "Formatting Go code"
)

# 静态检查目标
add_custom_target(vet-go
    COMMAND ${GO_EXECUTABLE} vet ./...
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    COMMENT "Running Go vet"
)

# 安装目标
install(PROGRAMS ${BINARY_PATH} DESTINATION bin)

# 跨平台编译目标
add_custom_target(build-all
    COMMAND ${CMAKE_COMMAND} -E make_directory ${OUTPUT_DIR}/release
    
    # Linux amd64
    COMMAND ${CMAKE_COMMAND} -E env GOOS=linux GOARCH=amd64 
            ${GO_EXECUTABLE} build -o ${OUTPUT_DIR}/release/gast-linux-amd64 .
    
    # Linux arm64
    COMMAND ${CMAKE_COMMAND} -E env GOOS=linux GOARCH=arm64 
            ${GO_EXECUTABLE} build -o ${OUTPUT_DIR}/release/gast-linux-arm64 .
    
    # Windows amd64
    COMMAND ${CMAKE_COMMAND} -E env GOOS=windows GOARCH=amd64 
            ${GO_EXECUTABLE} build -o ${OUTPUT_DIR}/release/gast-windows-amd64.exe .
    
    # macOS amd64
    COMMAND ${CMAKE_COMMAND} -E env GOOS=darwin GOARCH=amd64 
            ${GO_EXECUTABLE} build -o ${OUTPUT_DIR}/release/gast-darwin-amd64 .
    
    # macOS arm64 (Apple Silicon)
    COMMAND ${CMAKE_COMMAND} -E env GOOS=darwin GOARCH=arm64 
            ${GO_EXECUTABLE} build -o ${OUTPUT_DIR}/release/gast-darwin-arm64 .
    
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    COMMENT "Building for all platforms"
)

# 显示构建信息
message(STATUS "Project: ${PROJECT_NAME} v${PROJECT_VERSION}")
message(STATUS "Binary will be built to: ${BINARY_PATH}")
message(STATUS "Available targets:")
message(STATUS "  gast        - Build the main binary")
message(STATUS "  clean-go    - Clean build artifacts")
message(STATUS "  test-go     - Run tests")
message(STATUS "  format-go   - Format Go code")
message(STATUS "  vet-go      - Run Go vet")
message(STATUS "  build-all   - Cross-compile for all platforms")
message(STATUS "  install     - Install binary to system") 