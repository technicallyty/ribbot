package MediaBot

// Manager describes a type that can validate URLs and download MediaBot based on those urls
type Manager interface {
	Download() (string, string, error) // downloads a file and returns the save path, name of the downloaded item, + errors, if any
	ResourceURL() string               // returns the url used to fetch the resource
	IsValidURL() bool                  // validates the resource URL
	GetMediaDir() string               // gets
}
