package utils

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockResponseWriter struct {
	HeaderMap http.Header
	Body      *bytes.Buffer
	Code      int
}

func (m *MockResponseWriter) Header() http.Header {
	return m.HeaderMap
}

func (m *MockResponseWriter) Write(data []byte) (int, error) {
	return m.Body.Write(data)
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.Code = statusCode
}

func TestJSONResponse(t *testing.T) {
	// Create a mock response writer
	mockWriter := &MockResponseWriter{
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
	}

	// Sample data to return
	data := map[string]string{"message": "Success"}

	// Call the JSONResponse function
	JSONResponse(mockWriter, http.StatusOK, data)

	// Assert the status code is correct
	assert.Equal(t, http.StatusOK, mockWriter.Code)

	// Assert the content-type header is set to "application/json"
	assert.Equal(t, "application/json", mockWriter.Header().Get("Content-Type"))

	// Assert the body contains the correct JSON
	expectedBody := `{"message":"Success"}`
	assert.JSONEq(t, expectedBody, mockWriter.Body.String())
}

func TestJSONError(t *testing.T) {
	// Create a mock response writer
	mockWriter := &MockResponseWriter{
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
	}

	// Call the JSONError function
	JSONError(mockWriter, http.StatusBadRequest, "Invalid request")

	// Assert the status code is correct
	assert.Equal(t, http.StatusBadRequest, mockWriter.Code)

	// Assert the content-type header is set to "application/json"
	assert.Equal(t, "application/json", mockWriter.Header().Get("Content-Type"))

	// Assert the body contains the correct JSON error message
	expectedBody := `{"error":"Invalid request"}`
	assert.JSONEq(t, expectedBody, mockWriter.Body.String())
}
