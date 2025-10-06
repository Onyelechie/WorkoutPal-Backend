package model

type Exercise struct {
	ID                  int64    `json:"id"`
	Name                string   `json:"name"`
	Targets             []string `json:"targets"`
	Intensity           string   `json:"intensity"`
	Expertise           string   `json:"expertise"`
	Image               string   `json:"image"`
	Demo                string   `json:"demo"`
	RecommendedCount    int      `json:"recommendedCount"`
	RecommendedSets     int      `json:"recommendedSets"`
	RecommendedDuration int      `json:"recommendedDuration"`
	Custom              bool     `json:"custom"`
}

type ReadExerciseRequest struct {
	Target    string `json:"target"`
	Intensity string `json:"intensity"`
	Expertise string `json:"expertise"`
}

type CreateExerciseRequest struct {
	Name                string   `json:"name"`
	Targets             []string `json:"targets"`
	Intensity           string   `json:"intensity"`
	Expertise           string   `json:"expertise"`
	Image               string   `json:"image"`
	Demo                string   `json:"demo"`
	RecommendedCount    int      `json:"recommendedCount"`
	RecommendedSets     int      `json:"recommendedSets"`
	RecommendedDuration int      `json:"recommendedDuration"`
}
