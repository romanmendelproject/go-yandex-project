# go-yandex-project

Репозитория для системы хранения приватных данных

## Использование файла конфигурации сервера

**max_cpu** - максимальное количство используемых ЦПУ

**log_level** - уровень логирования

**db_ip** - адрес коллектора api_json

**db_dsn** - формат отправляемых данных(netflow|json)


## Запуск программы
1. Клонируем репозиторий и переходим в него
2. Запускаем БД
- cd docker
- docker-compose up -d
3. Запуск сервера
- go run cmd/server/main.go

## Запуск тестов
1. Клонируем репозиторий и переходим в него
2. Запускаем БД
- cd docker
- docker-compose up -d
3. go test ./...  -coverprofile cover.out 
4. go tool cover -func cover.out

## Аргументы клиентской части
***-a*** Тип действия которое наобходимо выполнить:
- ***reg*** - регистрация новго пользователя;
- ***auth*** - авторизация сущесивующего пользователя; 
- ***set*** - отправка данных на сервер;
- ***get*** - получение данных с сервера;
***-t*** Тип хранимых данных:
- ***cred*** - пара логин/пароль;
- ***text*** - произвольные текстовые данные; 
- ***byte*** - произвольные бинарные данные;
- ***card*** -данные банковских карт;

***-n*** Наименование хранимой единицы

***-u*** Имя пользователя(может использоваться как для авторизации на сервера, так и для типа cred)

***-p*** Пароль пользователя(может использоваться как для авторизации на сервера, так и для типа cred)

***-d*** Набор хранимых данных(используется для типов text, byte, card)

***-m*** Мета информация для любого из типа данных

## Примеры использования
1. Регистрация нового пользователя
   
***go run cmd/client/main.go -a reg  -u user -p password***

2. Авторизация пользователя

***go run cmd/client/main.go -a auth  -u user -p password***

3. Записать на сервер данные формата cred(пара логин/пароль) с именем creddata

***go run cmd/client/main.go -a set -t cred -n creddata -u newuser -p newpassword***

4. Получить данные формата cred с именем creddata
   
***go run cmd/client/main.go -a get -t cred -n creddata***

5. Записать на сервер данные формата text с именем textdata
   
***go run cmd/client/main.go -a set -t text -n textdata -d data***

6. Получить данные формата text с именем textdata

***go run cmd/client/main.go -a get -t text -n textdata***

7. Записать на сервер данные формата byte с именем bytedata

***go run cmd/client/main.go -a set -t byte -n bytedata -d data***

8. Получить данные формата cred с именем bytedata

***go run cmd/client/main.go -a get -t byte -n bytedata***

9. Записать на сервер данные формата card с именем carddata
    
***go run cmd/client/main.go -a set -t card -n carddata -d 5555555555555555***

10. Получить данные формата cred с именем carddata
    
***go run cmd/client/main.go -a get -t card -n carddata***
