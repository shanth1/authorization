# authorization

## Детализированный поток OIDC-провайдера

### 0. Определения

- **AuthBack**: Ваш OIDC-провайдер, бэкенд. Адрес: `https://api.auth.oidc.com`
- **AuthFrontPortal**: Портал управления аккаунтом пользователя. Адрес: `https://portal.auth.oidc.com`
- **AuthFrontSSO**: Страница входа в рамках OIDC-потока. Адрес: `https://sso.auth.oidc.com`
- **ClientFront**: Стороннее клиентское приложение, фронтенд. Адрес: `https://app.client.com`
- **ClientBack**: Стороннее клиентское приложение, бэкенд. Адрес: `https://api.client.com`
- **System Admin**: Администратор всей OIDC-инфраструктуры.
- **Client Admin**: Администратор, ответственный за `ClientFront` и `ClientBack`.
- **User**: Конечный пользователь.

---

### Часть 1: Первоначальная подготовка и запуск

Этот этап выполняется один раз для настройки всей системы.

#### 1.1. Запуск инфраструктуры `Auth` (System Admin)

1.  **Действие**: `System Admin` запускает `AuthBack` и связанные с ним сервисы (базу данных PostgreSQL, кэш Redis).
2.  **Действие**: `System Admin` запускает фронтенд-приложения `AuthFrontPortal` и `AuthFrontSSO`. На этом этапе они еще не могут полноценно работать, так как не зарегистрированы как клиенты.

#### 1.2. Генерация криптографических ключей (System Admin)

1.  **Описание**: Для подписи JWT-токенов (`id_token`, `access_token`) `AuthBack` должен иметь набор криптографических ключей (JWK - JSON Web Keys). Необходимо сгенерировать первый ключ.
2.  **Тип**: Команда в CLI или внутренний API-вызов.
3.  **Запрос**: `auth-cli keys generate --alg RS256 --use sig`
4.  **Логика `AuthBack`**:
    - Генерирует новую пару RSA-ключей (приватный и публичный).
    - Присваивает ключу уникальный идентификатор (`kid`).
    - Сохраняет зашифрованный приватный ключ и публичный ключ в таблицу `auth_db.jwks_keys`.
    - Помечает этот ключ как активный для подписи (`use: 'sig'`).
5.  **Результат**: Публичный ключ становится доступен через публичный эндпоинт `/.well-known/jwks.json`, что позволяет клиентам проверять подписи токенов.

#### 1.3. Регистрация `first-party` клиентов (System Admin)

1.  **Описание**: Чтобы портал (`AuthFrontPortal`) и страница входа (`AuthFrontSSO`) могли взаимодействовать с `AuthBack` через OIDC-потоки, их нужно зарегистрировать как клиентов.
2.  **Тип**: `POST` запрос на защищенный Admin API.
3.  **Запрос**:
    - **URL**: `POST https://api.auth.oidc.com/api/v1/admin/clients`
    - **Заголовки**: `Authorization: Bearer <SUPER_ADMIN_TOKEN>`
    - **Тело**:
      ```json
      {
      	"name": "Auth Portal",
      	"redirect_uris": ["https://portal.auth.oidc.com/auth/callback"],
      	"grant_types": ["authorization_code", "refresh_token"],
      	"token_endpoint_auth_method": "client_secret_post"
      }
      ```
4.  **Логика `AuthBack`**: Создает запись в `auth_db.clients`, генерирует `client_id` и `client_secret`.
5.  **Ответ**: `HTTP 201 Created` с `client_id` и `client_secret`.
6.  **Действие**: `System Admin` использует эти креды для настройки `AuthFrontPortal` и `AuthFrontSSO`.

#### 1.4. Регистрация стороннего клиента (Client Admin -> System Admin)

1.  **Описание**: `Client Admin` запрашивает у `System Admin` регистрацию своего приложения ("Магазин").
2.  **Тип**: `POST` запрос на Admin API (выполняется `System Admin`-ом).
3.  **Запрос**:
    - **URL**: `POST https://api.auth.oidc.com/api/v1/admin/clients`
    - **Заголовки**: `Authorization: Bearer <SUPER_ADMIN_TOKEN>`
    - **Тело**:
      ```json
      {
      	"name": "Online Store",
      	"redirect_uris": ["https://app.client.com/auth/callback"],
      	"backchannel_logout_uri": "https://api.client.com/logout/oidc",
      	"grant_types": ["authorization_code", "refresh_token"]
      }
      ```
