package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Source struct {
	Start  int    `json:"start"`
	End    int    `json:"end"`
	Column Column `json:"column"`
}

type Column struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type Finding struct {
	CweIDs           []string `json:"cwe_ids"`
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	DocumentationURL string   `json:"documentation_url"`
	LineNumber       int      `json:"line_number"`
	Filename         string   `json:"filename"`
	CodeExtract      string   `json:"code_extract"`
	Sources          Source   `json:"source"`
}

type SecurityReport struct {
	High []Finding `json:"high"`
	Low  []Finding `json:"low"`
}

// Function to generate a Markdown table from the security report
func generateMarkdownTable(report SecurityReport, url string, namespace string, branch string) string {
	var buffer bytes.Buffer

	if len(report.High) == 0 && len(report.Low) == 0 {
		buffer.WriteString("# Yuki - SECURITY REPORT :eyes: \n")
		buffer.WriteString("![alt text](https://i.pinimg.com/564x/5f/a0/04/5fa004c77dcc99f43701294e27dd7a64.jpg)\n")
		return buffer.String()
	} else {
		buffer.WriteString("# :rotating_light: Yuki - SECURITY REPORT :rotating_light: \n")
		buffer.WriteString("![alt text](https://thoropass.com/wp-content/uploads/2023/10/828wnh.jpg)\n")

		buffer.WriteString("### :rotating_light: High Severity Findings\n\n")
		buffer.WriteString("| CWE | Title | Filename | Line Number | File | Documentation |\n")
		buffer.WriteString("|-----|-------|----------|-------------|---------------|---------------|\n")
		for _, finding := range report.High {
			// Concatenate URL, namespace, and the filename
			fullLink := fmt.Sprintf("%s/%s/-/blob/%s/%s#L%d-L%d", url, namespace, branch, finding.Filename, finding.Sources.Start, finding.Sources.End)

			// Use the fullLink in the markdown formatting
			buffer.WriteString(fmt.Sprintf("| %s | %s | %s | %d | [%s](%s) | [Documentation](%s) |\n",
				finding.CweIDs[0], finding.Title, finding.Filename, finding.LineNumber, finding.Filename, fullLink, finding.DocumentationURL))
		}

		// Low severity findings
		buffer.WriteString("\n### :warning: Low Severity Findings\n\n")
		buffer.WriteString("| CWE | Title | Filename | Line Number | File | Documentation |\n")
		buffer.WriteString("|-----|-------|----------|-------------|---------------|---------------|\n")
		for _, finding := range report.Low {
			// Concatenate URL, namespace, and the filename
			fullLink := fmt.Sprintf("%s/%s/-/blob/%s/%s#L%d-L%d", url, namespace, branch, finding.Filename, finding.Sources.Start, finding.Sources.End)

			// Use the fullLink in the markdown formatting
			buffer.WriteString(fmt.Sprintf("| %s | %s | %s | %d | [%s](%s) | [Documentation](%s) |\n",
				finding.CweIDs[0], finding.Title, finding.Filename, finding.LineNumber, finding.Filename, fullLink, finding.DocumentationURL))
		}

		return buffer.String()
	}
}

func sendGitLabNote(git, markdown string, projectID, mergeRequestID int, token string) error {
	// GitLab API URL
	url := fmt.Sprintf("%s/api/v4/projects/%d/merge_requests/%d/notes", git, projectID, mergeRequestID)

	// Create the request body, using %q to escape markdown content
	reqBody := []byte(fmt.Sprintf(`{"body": %q}`, markdown))

	// Ensure that the markdown content isn't too large for GitLab's API
	if len(reqBody) > 65536 {
		return fmt.Errorf("markdown content exceeds GitLab's maximum note size of 65536 characters")
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PRIVATE-TOKEN", token)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("received non-2xx response: %d", resp.StatusCode)
	}

	log.Println("Successfully sent note to GitLab")
	return nil
}

func main() {
	// Define Flags
	reportFile := flag.String("f", "", "Path to the file containing the security report")
	projectID := flag.Int("i", 0, "GitLab project ID $CI_PROJECT_ID")
	url := flag.String("u", "", "Project URL $GITLAB_URL")
	namespace := flag.String("n", "", "Project Namespace $$CI_PROJECT_PATH")
	mergeRequestID := flag.Int("m", 0, "Merge request ID $CI_MERGE_REQUEST_IID")
	gitlabToken := flag.String("t", "", "GitLab private token")
	flag.Parse()

	// Validate required flags
	if *reportFile == "" || *projectID == 0 || *mergeRequestID == 0 || *gitlabToken == "" || *url == "" || *namespace == "" {
		flag.Usage()
		os.Exit(1)
	}

	jsonData, err := os.ReadFile(*reportFile)
	if err != nil {
		log.Fatalf("Failed to read file")
	}

	// Check for empty JSON object
	if len(jsonData) == 2 && string(jsonData) == "{}" {
		log.Println("Success: The JSON file contains an empty object.")
		// Use a return statement to stop further execution
		return
	}

	// Parse the JSON data
	var report SecurityReport
	err = json.Unmarshal(jsonData, &report)
	if err != nil {
		log.Fatalf("Failed to parse JSON")
	}

	// Generate the Markdown table
	markdown := generateMarkdownTable(report, *url, *namespace, "development")

	// Send the Markdown table as a note to GitLab
	err = sendGitLabNote(*url, markdown, *projectID, *mergeRequestID, *gitlabToken)
	if err != nil {
		log.Fatalf("Failed to send note to GitLab")
	}
}
