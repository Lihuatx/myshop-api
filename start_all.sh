#!/bin/bash

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $SCRIPT_DIR

echo "启动 goods-web 服务..."
./start.sh

echo "启动 order-web 服务..."
cd ../../order-web/target
./start.sh

echo "启动 user-web 服务..."
cd ../../user-web/target
./start.sh

echo "启动 userop-web 服务..."
cd ../../userop-web/target
./start.sh

echo "所有服务启动完成！"

# 等待所有服务启动
sleep 5

# 检查服务状态
echo "检查服务状态："
for service in "goods-web-main" "order-web-main" "user-web-main" "userop-web-main"
do
    if pgrep -x "$service" > /dev/null
    then
        echo "$service 运行正常"
    else
        echo "警告: $service 可能未正常启动"
    fi
done