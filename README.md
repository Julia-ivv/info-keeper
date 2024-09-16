# info-keeper
___
Info-keeper - это клиент-серверная система, позволяющая пользователю хранить приватную информацию.

## Использование
____

### Запуск сервера
---
#### [Скачать](https://github.com/Julia-ivv/info-keeper/releases/tag/v1.0.0) исполняемый файл для своей ОС

#### Для запуска сервера с использованием файла конфигурации создать файл config.json следующей структуры

```json
{
   "database_dsn":"host=host port=port user=myuser password=xxxx dbname=mydb sslmode=disable",
   "grpc":":3200",
   "key":"abcdefg"
}
```
`database_dsn` - строка подключения к БД

`grpc` - порт для grpc

`key` - ключ для создания токена

Файл конфигурации должен находиться в директории с исполняемым файлом.
#### Для запуска сервера без использования файла конфигурации при запуске указать флаги:
```
    -g порт для grpc
    -d строка подключения к БД
    -k ключ для создания токена
```
#### или задать значения переменным окружения:
```
    GRPC_PORT порт для grpc
    DATABASE_DSN строка подключения к БД
    SKEY ключ для создания токена
```
#### Сервер использует БД PostgreSQL.
 Для установки PostgreSQL нужно [скачать](https://www.postgresql.org/download/) дистрибутив и запустить его на своей ОС.
 
 Для создания БД использовать команды консольного клиента psql:
 ```
    create database dbname;
    create user username with encrypted password 'userpassword';
    grant all privileges on database dbname to username;
```
Здесь `dbname` — это имя БД, `username` — имя пользователя, а `userpassword` — пароль пользователя.

Строка подключения к БД будеть иметь вид: 
    `host=host port=port user=myuser password=xxxx dbname=mydb sslmode=disable`

#### Пример запуска сервера:
с файлом конфигурации:
```
    keeperserver -c="./config.json"
```
с использованием флагов:
```   
    keeperserver -g ":3200" -d "host=host port=port user=myuser password=xxxx dbname=mydb sslmode=disabl" -k "key"
```
с использованием переменных окружения 
```
    keeperserver
``` 

### Запуск клиента
---
#### [Скачать](https://github.com/Julia-ivv/info-keeper/releases/tag/v1.0.0) исполняемый файл для своей ОС
#### Для запуска клиента с использованием файла конфигурации создать файл config.json следующей структуры
```json
{
   "database_uri":"keeper.db",
   "grpc":":3200"
}
```
`database_uri` - имя файла для БД

`grpc` - порт для grpc

Используется БД SQLite.

Файл конфигурации должен находиться в директории с исполняемым файлом.
#### Для запуска клиента без использования файла конфигурации при запуске указать флаги:
```
    -g порт для grpc
    -d имя файла для БД
```
#### или задать значения переменным окружения:
```
    GRPC_PORT порт для grpc
    DATABASE_NAME имя файла для БД
```

#### Пример запуска клиента:
с файлом конфигурации:
```
    keepercli -c="./config.json"
```
с использованием флагов:
```
    keepercli -g ":3200" -d "info.db"
```
с использованием переменных окружения 
```
    keepercli
```

Далее работать с клиентом, используя команды из документации или помощь -h.