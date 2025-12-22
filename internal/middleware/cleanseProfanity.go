package middleware

import (
	"github.com/GhostVox/ghostvox.io-backend/internal/utils"
	"net/http"
)

func Cleanse( next http.HandlerFunc, filter *utils.Trie) http.HandlerFunc {
return http.HandlerFunc(func (w http.ResponseWriter, r http.Request) http.HandlerFunc {
		Body := utils.GetRequestBody(r)
})
