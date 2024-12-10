package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/itsorganic/panto/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/gitlab"
)

var (
	GithubOAuthConfig *oauth2.Config
	GitlabOAuthConfig *oauth2.Config
)
var (
	GithubRepos []models.GithubRepo
	GitlabRepos []models.GitlabRepo
)

var FrontendURL = "https://panto-frontend-production.up.railway.app"

func init() {
	clientId := os.Getenv("GITLAB_CLIENT_ID")
	clientSecret := os.Getenv("GITLAB_CLIENT_SECRET")
	GitlabOAuthConfig = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  "https://panto-backend-production.up.railway.app/gitlab/auth/callback",
		Scopes:       []string{"read_user", "read_api"},
		Endpoint:     gitlab.Endpoint,
	}

	gclientId := os.Getenv("GITHUB_CLIENT_ID")
	gclientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	GithubOAuthConfig = &oauth2.Config{
		ClientID:     gclientId,
		ClientSecret: gclientSecret,
		RedirectURL:  "https://panto-backend-production.up.railway.app/github/auth/callback",
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
	}
}

func GitlabAuthHandler(c *gin.Context) {
	url := GitlabOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

func GithubAuthHandler(c *gin.Context) {
	url := GithubOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

func GithubAuthCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}
	// Exchange the authorization code for an access token
	token, err := GithubOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange tokens: " + err.Error()})
		return
	}
	c.SetCookie("gh-accessToken", token.AccessToken, 3600, "/", "panto-backend-production.up.railway.app", true, true)
	c.SetCookie("provider", "github", 3600, "/", "panto-backend-production.up.railway.app", true, true)
	c.Redirect(http.StatusFound, fmt.Sprintf("%s/dashboard", FrontendURL))
}

func GitlabAuthCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	// Exchange the authorization code for an access token
	token, err := GitlabOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange tokens %s" + err.Error()})
		return
	}
	var user models.GitlabUser

	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	resp, err := client.Get("https://gitlab.com/api/v4/user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user information %s" + err.Error()})
		return
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode the user information %s" + err.Error()})
		return
	}
	c.SetCookie("accessToken", token.AccessToken, 3600, "/", "panto-backend-production.up.railway.app", true, true)
	c.SetCookie("provider", "gitlab", 3600, "/", "panto-backend-production.up.railway.app", true, true)
	c.Redirect(http.StatusFound, fmt.Sprintf("%s/dashboard", FrontendURL))
}

func GetGithubUserDetails(c *gin.Context) {
	token, err := c.Cookie("gh-accessToken")
	if err != nil {
		c.Redirect(http.StatusFound, "https://panto-backend-production.up.railway.app/github")
		return
	}

	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user information: " + err.Error()})
		return
	}
	defer resp.Body.Close()
	var user models.GithubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode the user information: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func FetchGithubUserRepo(c *gin.Context) {
	// Get the GitHub token from the cookie
	token, err := c.Cookie("gh-accessToken")
	if err != nil {
		c.Redirect(http.StatusFound, "https://panto-backend-production.up.railway.app/github")
		return
	}

	// Use the standalone function to fetch repositories
	repos, err := FetchGithubRepos(token, 5) // Limit to 5 repositories
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the repositories
	c.JSON(http.StatusOK, gin.H{"user": repos})
}

func GetGitlabUserDetails(c *gin.Context) {
	token, err := c.Cookie("accessToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get the access token: " + err.Error()})
		return
	}

	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	resp, err := client.Get("https://gitlab.com/api/v4/user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user information: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	var user models.GitlabUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode the user information: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func FetchGitlabUserRepo(c *gin.Context) {
	// Get the token from the cookie
	token, err := c.Cookie("accessToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get the access token: " + err.Error()})
		return
	}

	// Use the standalone function to fetch repositories
	repos, err := FetchGitlabRepos(token, 5) // Limit to 5 repositories
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the repositories
	c.JSON(http.StatusOK, gin.H{"repos": repos})
}

func GitlabLogout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "accessToken",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Logout successfull"})
}

func GithubLogout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "gh-accessToken",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Logout successfull"})
}

func FetchGitlabRepos(token string, limit int) ([]models.GitlabRepo, error) {
	// Create an HTTP client with the access token
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))

	// Make the API request to fetch user's projects
	resp, err := client.Get("https://gitlab.com/api/v4/projects?membership=true&per_page=100")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %v", err)
	}
	defer resp.Body.Close()

	// Decode the repositories
	var repos []models.GitlabRepo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("failed to decode repositories: %v", err)
	}

	// Limit the number of repositories
	if len(repos) > limit {
		repos = repos[:limit]
	}

	// Attach review status to each repository
	for i := range repos {
		// Use a unique identifier for the repository
		key := fmt.Sprintf("gitlab:%d", repos[i].ID)
		repos[i].Review = GitlabRepoReviews[key]
	}

	return repos, nil
}

// Modify to persist review status globally and across sessions
var (
	GithubRepoReviews = make(map[string]bool) // Add a global map to track review status
)

func FetchGithubRepos(token string, limit int) ([]models.GithubRepo, error) {
	// Create an HTTP client with the access token
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))

	// Make the API request
	resp, err := client.Get("https://api.github.com/user/repos")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %v", err)
	}
	defer resp.Body.Close()

	// Reset GithubRepos before populating
	GithubRepos = []models.GithubRepo{}

	// Parse the response
	if err := json.NewDecoder(resp.Body).Decode(&GithubRepos); err != nil {
		return nil, fmt.Errorf("failed to decode repositories: %v", err)
	}

	// Limit the repositories
	if len(GithubRepos) > limit {
		GithubRepos = GithubRepos[:limit]
	}

	// Set review status from the persistent map
	for i := range GithubRepos {
		fullName := GithubRepos[i].FullName
		GithubRepos[i].Review = GithubRepoReviews[fullName]
	}

	return GithubRepos, nil
}

var GitlabRepoReviews = make(map[string]bool)

func HandleGithubToggleReview(c *gin.Context) {
	var requestBody struct {
		RepoFullName string `json:"repoFullName"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Toggle review status in the global map
	currentStatus := GithubRepoReviews[requestBody.RepoFullName]
	GithubRepoReviews[requestBody.RepoFullName] = !currentStatus

	// Update the status in GithubRepos if the repo is currently loaded
	for i := range GithubRepos {
		if GithubRepos[i].FullName == requestBody.RepoFullName {
			GithubRepos[i].Review = !currentStatus
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Repository review status toggled",
		"isReviewed": !currentStatus,
	})
}

func HandleGitlabToggleReview(c *gin.Context) {
	var requestBody struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Create a unique key for the repository
	key := fmt.Sprintf("gitlab:%d", requestBody.ID)

	// Toggle review status
	currentStatus := GitlabRepoReviews[key]
	GitlabRepoReviews[key] = !currentStatus

	// Update the status in the current GitlabRepos slice if loaded
	for i := range GitlabRepos {
		if GitlabRepos[i].ID == requestBody.ID {
			GitlabRepos[i].Review = !currentStatus
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Repository review status toggled",
		"isReviewed": !currentStatus,
		"repoId":     requestBody.ID,
	})
}

// HandleAutoReview toggles the "Auto Review" status for a given repository
