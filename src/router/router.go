package router

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/praction-networks/quantum-ISP365/webapp/src/handler"
	middlewares "github.com/praction-networks/quantum-ISP365/webapp/src/middleware"
)

func LoadRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middlewares.MaxBodySizeMiddleware(20 * 1024 * 1024))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json", "multipart/form-data"))
	r.Use(middleware.CleanPath)

	// r.Use(MethodRestriction)
	r.Use(middlewares.CORS)

	// Grouping all routes under /web/v1
	r.Route("/web/v1", func(r chi.Router) {
		// Public routes
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			response := map[string]string{"status": "UP"}
			json.NewEncoder(w).Encode(response)
		})

		authHandler := &handler.User{}
		r.Post("/login", authHandler.Login)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middlewares.JWTAuthMiddleware) // Apply JWTAuthMiddleware to all routes in this group

			planHandler := &handler.Plan{}
			r.Post("/plan", planHandler.CreatePlan)
			r.Get("/plans", planHandler.GetAllPlan)

			jonHandler := &handler.JobHandler{}

			r.Post("/job", jonHandler.CreateJobHandler)
			r.Get("/jobs", jonHandler.GetAllJobs)

			userIntrestHandler := &handler.UserIntrest{}

			r.Post("/customer/intrest", userIntrestHandler.ShowIntresrtHandler)
			r.Post("/customer/intrest/otp", userIntrestHandler.VerifyUserOTP)
			r.Post("/customer/intrest/resend", userIntrestHandler.ResendUserOTP)
			r.Get("/customer/intrest", userIntrestHandler.GetALl)

			userReferalHandler := &handler.UserReferal{}

			r.Post("/customer/referral", userReferalHandler.ReferUser)
			r.Post("/customer/referral/otp", userReferalHandler.VerifyUserOTP)
			r.Post("/customer/referral/resend", userReferalHandler.ResendUserOTP)
			r.Get("/customer/referral", userReferalHandler.GetALl)

			blogHandler := &handler.BlogHandler{}

			r.Post("/blog", blogHandler.CreateBlogHandler)
			r.Get("/blog", blogHandler.GetBlogHandler)
			r.Get("/blog/{id}", blogHandler.GetOneBlogHandler)

			blogCategoryHandler := &handler.BlogCategoryHandler{}

			r.Post("/blog/category", blogCategoryHandler.CreateBlogCategoryHandler)
			r.Get("/blog/category", blogCategoryHandler.GetAllBlogCategoryHandler)

			blogTagHandler := &handler.BlogTagHandler{}

			r.Post("/blog/tag", blogTagHandler.CreateBlogTagHandler)
			r.Get("/blog/tag", blogTagHandler.GetAllBlogTagHandler)

			// Image upload route
			imageUploadHandler := &handler.ImageUploadHandler{}
			r.Post("/blog/image/upload", imageUploadHandler.UploadImage)

			blogCommentsHandler := &handler.BlogCommentsHandler{}

			r.Post("/blog/comments", blogCommentsHandler.CreateBlogCommentsHandler)

		})
	})

	// Handle OPTIONS requests (CORS preflight)
	r.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Strict-Auth-Key")
		w.Header().Set("Access-Control-Max-Age", "3600") // Optional: caches the preflight response for 1 hour
		w.WriteHeader(http.StatusNoContent)              // 204 No Content
	})

	return r
}
