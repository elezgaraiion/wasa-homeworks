package api

import (
	"encoding/json"
	"net/http"
	
	"crypto/rand"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	
)

func saveImageFile(file multipart.File, handle *multipart.FileHeader) (string, error) {
	if err := os.MkdirAll("uploads", 0755); err != nil {
		return "", err
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	fileName := hex.EncodeToString(randomBytes) + filepath.Ext(handle.Filename)

	dstPath := filepath.Join("uploads", fileName)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return "/uploads/" + fileName, nil
}

func extractBearer(authorization string) string {
	var tokens = strings.Split(authorization, " ")
	if len(tokens) == 2 {
		return strings.Trim(tokens[1], " ")
	}
	return ""
}
func (rt *_router) jsonError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}