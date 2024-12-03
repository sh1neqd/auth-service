# auth-service

Небольшой сервис аутентификации

### Stack
  * Golang
  * Fiber
  * JWT
  * Postgres
  * Docker

### Возможные действия:
  * Генерация access и refresh токенов
  * Refresh операция на access и refresh токены

### Запуск проекта
  Необходимо склонировать репозиторий:
  ```
     git clone github.com/sh1neqd/auth-service
     git checkout main
```
  Далее запускаем проект с помощью docker-compose:
```
  docker-compose build
  docker-compose up
```
  
### Routes

* http://localhost:8000/token [GET]:  
Input parameter:
 ```user_id=UUID```  
Output:  
```
{
    "access_token" (UUID)
    "refresh_token" (UUID)
}
```

* http://localhost:8000/refresh [POST]:  
Input:  
```
{
    "access_token" (UUID)
    "refresh_token" (UUID)
}
```  
Output:  
```
{
    "access_token" (UUID)
    "refresh_token" (UUID)
}
```
