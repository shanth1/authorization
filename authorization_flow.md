## Authorization Flow

#### Сценарий A: Первая авторизация (пользователь не залогинен)

```
┌──────────────┐
│ ClientFront  │
│ app.client.  │
│    com       │
└──────┬───────┘
       │ 1. User нажимает "Войти"
       │
       ↓ 2. Генерирует state, nonce, code_verifier, code_challenge
       │    Сохраняет в sessionStorage
       │
       ↓ 3. Редирект на AuthBack
       │    GET /authorize?
       │        response_type=code
       │        &client_id=...
       │        &redirect_uri=...
       │        &scope=openid profile
       │        &state=...
       │        &nonce=...
       │        &code_challenge=...
       │        &code_challenge_method=S256
       │
┌──────▼────────────────────────────────────────────────┐
│ AuthBack: GET /authorize                              │
│ Handler: OIDCHandler.Authorize                        │
│ Service: AuthorizationService.InitiateAuthorizationRequest │
└──────┬────────────────────────────────────────────────┘
       │
       │ 4. Валидация параметров:
       │    - client_id существует и активен
       │    - redirect_uri в списке разрешенных
       │    - scope в списке разрешенных
       │
       │ 5. Проверка cookie user_session_id
       │    ❌ Cookie отсутствует или невалидна
       │
       │ 6. Создание AuthorizationRequest:
       │    - Генерируется request_id
       │    - Сохраняется в Redis с TTL 10 минут
       │    - Поля: client_id, redirect_uri, scope,
       │             state, nonce, pkce, prompt, max_age
       │    - UserSessionID = nil (не залогинен)
       │    - ConsentGranted = false
       │
       │ 7. Логирование события: AuthorizationRequested
       │
       ↓ 8. HTTP 302 Redirect
       │    Location: https://sso.auth.oidc.com?request_id=<request_id>
       │
┌──────▼───────────┐
│ AuthFrontSSO     │
│ sso.auth.oidc.   │
│    com           │
└──────┬───────────┘
       │
       │ 9. Загружается страница SSO
       │
       ↓ 10. Запрос списка провайдеров
       │     GET /api/v1/auth/providers?request_id=<request_id>
       │
┌──────▼────────────────────────────────────────────────┐
│ AuthBack: GET /api/v1/auth/providers                  │
│ Handler: ProviderHandler.List                         │
└──────┬────────────────────────────────────────────────┘
       │
       │ 11. Получение AuthorizationRequest по request_id
       │ 12. Получение Client по client_id
       │ 13. Возврат client.AllowedProviders
       │
       ↓ HTTP 200 OK
       │ {"providers": ["password", "google"]}
       │
┌──────▼───────────┐
│ AuthFrontSSO     │
└──────┬───────────┘
       │
       │ 14. Рендер UI:
       │     - Форма email/password
       │     - Кнопка "Войти через Google"
       │
       │ 15. User вводит email + password
       │
       ↓ 16. POST /api/v1/auth/login/password
       │     Body: {
       │       "email": "user@example.com",
       │       "password": "secret",
       │       "request_id": "<request_id>"
       │     }
       │
┌──────▼────────────────────────────────────────────────┐
│ AuthBack: POST /api/v1/auth/login/password            │
│ Handler: UserHandler.LoginPassword                    │
│ Service: AuthenticationService.AuthenticateWithPassword │
└──────┬────────────────────────────────────────────────┘
       │
       │ 17. Валидация credentials:
       │     - Поиск пользователя по email
       │     - Проверка активности
       │     - Проверка хеша пароля
       │
       │ 18. Создание UserSession:
       │     Service: SessionService.CreateSession
       │     - Генерируется user_session_id
       │     - Сохраняется в Redis с TTL 30 дней
       │     - Поля: user_id, provider="password",
       │              ip_address, user_agent, created_at,
       │              expires_at, last_used_at
       │
       │ 19. Связывание сессии с запросом:
       │     Service: AuthorizationService.CompleteAuthentication
       │     - Обновляется AuthorizationRequest.UserSessionID
       │
       │ 20. Логирование события: UserLoggedIn
       │
       ↓ HTTP 200 OK
       │ Set-Cookie: user_session_id=<session_id>;
       │             HttpOnly; Secure; SameSite=Lax;
       │             Max-Age=2592000; Domain=.auth.oidc.com
       │
┌──────▼───────────┐
│ AuthFrontSSO     │
└──────┬───────────┘
       │
       │ 21. Получен 200 OK и Set-Cookie
       │
       ↓ 22. Редирект обратно на /authorize с request_id
       │     GET /authorize?request_id=<request_id>
       │
┌──────▼────────────────────────────────────────────────┐
│ AuthBack: GET /authorize (второй заход)               │
│ Handler: OIDCHandler.Authorize                        │
│ Service: AuthorizationService                         │
└──────┬────────────────────────────────────────────────┘
       │
       │ 23. Получение AuthorizationRequest по request_id
       │ 24. Проверка cookie user_session_id
       │     ✅ Cookie валидна, UserSession найдена
       │
       │ 25. Проверка необходимости consent:
       │     - Если client.Type == "first_party" → consent не нужен
       │     - Если client.RequiresConsent == false → consent не нужен
       │     - Иначе проверяем ConsentRepository.Exists(user_id, client_id, scopes)
       │
       │ Предположим: first_party клиент → consent не требуется
       │
       │ 26. Выдача Authorization Code:
       │     Service: AuthorizationService.IssueAuthorizationCode
       │     - Генерируется code (60 сек TTL)
       │     - Сохраняется AuthorizationCode в Redis:
       │       {code, request_id, user_session_id, user_id,
       │        client_id, redirect_uri, scope, nonce, pkce}
       │     - Удаляется AuthorizationRequest из Redis
       │
       │ 27. Логирование события: AuthorizationGranted
       │
       ↓ HTTP 302 Redirect
       │ Location: https://app.client.com/auth/callback?
       │           code=<code>&state=<state>
       │
┌──────▼───────────┐
│ ClientFront      │
└──────┬───────────┘
       │
       │ 28. Проверка state из sessionStorage
       │ 29. Извлечение code_verifier из sessionStorage
       │
       ↓ 30. Отправка code на свой бэкенд
       │     POST /auth/callback
       │     Body: {"code": "...", "code_verifier": "..."}
       │
┌──────▼───────────┐
│ ClientBack       │
└──────┬───────────┘
       │
       ↓ 31. POST /token на AuthBack
       │     Body (application/x-www-form-urlencoded):
       │         grant_type=authorization_code
       │         &code=<code>
       │         &redirect_uri=<redirect_uri>
       │         &client_id=<client_id>
       │         &client_secret=<client_secret>
       │         &code_verifier=<code_verifier>
       │
┌──────▼────────────────────────────────────────────────┐
│ AuthBack: POST /token                                 │
│ Handler: OIDCHandler.Token                            │
│ Service: TokenService.ExchangeCodeForTokens           │
└──────┬────────────────────────────────────────────────┘
       │
       │ 32. Аутентификация клиента:
       │     - Проверка client_id + client_secret
       │
       │ 33. Валидация кода:
       │     - Code существует и не expired
       │     - Code не использован (Used == false)
       │     - Code.ClientID == client_id
       │     - Code.RedirectURI == redirect_uri
       │
       │ 34. Проверка PKCE:
       │     - Хеширование code_verifier (SHA256)
       │     - Сравнение с code_challenge из кода
       │
       │ 35. Пометка кода как использованного
       │
       │ 36. Генерация токенов:
       │     - AccessToken (15 мин TTL):
       │       {id, user_id, client_id, scope, jti, expires_at}
       │     - RefreshToken (30 дней TTL):
       │       {id, user_id, client_id, user_session_id,
       │        scope, jti, expires_at, revoked=false}
       │     - IDToken (JWT, подписанный):
       │       {iss, sub, aud, exp, iat, auth_time,
       │        nonce, email, name, picture}
       │
       │ 37. Сохранение токенов в Redis
       │ 38. Логирование события: TokensIssued
       │
       ↓ HTTP 200 OK
       │ Content-Type: application/json
       │ {
       │   "access_token": "<jwt_or_jti>",
       │   "token_type": "Bearer",
       │   "expires_in": 900,
       │   "refresh_token": "<jti>",
       │   "id_token": "<jwt>"
       │ }
       │
┌──────▼───────────┐
│ ClientBack       │
└──────┬───────────┘
       │
       │ 39. Валидация ID Token:
       │     - Проверка подписи (jwks.json)
       │     - Проверка iss, aud, exp, nonce
       │
       │ 40. Создание собственной сессии:
       │     - Сохранение access_token и refresh_token
       │     - Генерация client_session_id
       │
       ↓ HTTP 200 OK
       │ Set-Cookie: client_session_id=...
       │
┌──────▼───────────┐
│ ClientFront      │
└──────────────────┘

       🎉 Пользователь залогинен!
```

