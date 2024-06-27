#!/bin/bash

export PATH=$PATH:/home/cat/go/bin
export GOROOT=$HOME/go
export GOPATH=$HOME/go

# 设置路径
WORK_DIR="$HOME/game_go"
LOG_FILE="$WORK_DIR/game.log"
PROCESS_NAME="main.go"
GO_PATH="/home/cat/go/bin/go"  # 使用你的 go 可执行文件路径

# 检查进程是否正在运行
if ! pgrep -f "$PROCESS_NAME" > /dev/null; then
    # 如果进程未运行，则启动它并记录日志
    cd "$WORK_DIR"
    nohup $GO_PATH run "$PROCESS_NAME" >> "$LOG_FILE" 2>&1 &
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $PROCESS_NAME started" >> "$LOG_FILE"
else
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $PROCESS_NAME is already running" >> "$LOG_FILE"
fi

