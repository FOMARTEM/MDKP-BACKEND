curl -X POST  http://127.0.0.1:8080/auth \
 -H "Content-Type: application/json" \
  -d ' {"email" : "admin@example.com", "password" : "Test12345"}'

curl -X GET http://127.0.0.1:8080/roles \
-H "Authorization: Bearer $JWT_TOKEN"

export JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwicm9sZSI6ItCQ0LTQvNC40L3QuNGB0YLRgNCw0YLQvtGAIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsIm5hbWUiOiLQn9C10YDQstGL0Lkg0JDQtNC80LjQvdC40YHRgtGA0LDRgtC-0YAiLCJleHAiOjE4MDYyNDQzMzksImlhdCI6MTc3NzQ0NDMzOX0.SSdgvF_N6r6hxX5Nr0Cxq0AdkSfmJYpWnGIHPnOMbOQ"

curl -X PUT http://localhost:8080/account \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -F "old_password=Test12345" \
  -F "new_password=Test123456" \
  -F "new_password_confirm=Test123456"

curl -X POST http://localhost:8080/user \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $JWT_TOKEN" \
-d '{
  "last_name": "Первый",
  "first_name": "Автор",
  "middle_name": "Платформы",
  "phone": "89876543210",
  "date_of_birth": "1985-05-15",
  "email": "author@example.com",
  "password": "Test1234",
  "role": "Автор"
}'

curl -X POST  http://127.0.0.1:8080/auth \
-H "Content-Type: application/json" \
-d ' {"email" : "author@example.com", "password" : "Test1234"}'

curl -X PUT  http://127.0.0.1:8080/user/active?email=author@example.com \
-H "Authorization: Bearer $JWT_TOKEN" 

curl -X GET http://localhost:8080/finduser \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $JWT_TOKEN" \
-d '{
  "role": "Автор"
}'

curl -X PUT  "http://127.0.0.1:8080/user/role?email=author@example.com&roleID=2" \
-H "Authorization: Bearer $JWT_TOKEN" 

curl -X GET  http://127.0.0.1:8080/user/list \
-H "Authorization: Bearer $JWT_TOKEN" 