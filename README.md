# Go Final Project

## Текущий статус

Проект уже содержит рабочую основу для аутентификации и клиентского API:

- проверка номера телефона: POST /auth/check-phone
- регистрация: POST /auth/register
- вход по телефону и PIN: POST /auth/login
- обновление токенов: POST /auth/refresh
- вход сотрудника: POST /employees/login
- создание и поиск клиентов: POST /clients, GET /clients
- получение, обновление и деактивация клиента: GET /clients/{id}, PUT /clients/{id}, DELETE /clients/{id}
- идентификация клиента по телефону/ТИН: POST /clients/identify
- пополнение счета: POST /transactions/deposit
- перевод между счетами: POST /transactions/transfer
- история операций счета: GET /transactions/account/{account_id}
- история операций клиента: GET /transactions/client/{client_id}

Также реализованы:
- слой usecase для auth и client
- репозитории для работы с телефонами, аккаунтами, клиентами и пользователями
- JWT-авторизация
- soft-delete через active_to IS NULL
- базовые unit-тесты для auth-логики

## Архитектура

- domain: модели домена, ошибки и структуры данных
- usecase/auth: бизнес-логика проверки номера, регистрации и входа
- usecase/client: бизнес-логика сотрудников и клиентов
- usecase/transaction: бизнес-логика пополнений, переводов и истории операций
- repository/postgres: доступ к PostgreSQL
- delivery/http: HTTP-обработчики и маршруты

## Примечания

- Поиск активных записей выполняется только по данным с active_to IS NULL.
- Регистрация и ключевые изменения выполняются в транзакциях.
- Все изменения баланса и записи операций выполняются внутри транзакций БД.
- PIN и пароли сотрудников хранятся как bcrypt-хеши.
- Для запуска сервера используйте: go run ./cmd/api
- Сервер слушает порт 8080.
