package redis

import (
	"encoding/json"
	redisWorker "github.com/topfreegames/go-workers"
)

func GetJobPayload(jobArg *redisWorker.Msg, payload interface{}) error {
	bts, err := jobArg.Args().MarshalJSON()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bts, payload); err != nil {
		return err
	}

	return nil
}