4.  **Ответ**: `HTTP 201 Created` с `client_id` и `client_secret`.
5.  **Действие**: `System Admin` безопасно передает `client_id` и `client_secret` `Client Admin`-у.

#### 1.5. Настройка и запуск клиентского приложения (Client Admin)

1.  **Действие**: `Client Admin` настраивает переменные окружения для `ClientBack` (`CLIENT_ID`, `CLIENT_SECRET`, `AUTH_ISSUER_URL=https://api.auth.oidc.com`).
2.  **Действие**: `Client Admin` запускает `ClientBack`.
3.  **Логика `ClientBack` при старте (OIDC Discovery)**:
    - **Запрос**: `GET https://api.auth.oidc.com/.well-known/openid-configuration`.
    - **Ответ**: JSON со всеми эндпоинтами (`authorization_endpoint`, `token_endpoint`, `jwks_uri` и т.д.).[10]
    - `ClientBack` кэширует эти URL.
    - **Запрос**: `GET` на `jwks_uri` из предыдущего шага.
    - **Ответ**: JSON с публичными ключами `AuthBack`. `ClientBack` кэширует эти ключи для валидации токенов.[11]
4.  **Действие**: `Client Admin` запускает `ClientFront`, который настроен на взаимодействие с `ClientBack`.

---

### Часть 2: Основной поток авторизации (Authorization Code Flow)

#### 2.1. Инициация (User @ ClientFront)

1.  **Действие**: `User` на `https://app.client.com` нажимает "Войти".
2.  **Логика `ClientFront`**:
    - Генерирует и сохраняет в `sessionStorage`: `state` (от CSRF), `nonce` (от replay-атак), `code_verifier` (для PKCE).
    - Вычисляет `code_challenge` из `code_verifier` (методом S256).
3.  **Редирект**: Браузер перенаправляется на `AuthBack`.
    - **Тип**: `GET`
    - **URL**: `https://api.auth.oidc.com/authorize?response_type=code&client_id=...&redirect_uri=...&scope=openid%20profile&state=...&nonce=...&code_challenge=...&code_challenge_method=S256`

#### 2.2. Обработка запроса на `/authorize` (Исправленная логика)

1.  **Описание**: `AuthBack` получает `GET` запрос на `/authorize`. Это единая точка входа, которая должна определить, залогинен ли пользователь в SSO-сессии.
2.  **Логика `/authorize`**:
    - **ШАГ 1: Проверка SSO-сессии.**
      - Бэкенд пытается прочитать cookie `user_session_id` из запроса.[1][2]
      - Если cookie есть, он валидирует ее значение в Redis.
      - **Если cookie найдена и валидна**: Это **Сценарий Б**. Пользователь уже аутентифицирован. Переходим к **Шагу 2.4**.
      - **Если cookie не найдена или невалидна**: Это **Сценарий А**. Пользователь не аутентифицирован. Выполняем следующие шаги.
    - **ШАГ 2 (Только для Сценария А): Обработка нового входа.**
      - Валидируем параметры из URL: `client_id` существует, `redirect_uri` совпадает с зарегистрированным и т.д.
      - Генерируем `request_id`.
      - **Сохраняем в Redis**: `auth_request:<request_id>` со всеми параметрами изначального запроса (`client_id`, `scope`, `nonce`, `code_challenge` и др.). TTL ~ 10 минут.[3]
      - **Ответ**: Редирект на страницу входа с `request_id`.
        - `HTTP 302 Location: https://sso.auth.oidc.com?request_id=<request_id>`

#### 2.3. Аутентификация на `AuthFrontSSO` (Дополненная логика)

1.  **Действие**: `User` попадает на страницу `sso.auth.oidc.com`.
2.  **Запрос `AuthFrontSSO` -> `AuthBack` для получения провайдеров (НОВЫЙ ШАГ)**:
    - **Описание**: Фронтенд `AuthFrontSSO` должен отобразить правильные способы входа.
    - **Тип**: `GET`
    - **URL**: `https://api.auth.oidc.com/api/v1/auth/providers?request_id=<request_id>` (где `request_id` взят из URL).
    - **Логика `AuthBack`**: `AuthBack` по `request_id` находит в Redis `auth_request`, из него берет `client_id` и возвращает список провайдеров, разрешенных для этого клиента.
    - **Ответ**: `HTTP 200 OK` с JSON: `["password", "google"]`.
