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
		r.Post("/logout", authHandler.Logout)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middlewares.JWTAuthMiddleware) // Apply JWTAuthMiddleware to all routes in this group
			r.Post("/user/set_password", authHandler.SetPassword)
			r.Post("/user", http.HandlerFunc(
				middlewares.JWTAuthMiddleware(
					middlewares.RoleAuthMiddleware([]string{"admin"}, http.HandlerFunc(authHandler.Create)),
				).ServeHTTP,
			))

			r.Get("/user", http.HandlerFunc(
				middlewares.JWTAuthMiddleware(
					middlewares.RoleAuthMiddleware([]string{"admin"}, http.HandlerFunc(authHandler.GetAll)),
				).ServeHTTP,
			))

			r.Get("/user/{id}", http.HandlerFunc(
				middlewares.JWTAuthMiddleware(
					middlewares.RoleAuthMiddleware([]string{"admin"}, http.HandlerFunc(authHandler.GetOne)),
				).ServeHTTP,
			))
			r.Post("/user/{id}/password", http.HandlerFunc(
				middlewares.JWTAuthMiddleware(
					middlewares.RoleAuthMiddleware([]string{"admin"}, http.HandlerFunc(authHandler.SetPasswordByAdmin)),
				).ServeHTTP,
			))

			r.Delete("/user", http.HandlerFunc(
				middlewares.JWTAuthMiddleware(
					middlewares.RoleAuthMiddleware([]string{"admin"}, http.HandlerFunc(authHandler.Create)),
				).ServeHTTP,
			))

			planHandler := &handler.Plan{}
			r.Post("/plan", planHandler.CreatePlan)
			r.Get("/plans", planHandler.GetAllPlan)
			r.Get("/plan/{uuid}", planHandler.GetOne)

			unAvailableAreaHandler := &handler.UnAvailableAreaRegister{}
			r.Post("/customer/unavailable", unAvailableAreaHandler.UnavailavleAreaUserIntrest)
			r.Get("/customer/unavailable", unAvailableAreaHandler.GetALl)

			planRegisterHandler := &handler.PlanRegister{}
			r.Post("/plan/register", planRegisterHandler.RegisterUserIntrestForPlan)
			r.Post("/plan/register/otp", planRegisterHandler.VerifyUserOTPForPlan)
			r.Post("/plan/register/resend", planRegisterHandler.ResendUserOTPForPlan)
			r.Get("/plans/register", planRegisterHandler.GetALl)

			jobHandler := &handler.JobHandler{}

			r.Post("/job", jobHandler.CreateJobHandler)
			r.Get("/jobs", jobHandler.GetAllJobs)
			r.Get("/job/{id}", jobHandler.GetOneJobandler)

			userIntrestHandler := &handler.UserIntrest{}

			r.Post("/customer/intrest", userIntrestHandler.ShowIntresrtHandler)
			r.Patch("/customer/intrest/id", userIntrestHandler.UpdateStatus)
			r.Post("/customer/intrest/otp", userIntrestHandler.VerifyUserOTP)
			r.Post("/customer/intrest/resend", userIntrestHandler.ResendUserOTP)
			r.Get("/customer/intrest", userIntrestHandler.GetALl)

			userReferalHandler := &handler.UserReferal{}

			r.Post("/customer/referral", userReferalHandler.ReferUser)
			r.Post("/customer/referral/otp", userReferalHandler.VerifyUserOTP)
			r.Post("/customer/referral/resend", userReferalHandler.ResendUserOTP)
			r.Get("/customer/referral", userReferalHandler.GetALl)

			blogHandler := &handler.BlogHandler{}

			r.Post("/blog", http.HandlerFunc(
				middlewares.JWTAuthMiddleware(
					middlewares.RoleAuthMiddleware([]string{"admin", "manager", "user"}, http.HandlerFunc(blogHandler.CreateBlogHandler)),
				).ServeHTTP,
			))
			r.Get("/blog", blogHandler.GetBlogHandler)
			r.Get("/blog/{id}", blogHandler.GetOneBlogHandler)
			r.Post("/blog/view", blogHandler.AddView)
			r.Post("/blog/share", blogHandler.AddShare)
			// Blog update route with role-based access control
			r.Put("/blog/{id}", http.HandlerFunc(
				middlewares.JWTAuthMiddleware(
					middlewares.RoleAuthMiddleware([]string{"admin", "manager", "user"}, http.HandlerFunc(blogHandler.UpdateBlogHandler)),
				).ServeHTTP,
			))
			r.Delete("/blog/{id}", http.HandlerFunc(
				middlewares.JWTAuthMiddleware(
					middlewares.RoleAuthMiddleware([]string{"admin"}, http.HandlerFunc(blogHandler.DeleteBlogHandler)),
				).ServeHTTP,
			))
			r.Get("/blog/all", http.HandlerFunc(
				middlewares.JWTAuthMiddleware(
					middlewares.RoleAuthMiddleware([]string{"admin", "manager", "user"}, http.HandlerFunc(blogHandler.GeAlltBlogHandler)),
				).ServeHTTP,
			))
			r.Put("/blog/approve/{id}",
				http.HandlerFunc(
					middlewares.JWTAuthMiddleware(
						middlewares.RoleAuthMiddleware([]string{"admin", "manager"}, http.HandlerFunc(blogHandler.ApproveBlogHandler)),
					).ServeHTTP,
				))
			r.Put("/blog/publish/{id}", http.HandlerFunc(
				middlewares.JWTAuthMiddleware(
					middlewares.RoleAuthMiddleware([]string{"admin"}, http.HandlerFunc(blogHandler.PublishBlogHandler)),
				).ServeHTTP,
			))

			blogCategoryHandler := &handler.BlogCategoryHandler{}

			r.Post("/blog/category",
				http.HandlerFunc(
					middlewares.JWTAuthMiddleware(
						middlewares.RoleAuthMiddleware([]string{"admin", "manager", "user"}, http.HandlerFunc(blogCategoryHandler.CreateBlogCategoryHandler)),
					).ServeHTTP,
				))
			r.Get("/blog/category", http.HandlerFunc(
				middlewares.JWTAuthMiddleware(
					middlewares.RoleAuthMiddleware([]string{"admin", "manager", "user"}, http.HandlerFunc(blogCategoryHandler.GetAllBlogCategoryHandler)),
				).ServeHTTP,
			))

			r.Delete("/blog/category/{id}", http.HandlerFunc(
				middlewares.JWTAuthMiddleware(
					middlewares.RoleAuthMiddleware([]string{"admin", "manager", "user"}, http.HandlerFunc(blogCategoryHandler.DeleteBlogCatgoryByID)),
				).ServeHTTP,
			))

			blogTagHandler := &handler.BlogTagHandler{}

			r.Post("/blog/tag",
				http.HandlerFunc(
					middlewares.JWTAuthMiddleware(
						middlewares.RoleAuthMiddleware([]string{"admin", "manager", "user"}, http.HandlerFunc(blogTagHandler.CreateBlogTagHandler)),
					).ServeHTTP,
				))

			r.Delete("/blog/tag/{id}",
				http.HandlerFunc(
					middlewares.JWTAuthMiddleware(
						middlewares.RoleAuthMiddleware([]string{"admin", "manager", "user"}, http.HandlerFunc(blogTagHandler.DeleteBlogTag)),
					).ServeHTTP,
				))

			r.Get("/blog/tag", http.HandlerFunc(
				middlewares.JWTAuthMiddleware(
					middlewares.RoleAuthMiddleware([]string{"admin", "manager", "user"}, http.HandlerFunc(blogTagHandler.GetAllBlogTagHandler)),
				).ServeHTTP,
			))

			// Image upload route
			imageUploadHandler := handler.ImageUploadHandler{}
			r.Post("/image/upload", http.HandlerFunc(
				middlewares.JWTAuthMiddleware(
					middlewares.RoleAuthMiddleware([]string{"admin", "manager", "user"}, http.HandlerFunc(imageUploadHandler.UploadImage)),
				).ServeHTTP,
			))

			r.Get("/image", imageUploadHandler.GetAllImage)
			r.Delete("/image/{id}", imageUploadHandler.DeleteImage)

			blogCommentsHandler := &handler.BlogCommentsHandler{}

			r.Post("/blog/comments", blogCommentsHandler.CreateBlogCommentsHandler)

			// Service Area

			serviceAreaHandler := &handler.ServiceAreaHandler{}

			r.Post("/servicearea", serviceAreaHandler.CreateServiceArea)
			r.Post("/servicearea/check", serviceAreaHandler.CheckServiceArea)
			r.Get("/servicearea", serviceAreaHandler.GetAllServiceArea)
			r.Patch("/servicearea/{id}", serviceAreaHandler.UpdateServiceArea)
			r.Patch("/servicearea/area/{id}", serviceAreaHandler.ModifyServiceArea)

		})
	})

	// Handle OPTIONS requests (CORS preflight)
	r.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT,PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Strict-Auth-Key")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
	})

	return r
}

func WithRoleAuth(requiredRoles []string, next http.Handler) http.Handler {
	return middlewares.JWTAuthMiddleware(
		middlewares.RoleAuthMiddleware(requiredRoles, next),
	)
}
