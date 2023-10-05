package repository

// Learn more about YDB Go SDK
// https://ydb.tech/ru/docs/reference/ydb-sdk/example/go/
// https://github.com/ydb-platform/ydb-go-sdk

import (
	"context"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result/named"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"log"
	"os"
	"yc-serverless-backend/models"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	yc "github.com/ydb-platform/ydb-go-yc"
)

func authCredentials() ydb.Option {
	if os.Getenv("DEBUG") == "TRUE" {
		return yc.WithServiceAccountKeyFileCredentials("../authorized_key.json")
	} else {
		return yc.WithMetadataCredentials()
	}
}

func connect() (*ydb.Driver, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	db, err := ydb.Open(ctx,
		os.Getenv("YDB_CONNECTION_URL"),
		yc.WithInternalCA(),
		authCredentials(),
	)
	if err != nil {
		cancel()
		panic(err)
	}
	return db, ctx, cancel
}

func GetQuestions(from uint64) []models.Question {
	db, ctx, cancel := connect()
	defer cancel()
	defer func() {
		_ = db.Close(ctx)
	}()

	var (
		readTx = table.TxControl(
			table.BeginTx(
				table.WithOnlineReadOnly(),
			),
			table.CommitTx(),
		)
	)

	var questions []models.Question
	err := db.Table().Do(ctx,
		func(ctx context.Context, s table.Session) (err error) {
			var (
				res      result.Result
				id       uint64
				question *string
				a        *string
				b        *string
				c        *string
				d        *string
				answer   *string
			)
			_, res, err = s.Execute(
				ctx,
				readTx,
				`
        DECLARE $fromID AS Uint64;
        SELECT *
        FROM
          question
        WHERE
          id > $fromID;
      `,
				table.NewQueryParameters(
					table.ValueParam("$fromID", types.Uint64Value(from)), // подстановка в условие запроса
				),
			)
			if err != nil {
				log.Printf(err.Error())
				return nil
			}
			defer func(res result.Result) {
				err := res.Close()
				if err != nil {
					log.Printf(err.Error())
				}
			}(res) // закрытие result'а обязательно

			for res.NextResultSet(ctx) {
				for res.NextRow() {
					// в ScanNamed передаем имена колонок из строки сканирования,
					// адреса (и типы данных), куда следует присвоить результаты запроса
					err = res.ScanNamed(
						named.Required("id", &id),
						named.Optional("question", &question),
						named.Optional("A", &a),
						named.Optional("B", &b),
						named.Optional("C", &c),
						named.Optional("D", &d),
						named.Optional("answer", &answer),
					)
					if err != nil {
						log.Printf(err.Error())
						return nil
					}
					question := models.Question{
						Id:       id,
						Question: question,
						A:        a,
						B:        b,
						C:        c,
						D:        d,
						Answer:   answer,
					}
					questions = append(questions, question)
				}
			}
			return res.Err()
		},
	)
	if err != nil {
		// обработка ошибки выполнения запроса
		log.Printf(err.Error())
		return nil
	}
	return questions
}

func GetUser(userId int64) *models.User {
	db, ctx, cancel := connect()
	defer cancel()
	defer func() {
		_ = db.Close(ctx)
	}()

	var (
		readTx = table.TxControl(
			table.BeginTx(
				table.WithOnlineReadOnly(),
			),
			table.CommitTx(),
		)
	)

	var user *models.User
	err := db.Table().Do(ctx,
		func(ctx context.Context, s table.Session) (err error) {
			var (
				res            result.Result
				id             int64
				lastQuestionId uint64
				score          uint64
				successAnswers uint64
				tgData         *string
			)

			_, res, err = s.Execute(
				ctx,
				readTx,
				`
        DECLARE $id AS int64;
        SELECT *
        FROM
          user
        WHERE
          id = $id;
      `,
				table.NewQueryParameters(
					table.ValueParam("$id", types.Int64Value(userId)), // подстановка в условие запроса
				),
			)
			if err != nil {
				log.Printf(err.Error())
				return nil
			}
			defer func(res result.Result) {
				err := res.Close()
				if err != nil {
					log.Printf(err.Error())
				}
			}(res) // закрытие result'а обязательно

			for res.NextResultSet(ctx) {
				for res.NextRow() {
					// в ScanNamed передаем имена колонок из строки сканирования,
					// адреса (и типы данных), куда следует присвоить результаты запроса
					err = res.ScanNamed(
						named.Required("id", &id),
						named.Required("last_question_id", &lastQuestionId),
						named.Required("score", &score),
						named.Required("success_answers", &successAnswers),
						named.Optional("tg_data", &tgData),
					)
					if err != nil {
						log.Printf(err.Error())
						return nil
					}
					user = &models.User{
						Id:             id,
						LastQuestionId: lastQuestionId,
						Score:          score,
						SuccessAnswers: successAnswers,
						TgData:         tgData,
					}
					return nil
				}
			}
			return res.Err()
		},
	)
	if err != nil {
		// обработка ошибки выполнения запроса
		log.Printf(err.Error())
		return nil
	}
	return user
}