3.  **Действие**: `AuthFrontSSO` рендерит кнопки "Войти через Google" и форму для ввода пароля. Пользователь вводит логин/пароль.
4.  **Запрос `AuthFrontSSO` -> `AuthBack` для логина**:
    - **Тип**: `POST`
    - **URL**: `https://api.auth.oidc.com/api/v1/auth/login/password`
    - **Тело**: `{"email": "...", "password": "...", "request_id": "..."}`
5.  **Логика `AuthBack`**:
    - Проверяет учетные данные.
    - Генерирует `session_id`.
    - **Сохраняет в Redis**: `user_session:<session_id>` с `user_id`.
    - **Ответ `AuthBack` -> `AuthFrontSSO`**: `HTTP 200 OK` с заголовком `Set-Cookie: user_session_id=<session_id>; ...`.
6.  **Действие `AuthFrontSSO`**: Получив `200 OK`, `AuthFrontSSO` понимает, что пользователь успешно аутентифицирован. Теперь нужно вернуться в основной поток авторизации.
    - **Редирект**: `HTTP 302 Location: https://api.auth.oidc.com/authorize?request_id=<request_id>`

#### 2.4. Выдача Authorization Code

1.  **Описание**: `AuthBack` снова получает запрос на `/authorize`, но теперь с `user_session_id` cookie.
2.  **Логика `/authorize`**:
    - Извлекает `request_id`, находит `auth_request` в Redis.
    - Проверяет cookie, находит `user_session`, получает `user_id`.
    - _(Опционально: если клиент `third-party`, здесь может быть проверка и показ экрана Consent)_.
    - Генерирует `auth_code`.
    - **Сохраняет в Redis**: `auth_code:<code>` (одноразовый код, связанный с `user_id`, `client_id`, `code_challenge` и т.д.). TTL ~ 1 минута.[13]
3.  **Ответ**: Редирект браузера на `ClientFront` с кодом.
    - `HTTP 302 Location: https://app.client.com/auth/callback?code=...&state=...`

#### 2.5. Обмен кода на токены

1.  **Логика `ClientFront`**: Проверяет `state`. Отправляет `code` на свой бэкенд.
2.  **Запрос `ClientBack` -> `AuthBack`**:
    - **Тип**: `POST` (сервер-сервер)
    - **URL**: `https://api.auth.oidc.com/token`
    - **Тело (form-urlencoded)**: `grant_type=authorization_code&code=...&redirect_uri=...&client_id=...&client_secret=...&code_verifier=...`
3.  **Логика `AuthBack` (`/token`)**:
    - Аутентифицирует клиента по `client_id` и `client_secret`.
    - Находит и удаляет `auth_code` из Redis (он одноразовый).
    - **Проверяет PKCE**: хеширует `code_verifier` и сравнивает с `code_challenge`.
    - Генерирует `access_token`, `refresh_token`, `id_token`.
    - **Сохраняет в Redis**: `refresh_token:<jti>` для возможности отзыва.
4.  **Ответ `AuthBack` -> `ClientBack`**: `HTTP 200 OK` с JSON, содержащим все токены.

#### 2.6. Создание сессии в `ClientApp`

1.  **Логика `ClientBack`**:
    - Валидирует подпись, `iss`, `aud`, `exp`, `nonce` в `id_token`.
    - Создает собственную сессию, сохраняя в `client_cache` `access_token` и `refresh_token`.
2.  **Ответ `ClientBack` -> `ClientFront`**: `HTTP 200 OK` с `Set-Cookie: client_session_id=...`. Пользователь залогинен.

---

### Часть 3: Сценарии выхода (Logout)

##### Сценарий 1: Локальный выход (Revocation)

- **Назначение**: Выйти из приложения только на текущем устройстве.
- **Пример**: Выйти из "Магазина" в браузере, но остаться залогиненным на телефоне.
- **Поток**:
  1.  `User` нажимает "Выйти" в `ClientFront`.
  2.  `ClientFront` делает запрос на `ClientBack`: `POST /auth/logout`.
  3.  `ClientBack` извлекает `refresh_token` из своей сессии и отправляет запрос на `AuthBack`:
      - **Тип**: `POST`
      - **URL**: `POST https://api.auth.oidc.com/api/v1/revoke`
      - **Тело**: `token=<refresh_token>&client_id=...&client_secret=...`
  4.  `AuthBack` находит `refresh_token:<jti>` в Redis и удаляет его.
  5.  `ClientBack` удаляет свою сессию и очищает cookie у `ClientFront`.

