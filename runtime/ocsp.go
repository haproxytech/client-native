package runtime

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	native_errors "github.com/haproxytech/client-native/v6/errors"
	"github.com/haproxytech/client-native/v6/models"
)

type OcspResponseFmt int

const (
	OcspFmtText OcspResponseFmt = iota
	OcspFmtBase64
)

func parseCertificateIDs(response string) ([]*models.SslCertificateID, error) {
	if response == "" {
		return nil, native_errors.ErrNotFound
	}
	var certificateIDs []*models.SslCertificateID
	parts := strings.Split(response, "\n")
	certificateID := &models.SslCertificateID{
		CertificateID: &models.CertificateID{},
	}
	for _, p := range parts {
		before, after, found := strings.Cut(p, ":")
		if !found {
			continue
		}
		keyString := strings.TrimSpace(before)
		valueString := strings.TrimSpace(after)
		switch key := keyString; key {
		case "Certificate ID key":
			certificateID.CertificateIDKey = valueString
		case "Certificate path":
			certificateID.CertificatePath = valueString
		case "Issuer Name Hash":
			certificateID.CertificateID.IssuerNameHash = valueString
		case "Issuer Key Hash":
			certificateID.CertificateID.IssuerKeyHash = valueString
		case "Serial Number":
			certificateID.CertificateID.SerialNumber = valueString
			certificateIDs = append(certificateIDs, certificateID)
			certificateID = &models.SslCertificateID{
				CertificateID: &models.CertificateID{},
			}
		}
	}

	return certificateIDs, nil
}

// OCSP Response Data:
// OCSP Response Status: successful (0x0)
// Response Type: Basic OCSP Response
// Version: 1 (0x0)
// Responder Id: C = FR, O = HAProxy Technologies, CN = ocsp.haproxy.com
// Produced At: May 27 15:43:38 2021 GMT
// Responses:
// Certificate ID:
// Hash Algorithm: sha1
// Issuer Name Hash: 8A83E0060FAFF709CA7E9B95522A2E81635FDA0A
// Issuer Key Hash: F652B0E435D5EA923851508F0ADBE92D85DE007A
// Serial Number: 100A
// Cert Status: good
// This Update: May 27 15:43:38 2021 GMT
// Next Update: Oct 12 15:43:38 2048 GMT

func parseOcspResponse(response string) (*models.SslOcspResponse, error) {
	response = strings.Trim(response, " \n")
	if response == "" {
		return nil, native_errors.ErrNotFound
	}

	resp := &models.SslOcspResponse{
		Responses: &models.OCSPResponses{
			CertificateID: &models.CertificateID{},
		},
	}

	var parseErr error

	strings.SplitSeq(response, "\n")(func(line string) bool {
		line = strings.TrimSpace(line)
		key, val, found := strings.Cut(line, ": ")
		if !found {
			return true
		}

		switch key {
		case "OCSP Response Status":
			resp.OcspResponseStatus = strings.Split(val, " ")[0]
		case "Response Type":
			resp.ResponseType = val
		case "Version":
			resp.Version = strings.Split(val, " ")[0]
		case "Produced At":
			producedAt, err := time.Parse("Jan 2 15:04:05 2006 GMT", val)
			if err != nil {
				parseErr = fmt.Errorf("cannot parse produced at of ocsp response : %s %w", err.Error(), native_errors.ErrGeneral)
				return false
			}
			resp.ProducedAt = strfmt.Date(producedAt)
		case "Hash Algorithm":
			resp.Responses.CertificateID.HashAlgorithm = val
		case "Serial Number":
			resp.Responses.CertificateID.SerialNumber = val
		case "Issuer Name Hash":
			resp.Responses.CertificateID.IssuerNameHash = val
		case "Issuer Key Hash":
			resp.Responses.CertificateID.IssuerKeyHash = val
		case "Responder Id":
			resp.ResponderID = strings.Split(val, ", ")
		case "This Update":
			thisUpdate, err := time.Parse("Jan 2 15:04:05 2006 GMT", val)
			if err != nil {
				parseErr = fmt.Errorf("cannot parse ocsp response this update : %s %w", err.Error(), native_errors.ErrGeneral)
				return false
			}
			resp.Responses.ThisUpdate = strfmt.Date(thisUpdate)
		case "Next Update":
			nextUpdate, err := time.Parse("Jan 2 15:04:05 2006 GMT", val)
			if err != nil {
				parseErr = fmt.Errorf("cannot parse ocsp response next update : %s %w", err.Error(), native_errors.ErrGeneral)
				return false
			}
			resp.Responses.NextUpdate = strfmt.Date(nextUpdate)
		case "Cert Status":
			resp.Responses.CertStatus = val
		}
		return true
	})

	if parseErr != nil {
		return nil, parseErr
	}

	return resp, nil
}