func AddUser(id int64, initData string) *models.User {
	db, ctx, cancel := connect()
	defer cancel()
	defer func() {
		_ = db.Close(ctx)
	}()

	writeTx := table.TxControl(
		table.BeginTx(
			table.WithSerializableReadWrite(),
		),
		table.CommitTx(),
	)

	err := db.Table().Do(ctx,
		func(ctx context.Context, s table.Session) (err error) {
			_, _, err = s.Execute(ctx, writeTx, `
        DECLARE $id AS int64;
        DECLARE $initData AS string;
        UPSERT INTO
          user
        (`+"`id`, `last_question_id`, `score`, `success_answers`, `tg_data`"+` )
		VALUES ( $id, 0, 0, 0, $initData);
      `, table.NewQueryParameters(
				table.ValueParam("$id", types.Int64Value(id)),
				table.ValueParam("$initData", types.StringValueFromString(initData)),
			))
			return err
		},
	)
	if err != nil {
		log.Printf(err.Error())
		return nil
	} else {
		return &models.User{
			Id:             id,
			LastQuestionId: 0,
			Score:          0,
			SuccessAnswers: 0,
			TgData:         &initData,
		}
	}
}

func UpdateUser(user *models.User) {
	db, ctx, cancel := connect()
	defer cancel()
	defer func() {
		_ = db.Close(ctx)
	}()

	writeTx := table.TxControl(
		table.BeginTx(
			table.WithSerializableReadWrite(),
		),
		table.CommitTx(),
	)

	err := db.Table().Do(ctx,
		func(ctx context.Context, s table.Session) (err error) {
			_, _, err = s.Execute(ctx, writeTx, `
        DECLARE $id AS int64;
        DECLARE $lastQuestionId AS Uint64;
        DECLARE $score AS Uint64;
        DECLARE $successAnswers AS Uint64;
        DECLARE $initData AS string;
        UPDATE user
		SET `+
				"`last_question_id` = $lastQuestionId, "+
				"`score` = $score, "+
				"`success_answers` = $successAnswers, "+
				"`tg_data` = $initData "+`
		WHERE id = $id;`,
				table.NewQueryParameters(
					table.ValueParam("$id", types.Int64Value(user.Id)),
					table.ValueParam("$lastQuestionId", types.Uint64Value(user.LastQuestionId)),
					table.ValueParam("$score", types.Uint64Value(user.Score)),
					table.ValueParam("$successAnswers", types.Uint64Value(user.SuccessAnswers)),
					table.ValueParam("$initData", types.StringValueFromString(*user.TgData)),
				))
			return err
		},
	)
	if err != nil {
		log.Printf(err.Error())
	}
}

func GetUsersUpperPositionCount(score uint64) int {
	db, ctx, cancel := connect()
	defer cancel()
	defer func() {
		_ = db.Close(ctx)
	}()

	var (
		readTx = table.TxControl(
			table.BeginTx(
				table.WithOnlineReadOnly(),
			),
			table.CommitTx(),
		)
	)

	count := 0
	err := db.Table().Do(ctx,
		func(ctx context.Context, s table.Session) (err error) {
			var res result.Result

			_, res, err = s.Execute(
				ctx,
				readTx,
				`
        DECLARE $score AS uint64;
        SELECT *
        FROM
          user
        WHERE
          score > $score;
      `,
				table.NewQueryParameters(
					table.ValueParam("$score", types.Uint64Value(score)), // подстановка в условие запроса
				),
			)
			if err != nil {
				log.Printf(err.Error())
				return nil
			}
			defer func(res result.Result) {
				err := res.Close()
				if err != nil {
					log.Printf(err.Error())
				}
			}(res) // закрытие result'а обязательно

			for res.NextResultSet(ctx) {
				for res.NextRow() {
					count += 1
				}
			}
			return res.Err()
		},
	)
	if err != nil {
		// обработка ошибки выполнения запроса
		log.Printf(err.Error())
		return 0
	}
	return count
}