##### Сценарий 2: Выход из всех сессий одного клиента (Client-Wide Logout)

- **Назначение**: Выйти из всех сессий конкретного приложения (например, "Магазина") на всех устройствах.
- **Пример**: Пользователь в `ClientFront` нажимает "Выйти со всех устройств Магазина".
- **Поток**:
  1.  `User` инициирует действие в `ClientFront`.
  2.  `ClientFront` делает запрос на `ClientBack`: `POST /auth/logout-all`.
  3.  `ClientBack` использует свой `access_token` для вызова `AuthBack`:
      - **Тип**: `POST`
      - **URL**: `POST https://api.auth.oidc.com/api/v1/logout-client-wide`
      - **Заголовки**: `Authorization: Bearer <access_token>`
  4.  `AuthBack` валидирует `access_token`, извлекает `sub` (UserID) и `aud` (ClientID), находит и удаляет все `refresh_token` для этой пары `(UserID, ClientID)`, инициируя Back-Channel Logout для каждого.

##### Сценарий 3: Удаленное завершение сессии (Remote Session Termination)

- **Назначение**: Управлять сессиями с центрального пульта.
- **Пример**: Пользователь на `AuthFrontPortal` видит список своих сессий и удаляет сессию "Магазин-Телефон".
- **Поток**:
  1.  `User` нажимает "Завершить" на портале `AuthFrontPortal`.
  2.  Портал отправляет запрос на `AuthBack` (аутентификация по cookie портала):
      - **Тип**: `DELETE`
      - **URL**: `DELETE https://api.auth.oidc.com/api/v1/users/me/sessions/<jti_сессии_телефона>`
  3.  `AuthBack` находит `UserID` по cookie, находит `refresh_token` по `jti`, проверяет, что он принадлежит этому пользователю, и удаляет его, инициируя Back-Channel Logout.

##### Сценарий 5: Глобальный выход (Global Logout)

- **Назначение**: Полностью завершить сеанс аутентификации, выйти отовсюду.
- **Пример**: Пользователь на `AuthFrontPortal` нажимает "Выйти со всех устройств".
- **Поток**:
  1.  Портал отправляет запрос на `AuthBack`:
      - **Тип**: `POST`
      - **URL**: `POST https://api.auth.oidc.com/api/v1/users/me/logout-global`
  2.  `AuthBack` находит `UserID` по cookie, находит **ВСЕ** `refresh_token`'ы и **ВСЕ** `user_session`'ы этого пользователя и удаляет их, рассылая уведомления Back-Channel Logout.

---

### Часть 4: Админские сценарии

##### 4.1. Ротация ключей (Key Rotation)

- **Назначение**: Регулярно менять ключи подписи токенов для повышения безопасности.
- **Поток**:
  1.  `System Admin` инициирует ротацию:
      - **Тип**: `POST`
      - **URL**: `POST https://api.auth.oidc.com/api/v1/admin/keys/rotate`
      - **Заголовки**: `Authorization: Bearer <SUPER_ADMIN_TOKEN>`
  2.  **Логика `AuthBack`**:
      - Генерирует новый ключ (`kid-2`), как в шаге 1.2, и делает его активным для подписи (`use: 'sig'`).
      - Старый ключ (`kid-1`) перестает быть активным для подписи, но остается в `jwks.json` для валидации уже выпущенных токенов.
      - Новые токены подписываются новым ключом `kid-2`.
      - Клиенты, получив токен с новым `kid`, обновляют свой кэш ключей с эндпоинта `jwks.json`.

##### 4.2. Аудит логов (Audit Logs)

- **Назначение**: Расследовать инциденты безопасности, отслеживать важные события.
- **Поток**:
  1.  `System Admin` запрашивает логи:
      - **Тип**: `GET`
      - **URL**: `GET https://api.auth.oidc.com/api/v1/admin/audit-logs?event=TokensRevoked&user_id=...`
      - **Заголовки**: `Authorization: Bearer <SUPER_ADMIN_TOKEN>`
  2.  **Логика `AuthBack`**:
      - На каждое важное событие (`UserLoggedIn`, `TokensIssued`, `TokensRevoked`, `BackchannelLogoutSent` и т.д.) `AuthBack` должен создавать запись в таблице `audit_logs`.[14]
      - Этот эндпоинт предоставляет доступ к этим записям с возможностью фильтрации.
  3.  **Ответ**: `HTTP 200 OK` с массивом записей аудита.
