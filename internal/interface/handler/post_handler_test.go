package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"my-blog-engine/internal/domain/entity"
	"my-blog-engine/internal/usecase"
)

type recordingPostUseCase struct {
	usecase.PostUseCase
	gotLimit, gotOffset int
}

func (u *recordingPostUseCase) List(_ context.Context, limit, offset int) ([]*entity.Post, int, error) {
	u.gotLimit, u.gotOffset = limit, offset
	return nil, 0, nil
}

func (u *recordingPostUseCase) ListPublished(_ context.Context, limit, offset int) ([]*entity.Post, int, error) {
	u.gotLimit, u.gotOffset = limit, offset
	return nil, 0, nil
}

func TestPostHandlerListOffset(t *testing.T) {
	cases := []struct {
		name  string
		query string
		want  int
	}{
		{"negative offset becomes zero", "?offset=-1", 0},
		{"large negative offset becomes zero", "?offset=-999999", 0},
		{"missing offset is zero", "", 0},
		{"non-numeric offset is zero", "?offset=abc", 0},
		{"zero offset is kept", "?offset=0", 0},
		{"positive offset is kept", "?offset=40", 40},
	}
	handlers := map[string]func(*PostHandler) http.HandlerFunc{
		"List":          func(h *PostHandler) http.HandlerFunc { return h.List },
		"ListPublished": func(h *PostHandler) http.HandlerFunc { return h.ListPublished },
	}
	for hname, pick := range handlers {
		for _, tc := range cases {
			t.Run(hname+"/"+tc.name, func(t *testing.T) {
				uc := &recordingPostUseCase{gotOffset: -12345}
				rec := httptest.NewRecorder()

				pick(NewPostHandler(uc))(rec, httptest.NewRequest(http.MethodGet, "/posts"+tc.query, nil))

				if rec.Code != http.StatusOK {
					t.Fatalf("status = %d, want 200", rec.Code)
				}
				if uc.gotOffset != tc.want {
					t.Errorf("offset passed to use case = %d, want %d", uc.gotOffset, tc.want)
				}
				var body struct {
					Offset int `json:"offset"`
				}
				if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if body.Offset != tc.want {
					t.Errorf("offset in response = %d, want %d", body.Offset, tc.want)
				}
			})
		}
	}
}
