package gdax

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/imdario/mergo"
)

// Formats
const (
	Fills = "fills"
	Pdf   = "pdf"
	Csv   = "csv"
)

// A ReportParams stores the start and end date of a report response.
type ReportParams struct {
	StartDate *time.Time `json:"start_date,string,omitempty"`
	EndDate   *time.Time `json:"end_date,string,omitempty"`
}

// A Report represents a report.
type Report struct {
	Type      string     `json:"type"`
	StartDate *time.Time `json:"start_date,string"`
	EndDate   *time.Time `json:"end_date,string"`
	ProductID string     `json:"product_id,omitempty"`
	AccountID *uuid.UUID `json:"account_id,string,omitempty"`
	Format    string     `json:"format,omitempty"`
	Email     string     `json:"email,omitempty"`

	// response params
	ID          *uuid.UUID    `json:"id,string,omitempty"`
	Status      string        `json:"status,omitempty"`
	CreatedAt   *time.Time    `json:"created_at,string,omitempty"`
	CompletedAt *time.Time    `json:"completed_at,string,omitempty"`
	ExpiresAt   *time.Time    `json:"expires_at,string,omitempty"`
	FileURL     string        `json:"file_url,omitempty"`
	Params      *ReportParams `json:"params,omitempty"`
}

// CreateReport submits a report request.
func (accessInfo *AccessInfo) CreateReport(report *Report) (*Report, error) {
	// POST /reports
	var reportResponse Report
	jsonBytes, err := json.Marshal(*report)
	if err != nil {
		return nil, err
	}
	_, err = accessInfo.request(http.MethodPost, "/reports", string(jsonBytes), &reportResponse)
	if err != nil {
		return nil, err
	}
	if err = mergo.Merge(&reportResponse, *report); err != nil {
		return nil, err
	}

	return &reportResponse, err
}

// GetReportStatus retrieves the status of a submitted report.
func (accessInfo *AccessInfo) GetReportStatus(reportID *uuid.UUID) (*Report, error) {
	// GET /reports/:report_id
	var reportStatus Report
	_, err := accessInfo.request(http.MethodGet, "/reports/"+reportID.String(), "", &reportStatus)
	if err != nil {
		return nil, err
	}
	return &reportStatus, nil
}
