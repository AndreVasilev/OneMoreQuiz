# One more Quiz! - Backend

- [Requirements](#requirements)
  - [Create your Telegram Bot](#create-your-telegram-bot)
  - [Create your database](#create-your-database)
  - [Local environment](#local-environment)
- [Build and Run](#build-and-run)
- [Debug](#debug)
- [Deploy](#deploy)
- [Alternative options](#alternative-options)

## Requirements

### Create your Telegram Bot
To run your own Mini App you will need a Telegram bot. Talk to [@BotFather](https://t.me/botfather) to create your bot and [obtain your bot's token](https://core.telegram.org/bots/tutorial#obtain-your-bot-token)

### Create your database
1. The default backend uses [Managed Service for YDB](https://cloud.yandex.com/en/services/ydb). You can register at [Yandex Cloud](https://cloud.yandex.com/) and [create](https://cloud.yandex.com/en/docs/ydb/quickstart) your own database. 
Or you can use any other database, but in that case you will need to implement your own repository with same [interface](https://github.com/AndreVasilev/OneMoreQuiz/tree/readme/yc-serverless-backend/repository/Interface.go).

[Create a service account](https://cloud.yandex.com/en/docs/iam/operations/sa/create) with role ```ydb.editor``` and generate authorized_key.json. 
Put the authorized_key.json to root of this repository. Or change a path to your authorized_key.json at [YDB.go](https://github.com/AndreVasilev/OneMoreQuiz/tree/readme/yc-serverless-backend/repository/YDB.go) line 23:
```go
func authCredentials() ydb.Option {
	if os.Getenv("DEBUG") == "TRUE" {
		return yc.WithServiceAccountKeyFileCredentials("../authorized_key.json")
	} else {
		return yc.WithMetadataCredentials()
	}
}
```

3. Create two tables

<details>
  <summary>Question</summary>
  
  | Variable | Type |
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
  
  | Variable | Type |
  | ------ | ------ |
  | id | int64 |
  | last_question_id | Uint64 |
  | score | Uint64 |
  | success_answers | Uint64 |
  | tg_data | String |
</details>

3. Fill your Question table with data using [sql script](https://github.com/AndreVasilev/OneMoreQuiz/blob/readme/yc-serverless-backend/repository/init_001.sql)
<details>
  <summary>Sql script exmample</summary>
  
  ```sql
  UPSERT INTO `question`
    ( `id`, `question`, `A`, `B`, `C`, `D`, `answer` )
  VALUES (1, "A knish is traditionally stuffed with what filling?", "potato", "creamed corn", "lemon custard", "raspberry jelly", "A"),
  ...;
  ```
</details>

### Local environment

Update or install neccesary software
- [Golang](https://go.dev/doc/install) >= 1.19

Set environment variables
```sh
$ export VAR=abc
```
| Variable | Value |
| ------ | ------ |
| BOT_TOKEN | [Obtain your bot's token](https://core.telegram.org/bots/tutorial#obtain-your-bot-token) |
| DEBUG | TRUE |
| YDB_CONNECTION_URL | [Get connection url](https://cloud.yandex.com/en/docs/ydb/operations/connection#endpoint-and-path) of your YDB |

## Build and Run

Open terminal, move to the backend's project directory, update all modules and run project
```
cd /path/to/repo/dir/yc-serverless-backend/
go get -u
go mod tidy
go run .
```

## Debug

Write your code to debug any method in main() function

```go
func main() {
  questions := repository.GetQuestions(0)
  log.Printf("Questions total: %d", len(questions))
}
```
Remember, that this two functions mast be called by a Cloud Functions service

```go
func QuestionHandler(rw http.ResponseWriter, req *http.Request) {
  handlers.Question(rw, req)
}

func UserHandler(rw http.ResponseWriter, req *http.Request) {
  handlers.User(rw, req)
}
```

## Deploy

The source code of this backend is intended to be deployed to a [Yandex Cloud Functions](https://cloud.yandex.com/en/services/functions). But it is possible to deploy it to another service such as [AWS Lambda](https://aws.amazon.com/lambda/) or [Google Cloud Functions](https://cloud.google.com/functions) (and any other serverless service which supports Go 1.19 runtime).

1. [Create](https://cloud.yandex.com/en/docs/functions/quickstart/create-function/go-function-quickstart) two functions for both /question and /user endpoints.
   Get an ID and a name of each function

2. Install [Yandex Cloud CLI](https://cloud.yandex.com/en/docs/cli/operations/install-cli) and login with your account

3. Set environment variables
```sh
$ export VAR=abc
```
| Variable | Value |
| ------ | ------ |
| QUESTION_FUNCTION_ID | Function ID for Question endpoint created at step 1 |
| QUESTION_FUNCTION_NAME | Function name for Question endpoint created at step 1 |
| USER_FUNCTION_ID | Function ID for User endpoint created at step 1 |
| USER_FUNCTION_NAME | Function name for User endpoint created at step 1 |
| SERVICE_ACCOUNT_ID | [Service account id](#create-your-database) for editing database |
| BOT_TOKEN | [Obtain your bot's token](#create-your-telegram-bot) |
| YDB_CONNECTION_URL | [Get connection url](#create-your-database) of your YDB |

```sh
cd /path/to/repo/dir/yc-serverless-backend/
sh ./scripts/deploy.sh
```

## Alternative options

1. Deploy you code manually with [Yandex Cloud Console](https://cloud.yandex.com/en/docs/functions/quickstart/create-function/go-function-quickstart) or other options

2. Generate any backend you like with Generate Client option of [Swagger Editor](https://editor.swagger.io/) using [openapi.yaml](https://github.com/AndreVasilev/OneMoreQuiz/blob/readme/openapi.yaml) and deploy it enywhere you like :)
