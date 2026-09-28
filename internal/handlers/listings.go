package handlers

import (
	"database/sql"
	"encoding/json"
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
	Price       string    `json:"price"`
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
