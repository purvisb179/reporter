// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package api

import (
	"context"
	"github.com/gin-contrib/sessions"
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
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Router /ping [get]
func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Pong!",
	})
}

// downloadReportHandler streams an Excel report to the client, prompting a file download.
// This handler generates a report based on transactions and their labels, then returns the report
// as an Excel file. The Content-Disposition header suggests to the client that the response
// should be downloaded and saved as a file.
//
// @Summary Download Report
// @Description Downloads an Excel report containing transactions and their associated labels.
// @Tags reports
// @Accept  json
// @Produce  application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security BearerAuth
// @Param labels query []string false "Array of labels to filter the report by"
// @Success 200 {file} file "Excel report"
// @Router /reports/download [get]
func downloadReportHandler(c *gin.Context, reportService *service.ReportService) {
	labels := c.QueryArray("labels")

	// Generate the report with labels
	report, err := reportService.GenerateReport(labels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
		return
	}

	// Set headers to instruct the browser to download the file
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=report.xlsx")

	// Write the Excel file to the response
	err = report.Write(c.Writer)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
}

func callbackHandler(c *gin.Context, oidcService *service.OIDCService) {
	// Extract the authorization code and state from the query parameters
	code := c.Query("code")
	//state := c.Query("state")
	//todo Verify the state matches the one stored in the session or passed initially for CSRF protection

	if code == "" {
		// No code was found; handle the error appropriately
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is required"})
		return
	}

	// Exchange the authorization code for tokens
	oauth2Token, err := oidcService.Config.Exchange(context.Background(), code)
	if err != nil {
		// Handle error: Token exchange failed
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange authorization code for tokens"})
		return
	}

	// Extract the ID token from OAuth2 token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		// Handle error: ID Token not found in the OAuth2 token response
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID token not found in the token response"})
		return
	}

	// Parse and verify the ID token payload
	idToken, err := oidcService.Verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		// Handle error: ID Token verification failed
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ID token"})
		return
	}

	// Decode the token claims to get user information
	var claims struct {
		Email string `json:"email"`
		// Add other claims you need
	}
	if err := idToken.Claims(&claims); err != nil {
		// Handle error: Failed to parse token claims
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token claims"})
		return
	}

	// At this point, authentication was successful, so create a session for the user
	session := sessions.Default(c)
	session.Set("user", claims.Email) // Store the email in the session, or any other information needed
	if err := session.Save(); err != nil {
		// Handle error: Failed to save the session
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a session for the user"})
		return
	}

	// Redirect the user to the home page or another target page
	c.Redirect(http.StatusSeeOther, "/")
}