---

#### Сценарий B: Повторная авторизация (пользователь уже залогинен в SSO)

```
┌──────────────┐
│ ClientFront  │
└──────┬───────┘
       │ 1. User нажимает "Войти" в другом клиенте
       │
       ↓ 2. Редирект на AuthBack
       │    GET /authorize?...
       │    Cookie: user_session_id=<existing_session>
       │
┌──────▼────────────────────────────────────────────────┐
│ AuthBack: GET /authorize                              │
│ Service: AuthorizationService.InitiateAuthorizationRequest │
└──────┬────────────────────────────────────────────────┘
       │
       │ 3. Валидация параметров (как в Сценарии A)
       │
       │ 4. Проверка cookie user_session_id
       │    ✅ Cookie найдена и валидна
       │
       │ 5. Получение UserSession из Redis
       │    - Проверка expires_at
       │    - Обновление last_used_at
       │
       │ 6. Создание AuthorizationRequest с заполненным UserSessionID
       │
       │ 7. Проверка prompt и max_age:
       │    - Если prompt=login → требуется повторная аутентификация
       │    - Если max_age превышен → требуется повторная аутентификация
       │    - Иначе → аутентификация не требуется
       │
       │ Предположим: повторная аутентификация НЕ требуется
       │
       │ 8. Проверка необходимости consent:
       │    - Проверяется наличие ConsentGrant
       │    - Если consent есть для всех запрошенных scopes → consent не нужен
       │
       │ Предположим: consent не требуется (first_party или уже дан)
       │
       │ 9. Немедленная выдача Authorization Code:
       │    (пропускаем редирект на SSO)
       │    Service: AuthorizationService.IssueAuthorizationCode
       │
       ↓ HTTP 302 Redirect
       │ Location: https://app.client.com/auth/callback?
       │           code=<code>&state=<state>
       │
       (Дальше как в Сценарии A, шаги 28-40)
```

