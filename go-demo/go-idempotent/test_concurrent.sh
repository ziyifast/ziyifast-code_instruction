#!/bin/bash

# 生成测试token
TOKEN=$(uuidgen)
echo "测试并发请求，使用token: $TOKEN"

# 并发发送10个相同的请求
for i in {1..10}; do
  curl -X POST http://localhost:8080/payment \
    -H "Content-Type: application/json" \
    -d '{
      "idempotent_token": "'$TOKEN'",
      "amount": '$i'.00,
      "user_id": "user_'$i'",
      "description": "并发测试请求 '$i'"
    }' &
done

# 等待所有请求完成
wait

echo "并发测试完成，查询结果:"
curl "http://localhost:8080/payment/status?token=$TOKEN"
