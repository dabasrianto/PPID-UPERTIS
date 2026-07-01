package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	_ "github.com/joho/godotenv/autoload"
)

// validateEnv checks all critical environment variables at startup.
// App will fail fast with clear log if any are missing.
func validateEnv() {
	required := []string{"DB_HOST", "DB_NAME", "JWT_SECRET", "PORT"}
	missing := []string{}
	for _, key := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		log.Fatalf("FATAL: Missing required environment variables: %v\nSet them in backend/.env", missing)
	}
}

// appReady indicates whether the app is fully initialized and ready to serve traffic
var appReady bool

func main() {
	// Fail fast if critical ENV missing
	validateEnv()

	// Initialize structured logger first
	initLogger()

	initDB()
	initJWTSecret()
	StartBackgroundPostSync()
	initWhatsApp()

	// Super admin credentials from environment variables
	adminEmail := getEnv("ADMIN_EMAIL", "")
	adminPassword := getEnv("ADMIN_PASSWORD", "")
	adminName := getEnv("ADMIN_NAME", "Super Admin")
	if adminEmail != "" && adminPassword != "" {
		if err := ensureSuperAdmin(adminEmail, adminPassword, adminName); err != nil {
			zlog.Fatal().Err(err).Msg("Failed to ensure super admin")
		}
	} else {
		zlog.Info().Msg("ADMIN_EMAIL/ADMIN_PASSWORD not set. Skipping super admin creation.")
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		BodyLimit:    50 * 1024 * 1024, // Allow up to 50MB file uploads
		ProxyHeader:  "CF-Connecting-IP",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			rid, _ := c.Locals(RequestIDKey).(string)
			return c.Status(code).JSON(fiber.Map{
				"error":      err.Error(),
				"request_id": rid,
			})
		},
	})

	// Middleware stack (order matters)
	app.Use(requestIDMiddleware)       // 1. Generate request_id
	app.Use(panicRecoverMiddleware)    // 2. Catch panics (logs internally, generic to client)
	app.Use(structuredLoggerMiddleware) // 3. Structured JSON logging
	// 4. Security headers
	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		return c.Next()
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     getEnv("CORS_ORIGINS", "http://localhost:5173,http://localhost:5174,http://localhost:5175"),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
	}))
	app.Static("/uploads", filepath.Join("..", "public", "uploads"))

	app.Get("/robots.txt", func(c *fiber.Ctx) error {
		domain := getEnv("SITE_DOMAIN", "https://upertis.ac.id")
		c.Set("Content-Type", "text/plain")
		return c.SendString(fmt.Sprintf("User-agent: *\nAllow: /\nDisallow: /admin\nDisallow: /api/v1/admin/\nDisallow: /api/v1/auth/\nSitemap: %s/sitemap.xml\n", domain))
	})

	// Sitemap.xml — dynamic from all public content
	app.Get("/sitemap.xml", func(c *fiber.Ctx) error {
		domain := getEnv("SITE_DOMAIN", "https://upertis.ac.id")

		var sb strings.Builder
		sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
		sb.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

		// Static pages
		sb.WriteString(fmt.Sprintf(`<url><loc>%s/</loc><changefreq>daily</changefreq><priority>1.0</priority></url>`, domain))
		sb.WriteString(fmt.Sprintf(`<url><loc>%s/posts</loc><changefreq>daily</changefreq><priority>0.8</priority></url>`, domain))
		sb.WriteString(fmt.Sprintf(`<url><loc>%s/berita</loc><changefreq>daily</changefreq><priority>0.8</priority></url>`, domain))
		sb.WriteString(fmt.Sprintf(`<url><loc>%s/events</loc><changefreq>weekly</changefreq><priority>0.7</priority></url>`, domain))
		sb.WriteString(fmt.Sprintf(`<url><loc>%s/fakultas</loc><changefreq>weekly</changefreq><priority>0.8</priority></url>`, domain))
		sb.WriteString(fmt.Sprintf(`<url><loc>%s/download</loc><changefreq>weekly</changefreq><priority>0.7</priority></url>`, domain))
		sb.WriteString(fmt.Sprintf(`<url><loc>%s/media</loc><changefreq>weekly</changefreq><priority>0.7</priority></url>`, domain))
		sb.WriteString(fmt.Sprintf(`<url><loc>%s/pmb</loc><changefreq>weekly</changefreq><priority>0.8</priority></url>`, domain))

		// Posts (blog/articles)
		rows, err := db.QueryContext(context.Background(),
			`SELECT slug, COALESCE(updated_at, created_at) FROM posts WHERE status='published' AND deleted_at IS NULL ORDER BY updated_at DESC LIMIT 500`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var slug, updatedAt string
				if rows.Scan(&slug, &updatedAt) == nil {
					lastmod := ""
					if len(updatedAt) >= 10 {
						lastmod = fmt.Sprintf(`<lastmod>%s</lastmod>`, updatedAt[:10])
					}
					sb.WriteString(fmt.Sprintf(`<url><loc>%s/posts/%s</loc>%s<changefreq>weekly</changefreq><priority>0.6</priority></url>`, domain, slug, lastmod))
				}
			}
		}

		// News (berita)
		rows2, err := db.QueryContext(context.Background(),
			`SELECT id, COALESCE(updated_at, created_at) FROM news WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT 500`)
		if err == nil {
			defer rows2.Close()
			for rows2.Next() {
				var id, updatedAt string
				if rows2.Scan(&id, &updatedAt) == nil {
					lastmod := ""
					if len(updatedAt) >= 10 {
						lastmod = fmt.Sprintf(`<lastmod>%s</lastmod>`, updatedAt[:10])
					}
					sb.WriteString(fmt.Sprintf(`<url><loc>%s/berita/%s</loc>%s<changefreq>weekly</changefreq><priority>0.6</priority></url>`, domain, id, lastmod))
				}
			}
		}

		// Pages (halaman)
		rows3, err := db.QueryContext(context.Background(),
			`SELECT slug, COALESCE(updated_at, created_at) FROM pages WHERE deleted_at IS NULL ORDER BY updated_at DESC LIMIT 100`)
		if err == nil {
			defer rows3.Close()
			for rows3.Next() {
				var slug, updatedAt string
				if rows3.Scan(&slug, &updatedAt) == nil {
					lastmod := ""
					if len(updatedAt) >= 10 {
						lastmod = fmt.Sprintf(`<lastmod>%s</lastmod>`, updatedAt[:10])
					}
					sb.WriteString(fmt.Sprintf(`<url><loc>%s/halaman/%s</loc>%s<changefreq>monthly</changefreq><priority>0.5</priority></url>`, domain, slug, lastmod))
				}
			}
		}

		// Events
		rows4, err := db.QueryContext(context.Background(),
			`SELECT slug, COALESCE(updated_at, created_at) FROM campus_events WHERE deleted_at IS NULL ORDER BY event_date DESC LIMIT 200`)
		if err == nil {
			defer rows4.Close()
			for rows4.Next() {
				var slug, updatedAt string
				if rows4.Scan(&slug, &updatedAt) == nil {
					lastmod := ""
					if len(updatedAt) >= 10 {
						lastmod = fmt.Sprintf(`<lastmod>%s</lastmod>`, updatedAt[:10])
					}
					sb.WriteString(fmt.Sprintf(`<url><loc>%s/events/%s</loc>%s<changefreq>weekly</changefreq><priority>0.5</priority></url>`, domain, slug, lastmod))
				}
			}
		}

		// Faculties
		rows5, err := db.QueryContext(context.Background(),
			`SELECT slug FROM faculties WHERE deleted_at IS NULL ORDER BY name`)
		if err == nil {
			defer rows5.Close()
			for rows5.Next() {
				var slug string
				if rows5.Scan(&slug) == nil {
					sb.WriteString(fmt.Sprintf(`<url><loc>%s/fakultas/%s</loc><changefreq>monthly</changefreq><priority>0.6</priority></url>`, domain, slug))
				}
			}
		}

		// Study Programs (Prodi)
		rows6, err := db.QueryContext(context.Background(),
			`SELECT slug FROM faculty_programs WHERE active = true AND slug IS NOT NULL ORDER BY name`)
		if err == nil {
			defer rows6.Close()
			for rows6.Next() {
				var slug string
				if rows6.Scan(&slug) == nil {
					sb.WriteString(fmt.Sprintf(`<url><loc>%s/prodi/%s</loc><changefreq>monthly</changefreq><priority>0.7</priority></url>`, domain, slug))
				}
			}
		}

		sb.WriteString(`</urlset>`)

		c.Set("Content-Type", "application/xml")
		c.Set("Cache-Control", "public, max-age=3600")
		return c.SendString(sb.String())
	})

	// Health check — checks DB connection
	app.Get("/health", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy",
				"db":     "disconnected",
			})
		}
		return c.JSON(fiber.Map{
			"status": "ok",
			"db":     "connected",
		})
	})

	// Readiness check — for reverse proxy / load balancer
	app.Get("/ready", func(c *fiber.Ctx) error {
		if !appReady {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not_ready",
			})
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not_ready",
				"reason": "db_unavailable",
			})
		}
		return c.JSON(fiber.Map{"status": "ready"})
	})

	// API routes
	api := app.Group("/api/v1")
	setupRoutes(api)

	// Serve static files from frontend build (SPA)
	app.Static("/", "../dist")

	// Dynamic SEO fallback for all page requests (direct browser entry, bot crawlers)
	app.Get("/*", handlePageRequest)

	// Start server
	port := os.Getenv("PORT")

	// Graceful shutdown with timeout
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan
		zlog.Info().Str("signal", sig.String()).Msg("Shutting down gracefully (10s timeout)")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		_ = app.ShutdownWithContext(shutdownCtx)
		_ = db.Close()
		zlog.Info().Msg("Shutdown complete")
		os.Exit(0)
	}()

	// Mark app as ready
	appReady = true
	startScheduledPublisher()
	startBatchAutoClose()
	ensureMultisiteColumns()
	migrateSitesTable()
	migrateLiveUsersTable()
	startSessionCleanup()
	syncIndexHTMLMetadata()
	zlog.Info().Str("port", port).Msg("Server starting")

	// Retry listen with backoff (handles port still in use during restart)
	var listenErr error
	for attempt := 1; attempt <= 5; attempt++ {
		listenErr = app.Listen(":" + port)
		if listenErr == nil {
			break
		}
		if attempt < 5 {
			zlog.Warn().Int("attempt", attempt).Err(listenErr).Msg("Port busy, retrying in 3s...")
			time.Sleep(3 * time.Second)
		}
	}
	if listenErr != nil {
		log.Fatal("Failed to start server after 5 attempts: ", listenErr)
	}
}