**Ключевое отличие:** Пользователь не видит страницу SSO, авторизация происходит мгновенно благодаря SSO-сессии.

---

#### Сценарий C: Авторизация с consent screen (third-party клиент)

```
(Шаги 1-25 как в Сценарии A)

┌──────────────────────────────────────────────────────┐
│ AuthBack: GET /authorize (после аутентификации)      │
└──────┬───────────────────────────────────────────────┘
       │
       │ 25. Проверка необходимости consent:
       │     - client.Type == "third_party"
       │     - client.RequiresConsent == true
       │     - ConsentRepository.Exists() == false
       │
       │ ❌ Consent требуется!
       │
       ↓ HTTP 302 Redirect
       │ Location: https://portal.auth.oidc.com/consent?
       │           request_id=<request_id>
       │
┌──────▼──────────────┐
│ AuthFrontPortal     │
│ (Consent Screen)    │
└──────┬──────────────┘
       │
       │ 26. GET /api/v1/auth/consent-info?request_id=...
       │     (получить информацию о клиенте и запрошенных scopes)
       │
       │ 27. Рендер UI:
       │     "Приложение 'Shop XYZ' запрашивает доступ к:"
       │     - Ваш email
       │     - Ваше имя
       │     [Разрешить] [Отказать]
       │
       │ 28. User нажимает "Разрешить"
       │
       ↓ POST /api/v1/auth/consent
       │ Body: {
       │   "request_id": "...",
       │   "granted": true,
       │   "scopes": ["openid", "profile", "email"]
       │ }
       │
┌──────▼────────────────────────────────────────────────┐
│ AuthBack: POST /api/v1/auth/consent                   │
│ Handler: AuthorizationHandler.SubmitConsent (NEW)     │
│ Service: AuthorizationService.CompleteConsent         │
└──────┬────────────────────────────────────────────────┘
       │
       │ 29. Получение AuthorizationRequest
       │ 30. Обновление AuthorizationRequest.ConsentGranted = true
       │ 31. Сохранение ConsentGrant в БД
       │ 32. Логирование события: ConsentGranted
       │
       ↓ HTTP 200 OK
       │
┌──────▼──────────────┐
│ AuthFrontPortal     │
└──────┬──────────────┘
       │
       ↓ 33. Редирект на /authorize?request_id=...
       │
┌──────▼────────────────────────────────────────────────┐
│ AuthBack: GET /authorize (третий заход)               │
└──────┬────────────────────────────────────────────────┘
       │
       │ 34. Проверка: UserSessionID ✅, ConsentGranted ✅
       │
       │ 35. Выдача Authorization Code
       │
       ↓ HTTP 302 Redirect → ClientFront
       │
       (Дальше как в Сценарии A)
```

