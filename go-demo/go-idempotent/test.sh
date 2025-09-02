   # 1. 生成幂等token
   TOKEN=$(uuidgen)
   echo "使用token: $TOKEN"

   # 2. 发送支付请求
   curl -X POST http://localhost:8080/payment \
     -H "Content-Type: application/json" \
     -d '{
       "idempotent_token": "'$TOKEN'",
       "amount": 100.50,
       "user_id": "user_123",
       "description": "测试支付"
     }'

   # 3. 使用相同token再次发送（测试幂等性）
   curl -X POST http://localhost:8080/payment \
     -H "Content-Type: application/json" \
     -d '{
       "idempotent_token": "'$TOKEN'",
       "amount": 100.50,
       "user_id": "user_123",
       "description": "测试支付"
     }'

   # 4. 查询支付状态
   curl "http://localhost:8080/payment/status?token=$TOKEN"

   # 5. 查看所有记录
   curl http://localhost:8080/payment/records
