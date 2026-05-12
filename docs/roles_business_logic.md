# Бизнес-логика по ролям

Документ описывает **ожидаемую** бизнес-логику по ролям в MDKP backend и её привязку к текущим HTTP-ручкам.

Важно:
- Сейчас в API **нет явной проверки ролей** на уровне middleware/хендлеров (кроме факта наличия JWT). Роль кладётся в контекст (`role`), но **не используется** для авторизации действий.
- Поэтому ниже приведено: (1) как **должно быть** по предметной области, (2) что **есть сейчас** в реализации (ручки/валидации).

## Общие сущности и процессы

- Пользователи: создание, активация/деактивация, смена роли, просмотр списка.
- Задачи: создание, просмотр, список, удаление (с ограничениями), смена статуса.
- Материалы: загрузка файла к задаче, просмотр метаданных, список, скачивание.
- Версии: создание версии задачи, список версий, получение версии.
- Правки (ревизии): создание правки для версии, список правок, получение правки, смена статуса правки.
- Логи: просмотр логов активности с фильтрацией.

Набор ручек (по коду `internal/api/helpers.go`):
- `POST /auth` (без JWT)
- `GET /health` (без JWT)
- `PUT /account/password`, `GET /account/my`
- `GET /activitylog`
- `POST /user`, `PUT /user/active`, `PUT /user/role`, `GET /user/list`
- `GET /roles`, `GET /status`, `POST /finduser`
- `POST /task`, `GET /task/:id`, `DELETE /task/:id`, `PUT /task/:id/status`, `GET /task/list`
- `POST /material/:id`, `GET /material/:id`, `GET /material/list/:id`, `GET /material/download/:id`
- `POST /version/:id`, `GET /version/list/:id`, `GET /version/:id`
- `POST /edit/:id`, `PUT /edit/:id/status`, `GET /edit/list/:id`, `GET /edit/:id`

## Роль: Администратор

Функция | Реализация (ручки/правила)
---|---
Управление пользователями (создание) | `POST /user` — JSON bind + `validator.Struct`, проверка `phone != ""`, далее `Usecase.CreateUser`
Управление пользователями (активность) | `PUT /user/active?email=...` — проверка query, далее `Usecase.UserActiveChange`
Управление пользователями (смена роли) | `PUT /user/role?email=...&roleID=...` — проверка query + `roleID` int, далее `Usecase.UserRoleUpdate`
Просмотр пользователей | `GET /user/list` — `Usecase.GetUsers`
Просмотр ролей | `GET /roles` — `Usecase.GetRoles`
Просмотр статусов | `GET /status` — `Usecase.GetStatuses`
Поиск пользователей | `POST /finduser` — bind `entities.User`, далее `Usecase.GetUsersByRole` (по одному из полей: `role/email/last_name/first_name`)
Просмотр логов активности | `GET /activitylog?user_id=&email=&start_date=&end_date=&limit=&offset=` — query-парсинг и делегация в `Usecase.GetLogs`

Ожидаемые ограничения (как должно быть):
- Только администратор может: создавать пользователей, менять роли, активность, видеть все логи.
- Валидация при создании пользователя: email/пароль/телефон/ФИО, уникальность email, корректность роли.

## Роль: Руководитель

Функция | Реализация (ручки/правила)
---|---
Создание задач | `POST /task` — JSON bind + `validator.Struct`, `IdCreator` берётся из JWT (`id`), далее `Usecase.CreateTask`
Просмотр задачи | `GET /task/:id` — `taskId` int, далее `Usecase.TaskGetById`
Список своих задач | `GET /task/list` — `userID` из JWT, далее `Usecase.TasksList`
Удаление задачи | `DELETE /task/:id` — `id` int, далее `Usecase.TaskDelete` (ошибка мапится на “нельзя удалить когда взята в работу”)
Смена статуса задачи | `PUT /task/:id/status` — form `status` обязателен, далее `Usecase.TaskStatusUpdate`
Работа с версиями | `POST /version/:taskId`, `GET /version/list/:taskId`, `GET /version/:id`
Работа с материалами (к задаче) | `POST /material/:taskId` — multipart, файл обязателен; имя файла превращается в `Title/Extension`, `CreatorID` из JWT; затем `Usecase.CreateMaterial` и сохранение файла на диск
Просмотр материалов | `GET /material/:id`, `GET /material/list/:taskId`, `GET /material/download/:id`
Работа с правками | `POST /edit/:versionId`, `PUT /edit/:id/status`, `GET /edit/list/:versionId`, `GET /edit/:id`

Ожидаемые ограничения (как должно быть):
- Руководитель создаёт задачи и назначает исполнителей (Автор/Редактор), управляет жизненным циклом задачи.
- Руководитель может менять статус задачи по процессу (например: “Открыта” → “В работе” → “На проверке” → “Закрыта”).
- Руководитель видит задачи своей зоны ответственности (созданные им).

## Роль: Редактор

Функция | Реализация (ручки/правила)
---|---
Просмотр назначенных задач | `GET /task/list` (текущая реализация возвращает по `userID`, логика фильтрации на стороне БД/провайдера)
Работа с версиями | `POST /version/:taskId`, `GET /version/list/:taskId`, `GET /version/:id`
Создание правок к версии | `POST /edit/:versionId` — `CreatorID` из JWT, далее `Usecase.CreateRevision`
Смена статуса правки | `PUT /edit/:id/status` — form `status` обязателен, далее `Usecase.EditStatusUpdate`
Просмотр правок | `GET /edit/list/:versionId`, `GET /edit/:id`
Работа с материалами | `GET /material/...`, `GET /material/download/...` (и при необходимости `POST /material/:taskId`)

Ожидаемые ограничения (как должно быть):
- Редактор ведёт проверку/редактуру: создаёт правки, переводит правки/версии по статусам.
- Редактор не должен иметь доступ к админским операциям (пользователи/роли/логи).

## Роль: Автор

Функция | Реализация (ручки/правила)
---|---
Просмотр назначенных задач | `GET /task/list` (по `userID` из JWT)
Просмотр деталей задачи | `GET /task/:id`
Загрузка материалов к задаче | `POST /material/:taskId` — multipart, файл обязателен; материал сохраняется на диск
Просмотр/скачивание материалов | `GET /material/:id`, `GET /material/list/:taskId`, `GET /material/download/:id`
Работа с версиями (если разрешено процессом) | `POST /version/:taskId`, `GET /version/list/:taskId`, `GET /version/:id`

Ожидаемые ограничения (как должно быть):
- Автор работает с контентом задачи: добавляет материалы, может создавать версии (если это часть процесса) или только просматривать.
- Автор не должен: менять статусы задач произвольно, удалять задачи, управлять пользователями/ролями.

## Рекомендуемые правила авторизации (для реализации в коде)

Чтобы роли реально работали, обычно добавляют:
- Политики доступа по ролям (RBAC) на каждый маршрут (например middleware, проверяющий `c.Get("role")`).
- Проверку владения/назначения: руководитель — только свои задачи; автор/редактор — только задачи где назначены (`id_author/id_redactor`) или где они создатели сущности.
- Валидацию переходов статусов (FSM) для задач/правок/версий.

