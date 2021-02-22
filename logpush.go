package cloudflare

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// LogpushJob describes a Logpush job.
type LogpushJob struct {
	ID                 int        `json:"id,omitempty"`
	Dataset            string     `json:"dataset"`
	Enabled            bool       `json:"enabled"`
	Name               string     `json:"name"`
	LogpullOptions     string     `json:"logpull_options"`
	DestinationConf    string     `json:"destination_conf"`
	OwnershipChallenge string     `json:"ownership_challenge,omitempty"`
	LastComplete       *time.Time `json:"last_complete,omitempty"`
	LastError          *time.Time `json:"last_error,omitempty"`
	ErrorMessage       string     `json:"error_message,omitempty"`
}

// LogpushJobsResponse is the API response, containing an array of Logpush Jobs.
type LogpushJobsResponse struct {
	Response
	Result []LogpushJob `json:"result"`
}

// LogpushJobDetailsResponse is the API response, containing a single Logpush Job.
type LogpushJobDetailsResponse struct {
	Response
	Result LogpushJob `json:"result"`
}

// LogpushFieldsResponse is the API response for a datasets fields
type LogpushFieldsResponse struct {
	Response
	Result LogpushFields `json:"result"`
}

// LogpushFields is a map of available Logpush field names & descriptions
type LogpushFields map[string]string

// LogpushGetOwnershipChallenge describes a ownership validation.
type LogpushGetOwnershipChallenge struct {
	Filename string `json:"filename"`
	Valid    bool   `json:"valid"`
	Message  string `json:"message"`
}

// LogpushGetOwnershipChallengeResponse is the API response, containing a ownership challenge.
type LogpushGetOwnershipChallengeResponse struct {
	Response
	Result LogpushGetOwnershipChallenge `json:"result"`
}

// LogpushGetOwnershipChallengeRequest is the API request for get ownership challenge.
type LogpushGetOwnershipChallengeRequest struct {
	DestinationConf string `json:"destination_conf"`
}

// LogpushOwnershipChallengeValidationResponse is the API response,
// containing a ownership challenge validation result.
type LogpushOwnershipChallengeValidationResponse struct {
	Response
	Result struct {
		Valid bool `json:"valid"`
	}
}

// LogpushValidateOwnershipChallengeRequest is the API request for validate ownership challenge.
type LogpushValidateOwnershipChallengeRequest struct {
	DestinationConf    string `json:"destination_conf"`
	OwnershipChallenge string `json:"ownership_challenge"`
}

// LogpushDestinationExistsResponse is the API response,
// containing a destination exists check result.
type LogpushDestinationExistsResponse struct {
	Response
	Result struct {
		Exists bool `json:"exists"`
	}
}

// LogpushDestinationExistsRequest is the API request for check destination exists.
type LogpushDestinationExistsRequest struct {
	DestinationConf string `json:"destination_conf"`
}

// CreateLogpushJob creates a new LogpushJob for a zone.
//
// API reference: https://api.cloudflare.com/#logpush-jobs-create-logpush-job
func (api *API) CreateLogpushJob(ctx context.Context, zoneID string, job LogpushJob) (*LogpushJob, error) {
	uri := "/zones/" + zoneID + "/logpush/jobs"
	res, err := api.makeRequestContext(ctx, "POST", uri, job)
	if err != nil {
		return nil, errors.Wrap(err, errMakeRequestError)
	}
	var r LogpushJobDetailsResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, errors.Wrap(err, errUnmarshalError)
	}
	return &r.Result, nil
}

