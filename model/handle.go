package model

import "time"

type Gogs struct {
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	After      string `json:"after"`
	CompareURL string `json:"compare_url"`
	Commits    []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
		URL     string `json:"url"`
		Author  struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"author"`
		Committer struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"committer"`
		Added     []interface{} `json:"added"`
		Removed   []interface{} `json:"removed"`
		Modified  []string      `json:"modified"`
		Timestamp time.Time     `json:"timestamp"`
	} `json:"commits"`
	Repository struct {
		ID    int `json:"id"`
		Owner struct {
			ID        int    `json:"id"`
			Username  string `json:"username"`
			Login     string `json:"login"`
			FullName  string `json:"full_name"`
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
		} `json:"owner"`
		Name            string      `json:"name"`
		FullName        string      `json:"full_name"`
		Description     string      `json:"description"`
		Private         bool        `json:"private"`
		Fork            bool        `json:"fork"`
		Parent          interface{} `json:"parent"`
		Empty           bool        `json:"empty"`
		Mirror          bool        `json:"mirror"`
		Size            int         `json:"size"`
		HTMLURL         string      `json:"html_url"`
		SSHURL          string      `json:"ssh_url"`
		CloneURL        string      `json:"clone_url"`
		Website         string      `json:"website"`
		StarsCount      int         `json:"stars_count"`
		ForksCount      int         `json:"forks_count"`
		WatchersCount   int         `json:"watchers_count"`
		OpenIssuesCount int         `json:"open_issues_count"`
		DefaultBranch   string      `json:"default_branch"`
		CreatedAt       time.Time   `json:"created_at"`
		UpdatedAt       time.Time   `json:"updated_at"`
	} `json:"repository"`
	Pusher struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Login     string `json:"login"`
		FullName  string `json:"full_name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	} `json:"pusher"`
	Sender struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Login     string `json:"login"`
		FullName  string `json:"full_name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	} `json:"sender"`
}

type Config []struct {
	Platform  string `json:"platform"`
	Namespace string `json:"namespace"`
	Path      string `json:"path"`
	Branch    string `json:"branch"`
	Proto     string `json:"proto"`
	Password  string `json:"password"`
}

type ProjectConfig struct {
	Path     string `json:"path"`
	Head     string `json:"head"`
	Password string `json:"password"`
}
