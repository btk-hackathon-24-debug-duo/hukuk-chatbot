package models

type SafetyRating struct {
	Category    int  `json:"Category"`
	Probability int  `json:"Probability"`
	Blocked     bool `json:"Blocked"`
}

type Content struct {
	Parts []string `json:"Parts"`
	Role  string   `json:"Role"`
}

type Candidate struct {
	Index            int            `json:"Index"`
	Content          Content        `json:"Content"`
	FinishReason     int            `json:"FinishReason"`
	SafetyRatings    []SafetyRating `json:"SafetyRatings"`
	CitationMetadata interface{}    `json:"CitationMetadata"`
	TokenCount       int            `json:"TokenCount"`
}

type UsageMetadata struct {
	PromptTokenCount        int `json:"PromptTokenCount"`
	CachedContentTokenCount int `json:"CachedContentTokenCount"`
	CandidatesTokenCount    int `json:"CandidatesTokenCount"`
	TotalTokenCount         int `json:"TotalTokenCount"`
}

type DataModel struct {
	Candidates     []Candidate   `json:"Candidates"`
	PromptFeedback interface{}   `json:"PromptFeedback"`
	UsageMetadata  UsageMetadata `json:"UsageMetadata"`
}
