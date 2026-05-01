curl -vX POST  http://127.0.0.1:8080/auth \
 -H "Content-Type: application/json" \
  -d ' {"email" : "admin@example.com", "password" : "Test12345"}'

curl -vX GET http://127.0.0.1:8080/roles \
-H "Authorization: Bearer $JWT_TOKEN"

curl -vX GET http://127.0.0.1:8080/status \
-H "Authorization: Bearer $JWT_TOKEN"

export JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwicm9sZSI6ItCQ0LTQvNC40L3QuNGB0YLRgNCw0YLQvtGAIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsIm5hbWUiOiLQn9C10YDQstGL0Lkg0JDQtNC80LjQvdC40YHRgtGA0LDRgtC-0YAiLCJleHAiOjE4MDYyNDQzMzksImlhdCI6MTc3NzQ0NDMzOX0.SSdgvF_N6r6hxX5Nr0Cxq0AdkSfmJYpWnGIHPnOMbOQ"

curl -vX PUT http://localhost:8080/account \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -F "old_password=Test12345" \
  -F "new_password=Test123456" \
  -F "new_password_confirm=Test123456"

curl -vX POST http://localhost:8080/user \
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

curl -vX POST  http://127.0.0.1:8080/auth \
-H "Content-Type: application/json" \
-d ' {"email" : "rukovod@example.com", "password" : "Test1234"}'

curl -vX PUT  http://127.0.0.1:8080/user/active?email=author@example.com \
-H "Authorization: Bearer $JWT_TOKEN" 

curl -vX GET http://localhost:8080/finduser \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $JWT_TOKEN" \
-d '{
  "role": "Автор"
}'

curl -vX PUT  "http://127.0.0.1:8080/user/role?email=author@example.com&roleID=2" \
-H "Authorization: Bearer $JWT_TOKEN" 

curl -vX GET  http://127.0.0.1:8080/user/list \
-H "Authorization: Bearer $JWT_TOKEN" 


# Руководитель 
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Miwicm9sZSI6ItCg0YPQutC-0LLQvtC00LjRgtC10LvRjCIsImVtYWlsIjoicnVrb3ZvZEBleGFtcGxlLmNvbSIsIm5hbWUiOiLQn9C10YDQstGL0Lkg0KDRg9C60L7QstC-0LTQuNGC0LXQu9GMIiwiZXhwIjoxODA2MjUxOTg4LCJpYXQiOjE3Nzc0NTE5ODh9.CEB1OqZzcAbgkAbRE5Ruj33PQvt0FPh9nbhVvR3wtxM

export Ruk_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Miwicm9sZSI6ItCg0YPQutC-0LLQvtC00LjRgtC10LvRjCIsImVtYWlsIjoicnVrb3ZvZEBleGFtcGxlLmNvbSIsIm5hbWUiOiLQn9C10YDQstGL0Lkg0KDRg9C60L7QstC-0LTQuNGC0LXQu9GMIiwiZXhwIjoxODA2MjUxOTg4LCJpYXQiOjE3Nzc0NTE5ODh9.CEB1OqZzcAbgkAbRE5Ruj33PQvt0FPh9nbhVvR3wtxM"

curl -vX GET http://127.0.0.1:8080/account/my \
-H "Authorization: Bearer $Ruk_TOKEN"

curl -vX POST http://127.0.0.1:8080/task \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $Ruk_TOKEN" \
  -d '{
    "title": "Для удаления задача",
    "description": "Описание задачи",
    "date_created": "2026-04-29",
    "date_deadline": "2026-05-15",
    "priority": 5,
    "id_status": 1,
    "id_redactor": 4,
    "id_author": 6
  }'

curl -vX DELETE  http://127.0.0.1:8080/task/:id \
-H "Authorization: Bearer $Ruk_TOKEN" 

curl -vX GET  http://127.0.0.1:8080/task/1 -H "Authorization: Bearer $Ruk_TOKEN"

curl -vX PUT http://127.0.0.1:8080/task/1/status \
  -H "Authorization: Bearer $Ruk_TOKEN" \
  -F "status=В работе"


curl -vX GET  http://127.0.0.1:8080/task/list -H "Authorization: Bearer $Ruk_TOKEN"

curl -vX GET http://127.0.0.1:8080/material/1 \
-H "Authorization: Bearer $Ruk_TOKEN"

curl -X POST http://127.0.0.1:8080/material/1 \
  -H "Authorization: Bearer $Ruk_TOKEN" \
  -F "description=Описание материала" \
  -F "file=@/home/shahov/Documents/Test/Шахов-ПР1.1-1.4-1.pdf"

curl -X GET http://127.0.0.1:8080/material/8 \
  -H "Authorization: Bearer $Ruk_TOKEN"

curl -X GET http://127.0.0.1:8080/material/download/8 \
  -H "Authorization: Bearer $Ruk_TOKEN" \
  -J -L   --remote-name
