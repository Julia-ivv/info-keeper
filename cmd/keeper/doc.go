/*
# Пакет main.

Точка входа в приложение сервера.

При запуске сервера нужно указать файл конфигурации config.json.
Он должен находиться в директории с исполняемым файлом.

Пример файла конфигурации:

	{
	    "database_dsn":"",
	    "grpc":":3200",
	    "key":"byrhtvtyn"
	}

В параметре database_dsn указывается строка подключения к БД.
В параметре grpc - порт для grpc.
Параметр key - ключ для создания токена.

Используемая БД - PostgreSQL.

# Запуск сервера.

Скачайте исполняемый файл:

https://github.com/Julia-ivv/info-keeper/releases/tag/v1.0.0-server

Запустите сервер:

	keeperserver -c="./config.json"

# Описание возможностей.

Сервер предоставляет следующие возможности:

  - регистрация, авторизация и аутентификация пользователя.
  - добавление информации в базу данных.
  - получение информации из базы данных.
  - обновление информации в базе данных.
  - синхронизация данных с клиентом.

База данных может хранить следующую информацию:

  - реквизиты банковских карт.
  - пары логин - пароль.
  - текстовые данные.
  - бинарные данные.
*/
package main
