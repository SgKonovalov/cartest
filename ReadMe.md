Задача:
Разработать CRUD (REST API) для модели автомобиля, который имеет следующие поля:
1. Уникальный идентификатор (любой тип, общение с БД не является критерием чего-либо,
можно сделать и in-memory хранилище на время жизни сервиса)
2. Бренд автомобиля (текст)
3. Модель автомобиля (текст)
4. Цена автомобиля (целое, не может быть меньше 0)
5. Статус автомобиля (В пути, На складе, Продан, Снят с продажи)
6. Пробег (целое)
Формат ответа api - json api (https://jsonapi.org/)

Решение:
Приложение состоит из 2-х микросервисов:
1) Однозадачный (JOB): добавляет данные обо всех автомобилях из реляционной БД в Redis.
Цель: снизить нагрузку на сервер БД и увеличить скорость обработки запроса на получение данных об автомобиле
(пакет loadallcars) – далее МС1.
2) Реализует CRUD операции и является http-сервером для обработки запросов – далее МС2.
* для реализации HTTP-сервера, использован фреймворк Gin$
* для работы с БД использован драйвер pgx.

МС1 – состоит из 1 файла, содержащего методы обработки данных
(подробное описание функционала в файле ./loadallcars/main.go).

МС2 – состоит из пакетов:
1. .migration – все sql-скрипты, используемые в приложении;
2. definition – само понятие «автомобиль» в приложении (структура Car) + геттеры/сеттеры и конструкторы указанной структуры;
3. handlers – хендлеры и модели JSON-ответов (подробное описание функционала в файле ./handlers/handlers.go
и ./handlers/responsemodel.go);
4. repository – структуры и функции, отвечающие за работу с БД и Redis: выполнение CRUD запросов
(подробное описание функционала в файле ./repository/repository.go). На все функции пакеты написаны unit-тесты
(./repository/repository_test.go).
5. service – структуры и функции, отвечающие за вызов методов пакета repository, но с пост- или предобработкой результатов/аргументов для функций  repository (подробное описание функционала в файле ./service/service.go). На дополнительный функционал, не охватывающий функции repository написаны unit-тесты (./service/service_test.go).
6. Файл main.go - инициализация структур и запуск приложения.