// OCSP Certid | Path | Next Update | Last Update | Successes | Failures | Last Update Status | Last Update Status (str)
// 303b300906052b0e03021a050004148a83e0060faff709ca7e9b95522a2e81635fda0a0414f652b0e435d5ea923851508f0adbe92d85de007a02021015 |
// /path_to_cert/cert.pem | 30/Jan/2023:00:08:09 +0000 | - | 0 | 1 | 2 | HTTP error

func parseOcspUpdates(response string) ([]*models.SslOcspUpdate, error) {
	if response == "" {
		return nil, native_errors.ErrNotFound
	}
	lines := strings.Split(response, "\n")
	updates := make([]*models.SslOcspUpdate, 0, len(lines))
	for _, line := range lines[1:] {
		fields := strings.Split(line, "|")
		if len(fields) < 8 {
			return nil, fmt.Errorf("%w: error parsing ocsp updates: not enough fields", native_errors.ErrGeneral)
		}
		update := &models.SslOcspUpdate{
			CertID:              strings.TrimSpace(fields[0]),
			Path:                strings.TrimSpace(fields[1]),
			LastUpdateStatusStr: strings.TrimSpace(fields[7]),
		}
		if strings.TrimSpace(fields[2]) != "-" {
			update.NextUpdate = strings.TrimSpace(fields[2])
		}
		if strings.TrimSpace(fields[3]) != "-" {
			update.LastUpdate = strings.TrimSpace(fields[3])
		}
		update.Successes, _ = strconv.ParseInt(strings.TrimSpace(fields[4]), 10, 64)
		update.Failures, _ = strconv.ParseInt(strings.TrimSpace(fields[5]), 10, 64)
		update.LastUpdateStatus, _ = strconv.ParseInt(strings.TrimSpace(fields[6]), 10, 64)
		updates = append(updates, update)
	}
	return updates, nil
}

// ShowOcspResponses returns the IDs of the ocsp responses as well as the corresponding frontend
// certificate's path, the issuer's name and key hash and the serial number of
// the certificate for which the OCSP response was built
func (s *SingleRuntime) ShowOcspResponses() ([]*models.SslCertificateID, error) {
	cmd := "show ssl ocsp-response"
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound)
	}
	return parseCertificateIDs(response)
}

// ShowOcspResponse returns the contents of the corresponding OCSP response.
// Specifying the output format is only allowed when using a certificate ID, not a path!
func (s *SingleRuntime) ShowOcspResponse(idOrPath string, ofmt ...OcspResponseFmt) (*models.SslOcspResponse, error) {
	if idOrPath == "" {
		return nil, fmt.Errorf("%s %w", "No certificate ID or path provided", native_errors.ErrGeneral)
	}

	var format string
	useBase64 := false
	if len(ofmt) > 0 {
		switch ofmt[0] {
		case OcspFmtText:
			format = "text "
		case OcspFmtBase64:
			format = "base64 "
			useBase64 = true
		}
	}

	response, err := s.ExecuteWithResponse("show ssl ocsp-response " + format + idOrPath)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound)
	}

	if useBase64 {
		return &models.SslOcspResponse{
			Base64Response: strings.Trim(response, " \n"),
		}, nil
	}

	ocspResponse, err := parseOcspResponse(response)
	if err != nil {
		return nil, err
	}
	return ocspResponse, nil
}

// ShowOcspUpdates displays the entries concerned by the OCSP update
func (s *SingleRuntime) ShowOcspUpdates() ([]*models.SslOcspUpdate, error) {
	cmd := "show ssl ocsp-updates"
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound)
	}
	return parseOcspUpdates(response)
}

// SetOcspResponse update an OCSP Response for a certificate using a payload
func (s *SingleRuntime) SetOcspResponse(payload string) error {
	if payload == "" {
		return fmt.Errorf("%s %w", "Argument payload empty", native_errors.ErrGeneral)
	}
	cmd := fmt.Sprintf("set ssl ocsp-response <<\n%s\n", payload)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if !strings.Contains(response, "OCSP Response updated") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

// UpdateOcspResponse creates an OCSP request for the specified certFile and send it to the OCSP Responder
// whose URI is specified in the AIA field of the certificate
func (s *SingleRuntime) UpdateOcspResponse(certFile string) (*models.SslOcspResponse, error) {
	if certFile == "" {
		return nil, fmt.Errorf("%s %w", "Argument certFile empty", native_errors.ErrNotFound)
	}
	response, err := s.ExecuteWithResponse("update ssl ocsp-response " + certFile)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound)
	}
	return parseOcspResponse(response)
}
