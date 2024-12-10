package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsorganic/panto/auth"
	"github.com/itsorganic/panto/middlewares"
)

func main() {
	route := gin.Default()

	route.Use(middlewares.CORSMiddleware())
	route.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Server is working fine")
	})
	route.GET("/github", auth.GithubAuthHandler)
	route.GET("github/auth/callback", auth.GithubAuthCallback)

	route.GET("/gitlab", auth.GitlabAuthHandler)
	route.GET("gitlab/auth/callback", auth.GitlabAuthCallback)

	route.GET("/github/dashboard", auth.GetGithubUserDetails)
	route.GET("/github/dashboard/repo", auth.FetchGithubUserRepo)

	route.GET("/gitlab/dashboard", auth.GetGitlabUserDetails)
	route.GET("/gitlab/dashboard/repo", auth.FetchGitlabUserRepo)

	route.GET("/github/logout", auth.GithubLogout)
	route.GET("/gitlab/logout", auth.GitlabLogout)
	route.POST("/github/review", auth.HandleGithubToggleReview)
	route.POST("/gitlab/review", auth.HandleGitlabToggleReview)

	route.Run(":8080")
}
