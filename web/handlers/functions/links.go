package functions

import (
	"fmt"
	"strings"
	"text/template"

	visualizer "github.com/jenkins-x/jx-pipelines-visualizer"

	jenkinsv1 "github.com/jenkins-x/jx-api/v4/pkg/apis/jenkins.io/v1"
)

func TraceURLFunc(pipelineTraceURLTemplate *template.Template) func(string) string {
	return func(traceID string) string {
		return traceIDToTraceURL(traceID, pipelineTraceURLTemplate)
	}
}

func RepositoryURL(pipeline interface{}) string {
	switch p := pipeline.(type) {
	case visualizer.Pipeline:
		return repositoryURLForPipeline(p)
	case visualizer.RunningPipeline:
		return repositoryURLForPipeline(p.Pipeline)
	case *jenkinsv1.PipelineActivity:
		return repositoryURLForPipelineActivity(p)
	default:
		return ""
	}
}

func BranchURL(pipeline interface{}) string {
	switch p := pipeline.(type) {
	case visualizer.Pipeline:
		return branchURLForPipeline(p)
	case visualizer.RunningPipeline:
		return branchURLForPipeline(p.Pipeline)
	case *jenkinsv1.PipelineActivity:
		return branchURLForPipelineActivity(p)
	default:
		return ""
	}
}

func CommitURL(pipeline interface{}) string {
	switch p := pipeline.(type) {
	case *jenkinsv1.PipelineActivity:
		return commitURLForPipelineActivity(p)
	default:
		return ""
	}
}

func AuthorURL(pipeline interface{}) string {
	switch p := pipeline.(type) {
	case visualizer.Pipeline:
		return authorURLForPipeline(p)
	case visualizer.RunningPipeline:
		return authorURLForPipeline(p.Pipeline)
	case *jenkinsv1.PipelineActivity:
		return authorURLForPipelineActivity(p)
	default:
		return ""
	}
}

func authorURLForPipeline(pipeline visualizer.Pipeline) string {
	switch pipeline.Provider {
	case "github":
		return fmt.Sprintf("https://github.com/%s", pipeline.Author)
	default:
		return ""
	}
}

func authorURLForPipelineActivity(pa *jenkinsv1.PipelineActivity) string {
	switch pipelineActivityProvider(pa) {
	case "github":
		return fmt.Sprintf("https://github.com/%s", pa.Spec.Author)
	default:
		return ""
	}
}

func repositoryURLForPipeline(pipeline visualizer.Pipeline) string {
	switch pipeline.Provider {
	case "github":
		return fmt.Sprintf("https://github.com/%s/%s", pipeline.Owner, pipeline.Repository)
	default:
		return ""
	}
}

func repositoryURLForPipelineActivity(pa *jenkinsv1.PipelineActivity) string {
	switch pipelineActivityProvider(pa) {
	case "github":
		return fmt.Sprintf("https://github.com/%s/%s", pa.Spec.GitOwner, pa.Spec.GitRepository)
	default:
		return ""
	}
}

func branchURLForPipeline(pipeline visualizer.Pipeline) string {
	switch pipeline.Provider {
	case "github":
		if pipeline.PullRequestNumber() != "" {
			return fmt.Sprintf("https://github.com/%s/%s/pull/%s", pipeline.Owner, pipeline.Repository, pipeline.PullRequestNumber())
		}
		return fmt.Sprintf("https://github.com/%s/%s/tree/%s", pipeline.Owner, pipeline.Repository, pipeline.Branch)
	default:
		return ""
	}
}

func branchURLForPipelineActivity(pa *jenkinsv1.PipelineActivity) string {
	switch pipelineActivityProvider(pa) {
	case "github":
		if strings.HasPrefix(pa.Spec.GitBranch, "PR-") {
			return fmt.Sprintf("https://github.com/%s/%s/pull/%s", pa.Spec.GitOwner, pa.Spec.GitRepository, strings.TrimPrefix(pa.Spec.GitBranch, "PR-"))
		}
		return fmt.Sprintf("https://github.com/%s/%s/tree/%s", pa.Spec.GitOwner, pa.Spec.GitRepository, pa.Spec.GitBranch)
	default:
		return ""
	}
}

func commitURLForPipelineActivity(pa *jenkinsv1.PipelineActivity) string {
	if len(pa.Spec.LastCommitURL) > 0 {
		return pa.Spec.LastCommitURL
	}
	switch pipelineActivityProvider(pa) {
	case "github":
		return fmt.Sprintf("https://github.com/%s/%s/commit/%s", pa.Spec.GitOwner, pa.Spec.GitRepository, pa.Spec.LastCommitSHA)
	default:
		return ""
	}
}

func pipelineActivityProvider(pa *jenkinsv1.PipelineActivity) string {
	if provider := pa.Labels["provider"]; provider != "" {
		return provider
	}

	if strings.Contains(pa.Spec.GitURL, "github") {
		return "github"
	}

	return ""
}

func traceIDToTraceURL(traceID string, pipelineTraceURLTemplate *template.Template) string {
	if pipelineTraceURLTemplate == nil {
		return ""
	}
	if traceID == "" {
		return ""
	}

	sb := new(strings.Builder)
	err := pipelineTraceURLTemplate.Execute(sb, map[string]string{
		"TraceID": traceID,
	})
	if err != nil {
		return err.Error()
	}
	return sb.String()
}