---

#### Сценарий D: Вход через Google

```
┌──────────────┐
│ AuthFrontSSO │
└──────┬───────┘
       │ (После шагов 1-14 из Сценария A)
       │
       │ 15. User нажимает "Войти через Google"
       │
       ↓ 16. Генерация OAuth state для Google
       │     Сохранение в sessionStorage: {request_id, google_state}
       │
       │ 17. Редирект на Google OAuth
       │     Location: https://accounts.google.com/o/oauth2/v2/auth?
       │         client_id=<google_client_id>
       │         &redirect_uri=https://api.auth.oidc.com/api/v1/auth/providers/google/callback
       │         &response_type=code
       │         &scope=openid email profile
       │         &state=<google_state>
       │
┌──────▼───────────┐
│ Google OAuth     │
└──────┬───────────┘
       │ 18. User аутентифицируется в Google
       │ 19. Google редиректит обратно
       │
       ↓ GET /api/v1/auth/providers/google/callback?
       │     code=<google_code>&state=<google_state>
       │
┌──────▼────────────────────────────────────────────────┐
│ AuthBack: GET /auth/providers/google/callback         │
│ Handler: ProviderHandler.GoogleCallback               │
│ Service: AuthenticationService.AuthenticateWithProvider │
└──────┬────────────────────────────────────────────────┘
       │
       │ 20. Валидация state
       │ 21. Обмен google_code на Google токены
       │ 22. Получение Google UserInfo
       │ 23. Создание ProviderProfile:
       │     {provider: "google", subject: "<google_sub>",
       │      email, name, picture}
       │
       │ 24. Поиск/создание User:
       │     - UserRepository.GetByProviderID("google", google_sub)
       │     - Если нет → создать нового пользователя
       │     - Если есть → использовать существующего
       │
       │ 25. Создание UserSession:
       │     SessionService.CreateSession(user_id, "google", ...)
       │
       │ 26. Связывание с AuthorizationRequest:
       │     (необходимо передать request_id через state или cookie)
       │     AuthorizationService.CompleteAuthentication(request_id, user_session_id)
       │
       │ 27. Логирование события: UserLoggedIn
       │
       ↓ HTTP 302 Redirect
       │ Set-Cookie: user_session_id=<session_id>
       │ Location: https://api.auth.oidc.com/authorize?request_id=<request_id>
       │
       (Дальше как в Сценарии A, шаги 23-40)
```

