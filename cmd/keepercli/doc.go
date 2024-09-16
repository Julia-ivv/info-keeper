/*
# Пакет main.

Точка входа в приложение клиента.

При запуске клиента необходимо указать файл конфигурации config.json.
Он должен находиться в директории с исполняемым файлом.

Пример файла конфигурации:

	{
	   	"database_uri":"keeper.db",
	   	"grpc":":3200"
	}

В параметре database_uri указывается имя файла БД.
БД будет создана в директории с исполняемым файлом.

В параметре grpc - порт gRPC сервера.

# Запуск клиента.

Скачайте исполняемый файл:

https://github.com/Julia-ivv/info-keeper/releases/tag/v1.0.0

Запустите клиент, указав файл конфигурации:

	keepercli -c="./config.json"

# Описание возможностей.

Клиент предоставляет следующие возможности:

  - регистрация, авторизация и аутентификация пользователя.
  - добавление информации в базу данных.
  - получение информации из базы данных.
  - обновление информации в базе данных.
  - синхронизация данных с сервером.

База данных может хранить следующую информацию:

  - реквизиты банковских карт.
  - пары логин - пароль.
  - текстовые данные.
  - бинарные данные из файлов.

Для каждого типа информации можно хранить описание или комментарий и
короткую подсказку (prompt).

При использовании клиента, информация изменяется в БД клиента.
Синхронизация с БД сервера происходит при аутентификации и при выходе из приложения.

# Использование клиента.

Сраузу после запуска для регистрации или аутентификации нужно ввести соответствующую команду:

	--reg -u=<user_name> // для регистрации нового пользователя.
	--auth -u=<user_name> // для аутентификации существующего пользователя.

После этого будет запрошен ввод пароля и секретного ключа для шифрования данных клиента.

Не забывайте свой пароль и ключ!

Значения всех флагов нужно указывать после знака "=".
Например,

	--gpwd -p=new login -l=login

Далее для работы используются флаги:

	--ncard
		Добавляет информацию о банковской карте.
		Используется с флагами -p -n -e -v -m.
		Например, --ncard -p=prompt -n=12345 -e=12/24 -v=111 -m=comment
	--npwd
		Добавляет информацию о паре логин-пароль.
		Используется с флагами -p -l -m.
		Например, --npwd -p=prompt -l=login -m=comment
	--ntext
		Добавляет текстовые данные.
		Используется с флагами -p -t -m.
		Например, --ntext -p=prompt -t=text -m=comment
	--nbyte
		Добавляет бинарные данные.
		Используется с флагами -p -b -m.
		Например, --nbyte -p=prompt -b=file -m=comment

	--ucard
		Обновляет информацию о банковской карте.
		Используется с флагами -p -n -e -v -m.
		Например, --ucard -p=prompt -n=12345 -e=12/24 -v=111 -m=comment
	--upwd
		Обновляет информацию о паре логин-пароль.
		Используется с флагами -p -l -m.
		Например, --upwd -p=prompt -l=login -m=comment
	--utext
		Обновляет текстовую информацию.
		Используется с флагами -p -t -m.
		Например, --utext -p=prompt -t=text -m=comment
	--ubyte
		Обновляет бинарную информацию.
		Используется с флагами -p -b -m.
		Например, --ubyte -p=prompt -b=file -m=comment

	--gcard
		Получает информацию о банковской карте по ее номеру.
		Используется с флагом -n.
		Например, --gcard -n=12345
	--gpwd
		Получает информацию о паре логин-пароль по подсказке и логину
		Используется с флагами -p -l.
		Например, --gpwd -p=prompt -l=login
	--gtext
		Получает текстовую информацию по подсказке.
		Используется с флагом -p.
		Например, --gtext -p=prompt
	--gbyte
		Получает бинарную информацию в файл по подсказке.
		Создает файл с именем подсказки в директории с исполняемым файлом.
		Используется с флагом -p.
		Например, --gbyte -p=prompt

	--gcards
		Получает информацию обо всех банковских картах пользователя.
		Используется без дополнительных флагов.
		Например, --gcards
	--gpwds
		Получает информацию обо всех парах логин-пароль пользователя.
		Используется без дополнительных флагов.
		Например, --gpwds
	--gtexts
		Получает всю текстовую информацию пользователя.
		Используется без дополнительных флагов.
		Например, --gtexts
	--gbytes
		Получает всю бинарную информацию пользователя.
		Данные сохраняются в разные файлы.
		Используется без дополнительных флагов.
		Например, --gbytes

	--fcard
		Обновляет данные банковской карты на сервере.
		Предназначена для разрешения конфликтов при синхронизации.
		Используется с флагом -n.
		Например, --fcard -n=12345
	--fpwd
		Обновляет данные пары логин-пароль на сервере.
		Предназначена для разрешения конфликтов при синхронизации.
		Используется с флагами -p -l.
		Например, --fpwd -p=prompt -l=login
	--ftext
		Обновляет текстовые данные на сервере.
		Предназначена для разрешения конфликтов при синхронизации.
		Используется с флагом -p.
		Например, --ftext -p=prompt
	--fbyte
		Обновляет бинарные данные на сервере.
		Предназначена для разрешения конфликтов при синхронизации.
		Используется с флагом -p.
		Например, --fbyte -p=prompt

	--scard
		Получает информацию с сервера о банковской карте по номеру.
		Предназначена для разрешения конфликтов при синхронизации.
		Используется с флагом -n.
		Например, --scard -n=12345
	--spwd
		Получает информацию с сервера о паре логин-пароль по подсказке и логину.
		Предназначена для разрешения конфликтов при синхронизации.
		Используется с флагами -p -l.
		Например, --spwd -p=prompt -l=login
	--stext
		Получает текстовую информацию с сервера по подсказке.
		Предназначена для разрешения конфликтов при синхронизации.
		Используется с флагом -p.
		Например, --stext -p=prompt
	--sbyte
		Получает бинарную информацию с сервера по подсказке в файл.
		Предназначена для разрешения конфликтов при синхронизации.
		Используется с флагом -p.
		Например, --sbyte -p=prompt

	-u
		Используется для указания логина пользователя при регистрации и аутентификации.
	-p
		Используется для указания короткой подсказки для данных.
	-l
		Используется для указания логина из пары логин-пароль.
	-m
		Используется для указания комментария или описания для данных.
	-n
		Используется для указания номера банковской карты.
	-e
		Используется для указания даты срока действия банковской карты.
	-v
		Используется для указания кода банковской карты.
	-t
		Используется для указания текстовых данных.
	-b
		Используется для указания пути к файлу с данными.

	-x
		Используется для выхода из приложения.
	--version
		Используется для получения информации о версии и дате сборки приложения.
	-h
		Получение помощи.
*/
package main