// LogpushJobs returns all Logpush Jobs for a zone.
//
// API reference: https://api.cloudflare.com/#logpush-jobs-list-logpush-jobs
func (api *API) LogpushJobs(ctx context.Context, zoneID string) ([]LogpushJob, error) {
	uri := "/zones/" + zoneID + "/logpush/jobs"
	res, err := api.makeRequestContext(ctx, "GET", uri, nil)
	if err != nil {
		return []LogpushJob{}, errors.Wrap(err, errMakeRequestError)
	}
	var r LogpushJobsResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return []LogpushJob{}, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// LogpushJobsForDataset returns all Logpush Jobs for a dataset in a zone.
//
// API reference: https://api.cloudflare.com/#logpush-jobs-list-logpush-jobs-for-a-dataset
func (api *API) LogpushJobsForDataset(ctx context.Context, zoneID, dataset string) ([]LogpushJob, error) {
	uri := "/zones/" + zoneID + "/logpush/datasets/" + dataset + "/jobs"
	res, err := api.makeRequestContext(ctx, "GET", uri, nil)
	if err != nil {
		return []LogpushJob{}, errors.Wrap(err, errMakeRequestError)
	}
	var r LogpushJobsResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return []LogpushJob{}, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// LogpushFields returns fields for a given dataset.
//
// API reference: https://api.cloudflare.com/#logpush-jobs-list-logpush-jobs
func (api *API) LogpushFields(ctx context.Context, zoneID, dataset string) (LogpushFields, error) {
	uri := "/zones/" + zoneID + "/logpush/datasets/" + dataset + "/fields"
	res, err := api.makeRequestContext(ctx, "GET", uri, nil)
	if err != nil {
		return LogpushFields{}, errors.Wrap(err, errMakeRequestError)
	}
	var r LogpushFieldsResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return LogpushFields{}, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// LogpushJob fetches detail about one Logpush Job for a zone.
//
// API reference: https://api.cloudflare.com/#logpush-jobs-logpush-job-details
func (api *API) LogpushJob(ctx context.Context, zoneID string, jobID int) (LogpushJob, error) {
	uri := "/zones/" + zoneID + "/logpush/jobs/" + strconv.Itoa(jobID)
	res, err := api.makeRequestContext(ctx, "GET", uri, nil)
	if err != nil {
		return LogpushJob{}, errors.Wrap(err, errMakeRequestError)
	}
	var r LogpushJobDetailsResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return LogpushJob{}, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// UpdateLogpushJob lets you update a Logpush Job.
//
// API reference: https://api.cloudflare.com/#logpush-jobs-update-logpush-job
func (api *API) UpdateLogpushJob(ctx context.Context, zoneID string, jobID int, job LogpushJob) error {
	uri := "/zones/" + zoneID + "/logpush/jobs/" + strconv.Itoa(jobID)
	res, err := api.makeRequestContext(ctx, "PUT", uri, job)
	if err != nil {
		return errors.Wrap(err, errMakeRequestError)
	}
	var r LogpushJobDetailsResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return errors.Wrap(err, errUnmarshalError)
	}
	return nil
}

// DeleteLogpushJob deletes a Logpush Job for a zone.
//
// API reference: https://api.cloudflare.com/#logpush-jobs-delete-logpush-job
func (api *API) DeleteLogpushJob(ctx context.Context, zoneID string, jobID int) error {
	uri := "/zones/" + zoneID + "/logpush/jobs/" + strconv.Itoa(jobID)
	res, err := api.makeRequestContext(ctx, "DELETE", uri, nil)
	if err != nil {
		return errors.Wrap(err, errMakeRequestError)
	}
	var r LogpushJobDetailsResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return errors.Wrap(err, errUnmarshalError)
	}
	return nil
}

// GetLogpushOwnershipChallenge returns ownership challenge.
//
// API reference: https://api.cloudflare.com/#logpush-jobs-get-ownership-challenge
func (api *API) GetLogpushOwnershipChallenge(ctx context.Context, zoneID, destinationConf string) (*LogpushGetOwnershipChallenge, error) {
	uri := "/zones/" + zoneID + "/logpush/ownership"
	res, err := api.makeRequestContext(ctx, "POST", uri, LogpushGetOwnershipChallengeRequest{
		DestinationConf: destinationConf,
	})
	if err != nil {
		return nil, errors.Wrap(err, errMakeRequestError)
	}
	var r LogpushGetOwnershipChallengeResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, errors.Wrap(err, errUnmarshalError)
	}

	if !r.Result.Valid {
		return nil, errors.New(r.Result.Message)
	}

	return &r.Result, nil
}

// ValidateLogpushOwnershipChallenge returns ownership challenge validation result.
//
// API reference: https://api.cloudflare.com/#logpush-jobs-validate-ownership-challenge
func (api *API) ValidateLogpushOwnershipChallenge(ctx context.Context, zoneID, destinationConf, ownershipChallenge string) (bool, error) {
	uri := "/zones/" + zoneID + "/logpush/ownership/validate"
	res, err := api.makeRequestContext(ctx, "POST", uri, LogpushValidateOwnershipChallengeRequest{
		DestinationConf:    destinationConf,
		OwnershipChallenge: ownershipChallenge,
	})
	if err != nil {
		return false, errors.Wrap(err, errMakeRequestError)
	}
	var r LogpushGetOwnershipChallengeResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return false, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result.Valid, nil
}

// CheckLogpushDestinationExists returns destination exists check result.
//
// API reference: https://api.cloudflare.com/#logpush-jobs-check-destination-exists
func (api *API) CheckLogpushDestinationExists(ctx context.Context, zoneID, destinationConf string) (bool, error) {
	uri := "/zones/" + zoneID + "/logpush/validate/destination/exists"
	res, err := api.makeRequestContext(ctx, "POST", uri, LogpushDestinationExistsRequest{
		DestinationConf: destinationConf,
	})
	if err != nil {
		return false, errors.Wrap(err, errMakeRequestError)
	}
	var r LogpushDestinationExistsResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return false, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result.Exists, nil
}