---

### ЧАСТЬ 5: Сценарии Logout (детализация)

#### Сценарий Logout 1: Локальный выход (Local Logout)

**Цель:** Выйти из конкретного клиента на конкретном устройстве.

```
┌──────────────┐
│ ClientFront  │
└──────┬───────┘
       │ 1. User нажимает "Выйти"
       │
       ↓ POST /auth/logout (к своему бэкенду)
       │
┌──────▼───────────┐
│ ClientBack       │
└──────┬───────────┘
       │
       │ 2. Извлечение refresh_token из сессии
       │
       ↓ POST /revoke (к AuthBack)
       │ Body: token=<refresh_token>&client_id=...&client_secret=...
       │
┌──────▼────────────────────────────────────────────────┐
│ AuthBack: POST /revoke                                │
│ Handler: OIDCHandler.Revoke                           │
│ Service: TokenService.RevokeToken                     │
└──────┬────────────────────────────────────────────────┘
       │
       │ 3. Аутентификация клиента
       │ 4. Отзыв refresh_token по JTI:
       │    - Пометка Revoked = true
       │    - Или удаление из Redis
       │
       │ 5. Логирование: TokensRevoked
       │
       ↓ HTTP 200 OK
       │
┌──────▼───────────┐
│ ClientBack       │
└──────┬───────────┘
       │
       │ 6. Удаление своей сессии
       │ 7. Очистка cookie client_session_id
       │
       ↓ HTTP 200 OK
       │ Set-Cookie: client_session_id=; Max-Age=0
       │
┌──────▼──────────┐
│ ClientFront     │
└─────────────────┘

       ✅ Выход выполнен локально
       ⚠️  SSO-сессия все еще активна!
```

---

#### Сценарий Logout 2: Client-Wide Logout

**Цель:** Выйти из всех сессий конкретного клиента на всех устройствах.

```
┌──────────────┐
│ ClientFront  │
└──────┬───────┘
       │ 1. User нажимает "Выйти со всех устройств"
       │
       ↓ POST /auth/logout-all
       │
┌──────▼───────────┐
│ ClientBack       │
└──────┬───────────┘
       │
       │ 2. Извлечение access_token
       │
       ↓ POST /logout-client-wide (к AuthBack)
       │ Headers: Authorization: Bearer <access_token>
       │
┌──────▼────────────────────────────────────────────────┐
│ AuthBack: POST /logout-client-wide                    │
│ Handler: OIDCHandler.ClientWideLogout                 │
│ Service: LogoutService.ClientWideLogout               │
└──────┬────────────────────────────────────────────────┘
       │
       │ 3. Валидация access_token (Bearer middleware)
       │ 4. Извлечение user_id и client_id из токена
       │
       │ 5. Получение всех refresh_tokens для (user_id, client_id)
       │ 6. Отзыв всех этих токенов
       │
       │ 7. Для каждого отозванного токена:
       │    - Получение Client
       │    - Если BackchannelLogoutURI != nil:
       │      * Создание LogoutToken (JWT)
       │      * POST на backchannel_logout_uri клиента
       │
       │ 8. Логирование: TokensRevoked + BackchannelLogoutSent
       │
       ↓ HTTP 200 OK
       │
┌──────▼───────────┐
│ ClientBack       │
└──────────────────┘

       В фоне:

┌─────────────────────────────────────────────────────┐
│ POST https://api.client.com/logout/oidc             │
│ (Back-Channel Logout Endpoint у ClientBack)         │
│ Body: logout_token=<jwt>                            │
└──────┬──────────────────────────────────────────────┘
       │
┌──────▼───────────┐
│ ClientBack       │
│ (другое          │
│  устройство)     │
└──────┬───────────┘
       │
       │ Валидация logout_token:
       │ - Проверка подписи
       │ - Проверка events["backchannel-logout"]
       │ - Извлечение sub (user_id)
       │
       │ Удаление всех сессий этого пользователя
       │
       ↓ HTTP 204 No Content

       ✅ Все устройства получили уведомление о выходе
```

