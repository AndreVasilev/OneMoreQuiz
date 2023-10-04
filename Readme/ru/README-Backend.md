# One more Quiz! - Backend

- [Требования](#requirements)
  - [Создайте своего Telegram бота](#create-your-telegram-bot)
  - [Создайте свою базу данных](#create-your-database)
  - [Локальное окружение](#local-environment)
- [Сборка и запуск](#build-and-run)
- [Дебаг](#debug)
- [Деплой](#deploy)
- [Альтернативные опции](#alternative-options)

## Требования

### Создайте своего Telegram бота
Чтобы запустить свой Mini App вам понадобится Telegram бот. Поговорите с [@BotFather](https://t.me/botfather) чтобы создать своего бота и [получить токен](https://core.telegram.org/bots/tutorial#obtain-your-bot-token)

### Создайте свою базу данных
1. Этот бэкенд использует [Managed Service for YDB](https://cloud.yandex.com/en/services/ydb). Вы можете зарегистрироваться в [Yandex Cloud](https://cloud.yandex.com/) и [создать](https://cloud.yandex.com/en/docs/ydb/quickstart) свою базу данных. 
Или вы можете использовать любую другую базу данных, но тогда вам придется реализовать свой собственный репозиторий с таким же [интерфейсом](https://github.com/AndreVasilev/OneMoreQuiz/tree/readme/yc-serverless-backend/repository/Interface.go).

[Создайте сервисный аккаунт](https://cloud.yandex.com/en/docs/iam/operations/sa/create) с ролью ```ydb.editor``` и сгенерируйте authorized_key.json. 
Положите authorized_key.json в корневую папку репозитория. Или измените путь к вамшему authorized_key.json в файле [YDB.go](https://github.com/AndreVasilev/OneMoreQuiz/tree/readme/yc-serverless-backend/repository/YDB.go) на строке 23:
```go
func authCredentials() ydb.Option {
	if os.Getenv("DEBUG") == "TRUE" {
		return yc.WithServiceAccountKeyFileCredentials("../authorized_key.json")
	} else {
		return yc.WithMetadataCredentials()
	}
}
```

3. Создайте две таблицы

<details>
  <summary>Question</summary>
  
  | Свойство | Тип |
  | ------ | ------ |
  | id | Uint64 |
  | A | String |
  | B | String |
  | C | String |
  | D | String |
  | answer | String |
  | question | String |
</details>

<details>
  <summary>User</summary>
  
  | Свойство | Тип |
  | ------ | ------ |
  | id | int64 |
  | last_question_id | Uint64 |
  | score | Uint64 |
  | success_answers | Uint64 |
  | tg_data | String |
</details>

3. Заполните таблицу Questionданными, используя [sql скрипт](https://github.com/AndreVasilev/OneMoreQuiz/blob/readme/yc-serverless-backend/repository/init_001.sql)
<details>
  <summary>Пример sql скрипта</summary>
  
  ```sql
  UPSERT INTO `question`
    ( `id`, `question`, `A`, `B`, `C`, `D`, `answer` )
  VALUES (1, "A knish is traditionally stuffed with what filling?", "potato", "creamed corn", "lemon custard", "raspberry jelly", "A"),
  ...;
  ```
</details>

### Локальное окружение

Обновите или установите необходимое программное обеспечение
- [Golang](https://go.dev/doc/install) >= 1.19

Установите переменные окружения
```sh
$ export VAR=abc
```
| Переменная | Значение |
| ------ | ------ |
| BOT_TOKEN | [Получите токен бота](https://core.telegram.org/bots/tutorial#obtain-your-bot-token) |
| DEBUG | TRUE |
| YDB_CONNECTION_URL | [Получите url подключения](https://cloud.yandex.com/en/docs/ydb/operations/connection#endpoint-and-path) к вашей YDB |

## Сборка и запуск

Откройте терминал, перейдите в каталог проекта backends, обновите все модули и запустите проект
```
cd /path/to/repo/dir/yc-serverless-backend/
go get -u
go mod tidy
go run .
```

## Дебаг

Напишите свой код для отладки любого метода в функции main()

```go
func main() {
  questions := repository.GetQuestions(0)
  log.Printf("Questions total: %d", len(questions))
}
```
Помните, что эти две функции должны вызываться службой облачных функций

```go
func QuestionHandler(rw http.ResponseWriter, req *http.Request) {
  handlers.Question(rw, req)
}

func UserHandler(rw http.ResponseWriter, req *http.Request) {
  handlers.User(rw, req)
}
```

## Деплой

Исходный код этого бэкенда предназначен для развертывания в [Yandex Cloud Functions](https://cloud.yandex.com/en/services/functions). Но его можно развернуть в другой службе, такой как [AWS Lambda](https://aws.amazon.com/lambda/) или [Google Cloud Functions](https://cloud.google.com/functions) (и любой другой бессерверный сервис, поддерживающий среду выполнения Go 1.19).

1. [Создайте](https://cloud.yandex.com/en/docs/functions/quickstart/create-function/go-function-quickstart) две функции для адресов /question и /user.
   Получите ID и имя каждой функции

2. Установите [Yandex Cloud CLI](https://cloud.yandex.com/en/docs/cli/operations/install-cli) и авторизуйтесь в своем аккаунте

3. Установите переменные окружения
```sh
$ export VAR=abc
```
| Variable | Value |
| ------ | ------ |
| QUESTION_FUNCTION_ID | ID функции для адреса Question, созданная на шаге 1 |
| QUESTION_FUNCTION_NAME | Имя функции для адреса Question, созданная на шаге  1 |
| USER_FUNCTION_ID | ID функции для адреса User, созданная на шаге  1 |
| USER_FUNCTION_NAME | Имя функции для адреса User, созданная на шаге  1 |
| SERVICE_ACCOUNT_ID | [ID сервисного аккаунта](#create-your-database) for editing database |
| BOT_TOKEN | [Получите токен бота](#create-your-telegram-bot) |
| YDB_CONNECTION_URL | [Получите url подключения](#create-your-database) к вашей YDB |

```sh
cd /path/to/repo/dir/yc-serverless-backend/
sh ./scripts/deploy.sh
```

## Альтернативные опции

1. Разверните свой код вручную с помощью [Yandex Cloud Console](https://cloud.yandex.com/en/docs/functions/quickstart/create-function/go-function-quickstart) или другими способами

2. Сгенерируйте любой бэкенд, который вам нравится, с помощью опции Generate Client от [Swagger Editor](https://editor.swagger.io/), используя [openapi.yaml](https://github.com/AndreVasilev/OneMoreQuiz/blob/readme/openapi.yaml) и разверните его там, где захотите :)
