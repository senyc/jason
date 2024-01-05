package dbconv

import (
	"time"

	"github.com/senyc/jason/pkg/types"
)

func ToTaskResponse(r types.SqlTasksRow) (types.TaskReponse, error) {
	var completedDate *time.Time
	result := types.TaskReponse{
		Id:            r.Id,
		Title:         r.Title,
		Body:          "",
		Due:           types.NullTime{},
		Priority:      r.Priority,
		Completed:     r.Completed,
		CompletedDate: completedDate,
	}

	if r.Due.Valid {
		result.Due = types.NullTime{r.Due.Time}
	}

	if r.Body.Valid {
		result.Body = r.Body.String
	}

	if r.CompletedDate.Valid {
		result.CompletedDate = &r.CompletedDate.Time
	}

	return result, nil
}

func ToCompletedTaskResponse(r types.SqlTasksRow) (types.CompletedTaskResponse, error) {
	result := types.CompletedTaskResponse{
		Id:            r.Id,
		Title:         r.Title,
		Body:          "",
		Due:           types.NullTime{},
		Priority:      r.Priority,
		CompletedDate: time.Time{},
	}

	if r.Due.Valid {
		result.Due = types.NullTime{r.Due.Time}
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
		Due:      types.NullTime{},
		Priority: r.Priority,
	}

	if r.Due.Valid {
		result.Due = types.NullTime{r.Due.Time}
	}

	if r.Body.Valid {
		result.Body = r.Body.String
	}

	return result, nil
}
