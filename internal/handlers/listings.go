package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/abhishek8841/go-monolith/internal/httpx"
	"github.com/abhishek8841/go-monolith/internal/middleware"
)

type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListingHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		db:     db,
		logger: logger,
	}
}

// always pass context in things like api or requests like we did below

// using this constructor pattern we dont have to do dependency injection of dependencies
// so we can remove the envelop function we used from over the http.HandleFunc type of funcs
func (lh ListingHandler) List(w http.ResponseWriter, r *http.Request) {
	// request scoped context is a context that web server creates for us and calls the cancel function on it when request gets cancelled
	ctx := r.Context()
	// we control zombie query using this
	rows, err := lh.db.QueryContext(ctx,
		`SELECT id, title, description, price, city, created_at 
		FROM LISTINGS
		ORDER BY created_at DESC 
		LIMIT 100`)
	if err != nil {
		lh.logger.Info("listings query error", "err", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	listings := []listing{}
	for rows.Next() {
		var l listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			lh.logger.Error("rows scan error", "err", err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
		listings = append(listings, l)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Rows Err: %v", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	// w.Write([]byte(`{"status": "endpoint working"}`))

	_ = json.NewEncoder(w).Encode(listings)
}

func (lh ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)

	id := r.PathValue("id")

	_, err := lh.db.ExecContext(ctx, `DELETE FROM LISTINGS WHERE id = $1`, id)
	if err != nil {
		log.Printf("Delete: %v", err)
		// in form of key value pair after the init message
		lh.logger.Error("delete failed", "listing_id", id, "request_id", requestId, "err", err)
		// http.Error(w, "Internal Error", http.StatusInternalServerError)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// in slog.Debug we try to log as much as we can so thats why the default is info level so as to ignore the the hefty .Debug logs during peace times
// in go we have nominal typing i.e. directly passing (in the 4th arg of httpx.Error) "code_constant" can slip as Code type but var temp string = "code_constant gives error"

func (lh ListingHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)

	var req CreateListingRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		lh.logger.Error("failed to decode", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.CodeMalformedJson)
		return
	}

	if err := req.Validate(); err != nil {
		var verr *ValidationError
		errors.As(err, &verr)
		httpx.ValidationError(w, http.StatusUnprocessableEntity, err.Error(), httpx.CodeValidationError, verr.Field)
		return
	}

	row := lh.db.QueryRowContext(ctx, `
	INSERT INTO LISTINGS (title, description, price, city) VALUES ($1, $2, $3, $4) RETURNING id, title, created_at`, req.Title, req.Description, req.Price, req.City)
	// var id string
	var out CreateListingResponse
	if err := row.Scan(&out.ID, &out.Title, &out.CreatedAt); err != nil {
		lh.logger.Error("failed to insert", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	lh.logger.Info("listing created", "request_id", requestId, "listing_id", out.ID)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(out)

}
