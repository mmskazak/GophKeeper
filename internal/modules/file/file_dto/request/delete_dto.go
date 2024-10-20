package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DeleteFileDTO struct {
	FileID string `json:"file_id"`
	UserID int    `json:"user_id"`
}

func DeleteFileDTOFromHTTP(r *http.Request) (DeleteFileDTO, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return DeleteFileDTO{}, fmt.Errorf("reading body registration: %w", err)
	}
	var deletePwdDTO DeleteFileDTO
	err = json.Unmarshal(data, &deletePwdDTO)
	if err != nil {
		return DeleteFileDTO{}, fmt.Errorf("unmarshalling body registration: %w", err)
	}

	// Извлекаем userID из контекста
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		return DeleteFileDTO{}, fmt.Errorf("error getUserIDFromContext: %w", err)
	}

	deletePwdDTO.UserID = userID
	return deletePwdDTO, nil
}