func setupRoutes(api fiber.Router) {
	// Auth routes with strict rate limiting (10 req/min)
	authLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Terlalu banyak percobaan. Silakan coba lagi dalam 1 menit.",
			})
		},
	})

	// Public API rate limiter (200 req/min per IP)
	// Homepage loads ~10 endpoints simultaneously, so limit must accommodate SPA behavior
	publicLimiter := limiter.New(limiter.Config{
		Max:        200,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded. Try again later.",
			})
		},
	})

	// Sensitive action limiter (5 req/min — OTP, contact, chat)
	sensitiveLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Terlalu banyak request. Coba lagi nanti.",
			})
		},
	})

	api.Post("/auth/login", authLimiter, loginWithOTP)
	api.Post("/auth/login/verify-otp", authLimiter, verifyLoginOTP)
	api.Post("/auth/login/resend-otp", authLimiter, resendLoginOTP)
	api.Post("/auth/register", authLimiter, registerUser)
	api.Post("/auth/register/verify-otp", authLimiter, verifyRegisterOTP)
	api.Post("/auth/register/resend-otp", authLimiter, resendRegisterOTP)
	api.Post("/auth/logout", logout)
	api.Post("/auth/register-candidate", authLimiter, registerCandidate)
	api.Post("/auth/google", authLimiter, loginWithGoogle)
	api.Get("/auth/me", authMiddleware, getMe)

	// Public routes (with general rate limit)
	api.Get("/hero-slides", publicLimiter, getHeroSlides)
	api.Get("/faculties", publicLimiter, getFaculties)
	api.Get("/faculties/:slug", publicLimiter, getFacultyBySlug)
	api.Get("/programs", publicLimiter, getPrograms)
	api.Get("/programs/:id", publicLimiter, getProgramByID)
	api.Get("/lecturers", publicLimiter, getLecturers)
	api.Get("/lecturers/:id", publicLimiter, getLecturerByID)
	api.Get("/pddikti/search", publicLimiter, pddiktiSearchPT)
	api.Get("/news", publicLimiter, getNews)
	api.Get("/news/:id", publicLimiter, getNewsByID)
	api.Get("/testimonials", publicLimiter, getTestimonials)
	api.Get("/events", publicLimiter, getEvents)
	api.Get("/events/:id", publicLimiter, getEventBySlug)
	api.Post("/events/:id/register", sensitiveLimiter, registerForEvent)
	api.Get("/blog", publicLimiter, getBlogPosts)
	api.Get("/blog/:slug", publicLimiter, func(c *fiber.Ctx) error {
		// Backward compatibility: redirect to posts endpoint
		slug := c.Params("slug")
		return c.Redirect("/api/v1/posts/"+slug, fiber.StatusMovedPermanently)
	})
	api.Get("/pages", publicLimiter, getPages)
	api.Get("/pages/:slug", publicLimiter, getPageBySlug)
	api.Get("/popups", publicLimiter, getPopups)
	api.Get("/post-categories", publicLimiter, getPostCategories)
	api.Get("/posts", publicLimiter, getPosts)
	api.Get("/posts/:slug", publicLimiter, getPostBySlug)
	api.Get("/downloads", publicLimiter, getDownloads)
	api.Post("/downloads/:id/increment", publicLimiter, incrementDownloadCount)
	api.Get("/faqs", publicLimiter, getFAQs)
	api.Get("/gallery", publicLimiter, getGallery)

	api.Get("/tools/usage", publicLimiter, getToolUsages)
	api.Post("/tools/:slug/increment", publicLimiter, incrementToolUsage)
	api.Post("/contact", sensitiveLimiter, submitContact)
	api.Post("/chat", sensitiveLimiter, handleChat)
	api.Get("/settings", publicLimiter, getSettings)
	api.Get("/settings/chatbot", publicLimiter, getChatbotPublicConfig)
	api.Get("/comments", publicLimiter, getComments)
	api.Get("/comments/count", publicLimiter, getCommentCount)
	api.Post("/comments", authMiddleware, sensitiveLimiter, submitComment)
	api.Get("/user/messages", authMiddleware, getUserMessages)
	api.Post("/user/messages/:id/reply", authMiddleware, replyMessageUser)
	api.Post("/user/messages/mark-read", authMiddleware, markUserMessagesRead)
	
	// WhatsApp & OTP routes (strict rate limit)
	api.Get("/whatsapp/status", publicLimiter, getWhatsAppStatusPublic)
	api.Get("/whatsapp/qr", publicLimiter, getWhatsAppQR)
	api.Get("/site-config", publicLimiter, getSiteConfig)
	api.Post("/permohonan-informasi", sensitiveLimiter, submitPermohonan)
	api.Get("/permohonan-informasi/status/:ticket", publicLimiter, checkPermohonanStatus)
	api.Post("/permohonan-informasi/upload", publicLimiter, uploadPermohonanAttachment)
	api.Post("/visitor-ping", publicLimiter, visitorPing)
	api.Get("/visitor-counter.svg", publicLimiter, getVisitorCounterSVG)
	api.Post("/keepalive", authMiddleware, keepalive)
	api.Post("/otp/request", authLimiter, requestOTP)
	api.Post("/otp/verify", authLimiter, verifyOTP)

	// PMB routes (protected for candidate)
	pmb := api.Group("/pmb", authMiddleware, requireRole("candidate", "admin"))
	pmb.Get("/candidate-profile", getCandidateProfile)
	pmb.Post("/candidate-form", saveCandidateForm)
	pmb.Post("/payment", uploadPaymentProof)
	pmb.Post("/upload", candidateUploadFile)
	pmb.Post("/documents", uploadDocument)
	pmb.Get("/documents", getCandidateDocuments)
	pmb.Post("/documents/submit", submitAllDocuments)

	// Admin routes (protected)
	admin := api.Group("/admin", authMiddleware)

	// Profile routes (any authenticated admin/faculty_admin)
	admin.Get("/profile", getAdminProfile)
	admin.Put("/profile", updateAdminProfile)
	admin.Put("/change-password", sensitiveLimiter, changePassword)
	admin.Get("/security/devices", getTrustedDevices)
	admin.Delete("/security/devices", revokeAllTrustedDevices)
	
	// Only super admin
	superAdminAuth := requireRole("admin")
	admin.Get("/security/audit-log", superAdminAuth, getAuthAuditLog)
	admin.Get("/comments", superAdminAuth, getAdminComments)
	admin.Put("/comments/:id", superAdminAuth, moderateComment)
	admin.Get("/messages", superAdminAuth, getMessages)
	admin.Post("/messages/:id/reply", superAdminAuth, replyMessage)
	admin.Put("/messages/:id/read", superAdminAuth, markMessageRead)
	admin.Delete("/messages/:id", superAdminAuth, deleteMessage)
	admin.Get("/users", superAdminAuth, getUsers)
	admin.Post("/users", superAdminAuth, createUser)
	admin.Put("/users/:id", superAdminAuth, updateUser)
	admin.Delete("/users/:id", superAdminAuth, deleteUser)
	admin.Get("/settings", superAdminAuth, getAdminSettings)
	admin.Put("/settings", superAdminAuth, updateSettings)
	admin.Post("/chatbot-models", superAdminAuth, fetchChatbotModels)
	admin.Post("/chatbot-knowledge", superAdminAuth, handleKnowledgeUpload)
	admin.Get("/otp-logs", superAdminAuth, getOTPLogs)
	admin.Post("/upload-logo", superAdminAuth, uploadLogo)
	admin.Post("/wp-sync", superAdminAuth, syncWordPressPosts)
	admin.Post("/wp-xml-import", superAdminAuth, importWordPressXML)
	admin.Post("/import/wordpress", superAdminAuth, importWPXMLToPosts)
	admin.Post("/import/wordpress-live", superAdminAuth, importFromWPLive)
	admin.Delete("/import/rollback/:batchId", superAdminAuth, rollbackImportBatch)
	admin.Delete("/posts/soft/:id", superAdminAuth, softDeletePost)
	admin.Put("/posts/restore/:id", superAdminAuth, restorePost)
	admin.Post("/posts/sync-live", requireRole("admin", "ppid_admin"), handleManualLiveSync)
	admin.Get("/migration/check", superAdminAuth, checkMigrationIntegrity)
	admin.Post("/pddikti/preview-lecturers", superAdminAuth, previewPDDIKTIUPERTISLecturers)
	admin.Post("/pddikti/seed-lecturers", superAdminAuth, seedPDDIKTIUPERTISLecturers)
	admin.Get("/pddikti/search-lecturer", superAdminAuth, searchPDDIKTILecturer)
	admin.Post("/pddikti/import-single", superAdminAuth, importSinglePDDIKTILecturer)
	admin.Get("/whatsapp/status", superAdminAuth, getWhatsAppStatus)
	admin.Get("/whatsapp/reset", superAdminAuth, resetWhatsApp)

	// Google Scholar sync (superadmin)
	admin.Get("/scholar/status", superAdminAuth, scholarSyncStatus)
	admin.Post("/scholar/sync", superAdminAuth, scholarSyncAll)
	admin.Post("/scholar/sync/:id", superAdminAuth, scholarSyncOne)

	admin.Get("/sites", superAdminAuth, getAdminSites)
	admin.Get("/live-users", superAdminAuth, getLiveUsers)
	admin.Post("/sites", superAdminAuth, createSite)
	admin.Put("/sites/:id", requireRole("admin", "ppid_admin"), updateSite)
	admin.Delete("/sites/:id", superAdminAuth, deleteSite)

	// Both Admin and Faculty Admin
	facultyAuth := requireRole("admin", "faculty_admin", "ppid_admin")
	admin.Get("/dashboard", facultyAuth, getDashboardStats)
	admin.Get("/notifications", facultyAuth, getNotificationBadges)
	admin.Get("/pmb/dashboard", facultyAuth, getPMBDashboard)
	admin.Get("/pmb/payments", superAdminAuth, getPendingPayments)
	admin.Put("/pmb/payments/:id", superAdminAuth, verifyPayment)
	admin.Put("/pmb/candidates/:id/status", superAdminAuth, updateCandidateStatus)
	admin.Get("/pmb/export", superAdminAuth, exportCandidates)
	admin.Get("/analytics", facultyAuth, getAnalytics)
	admin.Get("/analytics/visitors", facultyAuth, getVisitorAnalytics)

	admin.Post("/uploads", facultyAuth, uploadMultipleFiles)
	admin.Post("/ai/generate", facultyAuth, handleAdminAIGenerate)
	admin.Get("/media", facultyAuth, getMediaLibrary)
	admin.Post("/media/bulk-delete", facultyAuth, bulkDeleteMedia)
	admin.Delete("/media/:id", facultyAuth, deleteMedia)
	admin.Get("/:table", facultyAuth, getRows)
	admin.Post("/:table/bulk", facultyAuth, bulkAction)
	admin.Post("/:table", facultyAuth, createRow)
	admin.Put("/:table/:id", facultyAuth, updateRow)
	admin.Delete("/:table/:id", facultyAuth, deleteRow)
}