---

#### Сценарий Logout 3: Remote Session Termination

**Цель:** Завершить конкретную сессию из портала управления.

```
┌──────────────────┐
│ AuthFrontPortal  │
└──────┬───────────┘
       │ 1. User открывает список сессий
       │
       ↓ GET /api/v1/users/me/sessions
       │ Cookie: user_session_id=<portal_session>
       │
┌──────▼────────────────────────────────────────────────┐
│ AuthBack: GET /users/me/sessions                      │
│ Handler: UserHandler.GetMySessions                    │
│ Service: SessionService.ListUserSessions              │
└──────┬────────────────────────────────────────────────┘
       │
       │ 2. Middleware UserSession извлекает user_id из cookie
       │ 3. Получение всех UserSession пользователя
       │ 4. Для каждой сессии - получение связанных refresh_tokens
       │    (чтобы показать, в каких клиентах активна сессия)
       │
       ↓ HTTP 200 OK
       │ [
       │   {
       │     "session_id": "sess_123",
       │     "device": "Chrome on Mac",
       │     "ip": "1.2.3.4",
       │     "created_at": "...",
       │     "last_used_at": "...",
       │     "clients": ["Shop", "Forum"]
       │   },
       │   ...
       │ ]
       │
┌──────▼──────────────┐
│ AuthFrontPortal     │
└──────┬──────────────┘
       │
       │ 5. User нажимает "Завершить" у сессии "sess_123"
       │
       ↓ DELETE /api/v1/users/me/sessions/sess_123
       │ Cookie: user_session_id=<portal_session>
       │
┌──────▼────────────────────────────────────────────────┐
│ AuthBack: DELETE /users/me/sessions/:jti              │
│ Handler: UserHandler.TerminateSession                 │
│ Service: LogoutService.RemoteSessionTermination       │
└──────┬────────────────────────────────────────────────┘
       │
       │ 6. Проверка владельца сессии (user_id)
       │ 7. Получение всех refresh_tokens для session_id
       │ 8. Удаление UserSession
       │ 9. Отзыв всех refresh_tokens этой сессии
       │
       │ 10. Для каждого клиента:
       │     - Отправка Back-Channel Logout
       │
       │ 11. Логирование: SessionTerminated + BackchannelLogoutSent
       │
       ↓ HTTP 204 No Content

       ✅ Сессия завершена, все клиенты уведомлены
```

---

#### Сценарий Logout 5: Global Logout

**Цель:** Выйти отовсюду - завершить все сессии и все клиенты.

```
┌──────────────────┐
│ AuthFrontPortal  │
└──────┬───────────┘
       │ 1. User нажимает "Выйти со всех устройств"
       │
       ↓ POST /api/v1/users/me/logout-global
       │ Cookie: user_session_id=<portal_session>
       │
┌──────▼────────────────────────────────────────────────┐
│ AuthBack: POST /users/me/logout-global                │
│ Handler: UserHandler.GlobalLogout                     │
│ Service: LogoutService.GlobalLogout                   │
└──────┬────────────────────────────────────────────────┘
       │
       │ 2. Middleware извлекает user_id
       │
       │ 3. Получение всех refresh_tokens пользователя
       │ 4. Группировка по client_id
       │
       │ 5. Удаление всех UserSession:
       │    SessionService.GlobalLogout(user_id)
       │
       │ 6. Отзыв всех refresh_tokens:
       │    TokenRepository.RevokeAllByUser(user_id)
       │
       │ 7. Для каждого уникального client_id:
       │    - Получение Client
       │    - Если BackchannelLogoutURI != nil:
       │      * Создание LogoutToken (без session_id, т.к. все сессии)
       │      * Отправка Back-Channel Logout
       │
       │ 8. Логирование: UserLoggedOut (global)
       │
       ↓ HTTP 200 OK
       │ Set-Cookie: user_session_id=; Max-Age=0

       ✅ Все сессии завершены
       ✅ Все клиенты уведомлены
       ✅ Portal cookie удалена
```
