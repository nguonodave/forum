package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorPage(t *testing.T) {
	type args struct {
		errorText  string
		statusCode int
	}

	mockHandler := func(arg args) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			ErrorPage(w, arg.errorText, arg.statusCode)
		}
	}

	tests := []struct {
		name       string
		args       args
		noTemplate bool
	}{
		{
			name: "Error 500",
			args: args{
				errorText:  "Internal Server Error",
				statusCode: 500,
			},
			noTemplate: false,
		},

		{
			name: "Error 404",
			args: args{
				errorText:  "Not found",
				statusCode: 404,
			},
			noTemplate: false,
		},

		{
			name: "Error 403",
			args: args{
				errorText:  "Forbidden",
				statusCode: 403,
			},
			noTemplate: false,
		},

		{
			name: "Error 400",
			args: args{
				errorText:  "Bad request",
				statusCode: 400,
			},
			noTemplate: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				var originalTemplateDir = templatesDir
				if tt.noTemplate {
					templatesDir = ""
					defer func() {
						templatesDir = originalTemplateDir
					}()
				}

				req := httptest.NewRequest("GET", "/", nil)
				w := httptest.NewRecorder()

				// Call the handler
				handler := mockHandler(tt.args)
				handler(w, req)

				if w.Code != tt.args.statusCode {
					t.Errorf("got %d, want %d", w.Code, tt.args.statusCode)
				}
			},
		)
	}
}
