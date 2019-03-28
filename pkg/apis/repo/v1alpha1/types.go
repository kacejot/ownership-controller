type Repo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RepoSpec   `json:"spec"`
	Status            RepoStatus `json:"status"`
}

type RepoSpec struct {
	RepoAddress   string   `json:"repoAddress"`
	DockerImage   string   `json:"dockerImage"`
	BuildCommands []string `json:"buildCommands"`
	Artifacts     []string `json:"artifacts"`
}

type RepoStatus struct {
	Log []string `json:"log"`
}
