package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reporter/internal/service"
)

// pingHandler returns a pong response
// @Summary Pong
// @Description get a pong response
// @Tags example
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func callbackHandler(c *gin.Context, oidcService *service.OIDCService) {
	// Extract the authorization code from the query parameters
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	// Exchange the authorization code for tokens
	oauth2Token, err := oidcService.Config.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Extract the ID token from oauth2Token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No id_token field in oauth2 token"})
		return
	}

	// Optionally, verify the ID token's integrity and authenticity
	idToken, err := oidcService.Verifier.Verify(c.Request.Context(), rawIDToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ID Token"})
		return
	}

	// Here, you can extract claims from idToken, create sessions, etc.
	// For example:
	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract claims from ID Token"})
		return
	}

	// Redirect the user or return a success response
	// For simplicity, let's redirect the user to the /ping route
	c.Redirect(http.StatusFound, "/ping")
}
