# Wishlist API

## Описание

REST API сервис для создания и управления вишлистами.

Функционал:

* регистрация и авторизация пользователей (JWT)
* создание, редактирование и удаление вишлистов
* добавление и управление item'ами
* публичный доступ по ссылке
* бронирование подарков

---

## Запуск проекта

### 1. Клонировать репозиторий

```bash
git clone https://github.com/bagg1487/Wishlist_API/
cd Wishlist_API
```

### 2. Запуск через Docker

```bash
docker compose up --build
```

### 3. Swagger UI

Открыть в браузере:

```
http://localhost:8000/swagger/index.html
```

---

## Авторизация

Используется JWT токен.

После логина вставить токен в Swagger:

```
Authorize → Bearer <token>
```

---

## Основные эндпоинты

### Auth

* POST /auth/register
* POST /auth/login

### Wishlists

* GET /wishlists
* POST /wishlists
* GET /wishlists/{id}
* PUT /wishlists/{id}
* DELETE /wishlists/{id}

### Items

* POST /wishlists/{wishlistId}/items
* GET /wishlists/{wishlistId}/items
* PUT /items/{id}
* DELETE /items/{id}

### Public

* GET /public/{token}
* POST /public/{token}/book/{itemId}

---

##  Примеры запросов

Все примеры доступны через Swagger UI:
http://localhost:8000/swagger/index.html

---

##  Технологии

* Go (Golang)
* Gorilla Mux
* PostgreSQL
* GORM
* Docker
* Swagger (swaggo)

---

## Структура проекта

```
.
├── controllers/
├── models/
├── database/
├── middleware/
├── utils/
├── docs/
├── main.go
├── docker-compose.yml
└── README.md
```
