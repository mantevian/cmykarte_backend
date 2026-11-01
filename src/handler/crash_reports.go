package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"regexp"
	"strings"
)

const maxCrashReportBodyBytes = 32 * 1024

var crashReportIDPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

type crashReportRequest struct {
	ReportID     string `json:"report_id"`
	ReportType   string `json:"report_type"`
	Message      string `json:"message"`
	StackTrace   string `json:"stack_trace"`
	AppVersion   string `json:"app_version"`
	UnityVersion string `json:"unity_version"`
	Platform     string `json:"platform"`
	Scene        string `json:"scene"`
}

func HandleCrashReport(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || mediaType != "application/json" {
			http.Error(w, "content type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxCrashReportBodyBytes)
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		var report crashReportRequest
		if err := decoder.Decode(&report); err != nil {
			var maxBytesError *http.MaxBytesError
			if errors.As(err, &maxBytesError) {
				http.Error(w, "crash report is too large", http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w, "invalid crash report", http.StatusBadRequest)
			return
		}
		if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
			http.Error(w, "invalid crash report", http.StatusBadRequest)
			return
		}

		if !crashReportIDPattern.MatchString(report.ReportID) ||
			(report.ReportType != "Error" && report.ReportType != "Assert" &&
				report.ReportType != "Exception" && report.ReportType != "Crash") ||
			!validCrashReportField(report.Message, 1024, true) ||
			!validCrashReportField(report.StackTrace, 4096, false) ||
			!validCrashReportField(report.AppVersion, 128, true) ||
			!validCrashReportField(report.UnityVersion, 128, true) ||
			!validCrashReportField(report.Platform, 128, true) ||
			!validCrashReportField(report.Scene, 256, false) {
			http.Error(w, "invalid crash report", http.StatusBadRequest)
			return
		}

		_, err = db.Exec(
			`insert into crash_reports (
				report_id, report_type, message, stack_trace,
				app_version, unity_version, platform, scene
			) values ($1, $2, $3, $4, $5, $6, $7, $8)
			on conflict (report_id) do nothing`,
			report.ReportID,
			report.ReportType,
			report.Message,
			report.StackTrace,
			report.AppVersion,
			report.UnityVersion,
			report.Platform,
			report.Scene,
		)
		if err != nil {
			log.Printf("failed to persist crash report %s: %v", report.ReportID, err)
			http.Error(w, "could not store crash report", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func validCrashReportField(value string, maxLength int, required bool) bool {
	if len(value) > maxLength || strings.IndexByte(value, 0) >= 0 ||
		(required && strings.TrimSpace(value) == "") {
		return false
	}
	return true
}
