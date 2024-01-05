package dbconv

import (
	"time"

	"github.com/senyc/jason/pkg/types"
)

func ToTaskResponse(r types.SqlTasksRow) (types.TaskReponse, error) {
	result := types.TaskReponse{
		Id:            r.Id,
		Title:         r.Title,
		Body:          "",
		Due:           types.NilTime{},
		Priority:      r.Priority,
		Completed:     r.Completed,
		CompletedDate: types.NilTime{},
	}

	if r.Due.Valid {
		result.Due = types.NilTime{r.Due.Time}
	}

	if r.Body.Valid {
		result.Body = r.Body.String
	}

	if r.CompletedDate.Valid {
		result.CompletedDate = types.NilTime{r.CompletedDate.Time}
	}

	return result, nil
}

func ToCompletedTaskResponse(r types.SqlTasksRow) (types.CompletedTaskResponse, error) {
	result := types.CompletedTaskResponse{
		Id:            r.Id,
		Title:         r.Title,
		Body:          "",
		Due:           types.NilTime{},
		Priority:      r.Priority,
		CompletedDate: time.Time{},
	}

	if r.Due.Valid {
		result.Due = types.NilTime{r.Due.Time}
	}

	if r.CompletedDate.Valid {
		result.CompletedDate = r.CompletedDate.Time
	}

	if r.Body.Valid {
		result.Body = r.Body.String
	}

	return result, nil
}

func ToIncompleteTaskResponse(r types.SqlTasksRow) (types.IncompleteTaskResponse, error) {
	result := types.IncompleteTaskResponse{
		Id:       r.Id,
		Title:    r.Title,
		Body:     "",
		Due:      types.NilTime{},
		Priority: r.Priority,
	}

	if r.Due.Valid {
		result.Due = types.NilTime{r.Due.Time}
	}

	if r.Body.Valid {
		result.Body = r.Body.String
	}

	return result, nil
}
