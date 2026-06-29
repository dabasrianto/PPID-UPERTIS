package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ─── Types ────────────────────────────────────────────────────────────────────

type ScholarPaper struct {
	Title    string `json:"title"`
	Authors  string `json:"authors"`
	Journal  string `json:"journal"`
	Year     int    `json:"year"`
	Citations int   `json:"citations"`
	URL      string `json:"url"`
}

type ScholarData struct {
	ScholarID        string            `json:"scholar_id"`
	ScholarURL       string            `json:"scholar_url"`
	TotalCitations   int               `json:"total_citations"`
	HIndex           int               `json:"h_index"`
	I10Index         int               `json:"i10_index"`
	CitationsPerYear map[string]int    `json:"citations_per_year"`
	Papers           []ScholarPaper    `json:"papers"`
	FetchedAt        string            `json:"fetched_at"`
}

// ─── Sync Job State ───────────────────────────────────────────────────────────

type scholarSyncState struct {
	mu        sync.Mutex
	running   bool
	total     int
	done      int
	errors    []string
	lastStart time.Time
}

var syncState = &scholarSyncState{
	errors: []string{},
}

func getSyncStatus() map[string]interface{} {
	syncState.mu.Lock()
	defer syncState.mu.Unlock()
	progress := 0
	if syncState.total > 0 {
		progress = syncState.done * 100 / syncState.total
	}
	// Always return non-null errors array
	errs := syncState.errors
	if errs == nil {
		errs = []string{}
	}
	return map[string]interface{}{
		"running":    syncState.running,
		"total":      syncState.total,
		"done":       syncState.done,
		"progress":   progress,
		"errors":     errs,
		"last_start": syncState.lastStart.Format(time.RFC3339),
	}
}

// ─── HTTP Client with Browser Headers ────────────────────────────────────────

var scholarClient = &http.Client{
	Timeout: 20 * time.Second,
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:124.0) Gecko/20100101 Firefox/124.0",
}

func randomUA() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func fetchURL(targetURL string) (string, error) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", randomUA())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "id-ID,id;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := scholarClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return "", fmt.Errorf("rate limited by Google (429)")
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// ─── Name Normalization ───────────────────────────────────────────────────────

// toTitleCase converts "EKA FITRIANDA" → "Eka Fitrianda"
func toTitleCase(s string) string {
	words := strings.Fields(strings.ToLower(s))
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

// ─── HTML Parsers ─────────────────────────────────────────────────────────────

var (
	reCitations   = regexp.MustCompile(`(?i)cited by (\d+)`)
	reYear        = regexp.MustCompile(`\b(19|20)\d{2}\b`)
	reScholarID   = regexp.MustCompile(`user=([A-Za-z0-9_-]+)`)
	reProfileLink = regexp.MustCompile(`href="(/citations\?[^"]*view_op=list_works[^"]*)"`)
	reSearchResult = regexp.MustCompile(`(?s)<div class="gs_ri">(.*?)</div>\s*</div>\s*</div>`)
	rePaperTitle  = regexp.MustCompile(`<h3[^>]*class="gs_rt"[^>]*>(.*?)</h3>`)
	rePaperLink   = regexp.MustCompile(`href="([^"]+)"[^>]*>\s*<h3`)
	reAuthors     = regexp.MustCompile(`<div class="gs_a">(.*?)</div>`)
	reTagStrip    = regexp.MustCompile(`<[^>]+>`)
	reCiteCount   = regexp.MustCompile(`Dikutip.*?(\d+)`)
	reStatTable   = regexp.MustCompile(`(?s)<table[^>]*id="gsc_rsb_st"[^>]*>(.*?)</table>`)
	reTdVal       = regexp.MustCompile(`<td[^>]*class="[^"]*gsc_rsb_std[^"]*"[^>]*>(\d+)</td>`)
	reYearBarVal  = regexp.MustCompile(`(?s)<a[^>]*class="[^"]*gsc_g_a[^"]*"[^>]*>(\d+)</a>`)
	reYearLabel   = regexp.MustCompile(`(?s)<span[^>]*class="[^"]*gsc_g_t[^"]*"[^>]*>(\d{4})</span>`)
	rePaperRow    = regexp.MustCompile(`(?s)<tr[^>]*class="[^"]*gsc_a_tr[^"]*"[^>]*>(.*?)</tr>`)
	rePaperTitleProfile = regexp.MustCompile(`<a[^>]*class="[^"]*gsc_a_at[^"]*"[^>]*>(.*?)</a>`)
	rePaperCiteProfile  = regexp.MustCompile(`<a[^>]*class="[^"]*gsc_a_ac[^"]*"[^>]*>(\d+)</a>`)
	rePaperYearProfile  = regexp.MustCompile(`<span[^>]*class="[^"]*gsc_a_hc[^"]*"[^>]*>(\d{4})</span>`)
	rePaperHref         = regexp.MustCompile(`href="(/citations\?[^"]*view_op=view_citation[^"]*)"`)
)

func stripTags(s string) string {
	return strings.TrimSpace(reTagStrip.ReplaceAllString(s, ""))
}

// parseScholarSearchHTML extracts scholar ID from search results page
func parseScholarSearchHTML(html string) (scholarID string, profileURL string) {
	// Look for user profile link
	matches := reProfileLink.FindAllStringSubmatch(html, -1)
	for _, m := range matches {
		href := m[1]
		idMatch := reScholarID.FindStringSubmatch(href)
		if len(idMatch) > 1 {
			return idMatch[1], "https://scholar.google.com" + href
		}
	}
	return "", ""
}

// parseScholarProfile extracts full data from scholar profile page
func parseScholarProfile(html string, scholarID string) ScholarData {
	data := ScholarData{
		ScholarID:        scholarID,
		ScholarURL:       fmt.Sprintf("https://scholar.google.com/citations?user=%s&hl=id", scholarID),
		CitationsPerYear: make(map[string]int),
		FetchedAt:        time.Now().UTC().Format(time.RFC3339),
	}

	// Extract stats table (total citations, h-index, i10-index)
	statMatch := reStatTable.FindStringSubmatch(html)
	if len(statMatch) > 1 {
		vals := reTdVal.FindAllStringSubmatch(statMatch[1], -1)
		// Usually: citations_all, citations_5yr, h_index_all, h_index_5yr, i10_all, i10_5yr
		if len(vals) >= 1 {
			data.TotalCitations, _ = strconv.Atoi(vals[0][1])
		}
		if len(vals) >= 3 {
			data.HIndex, _ = strconv.Atoi(vals[2][1])
		}
		if len(vals) >= 5 {
			data.I10Index, _ = strconv.Atoi(vals[4][1])
		}
	}

	// Extract citations per year bar chart
	yearLabels := reYearLabel.FindAllStringSubmatch(html, -1)
	yearVals := reYearBarVal.FindAllStringSubmatch(html, -1)
	for i, y := range yearLabels {
		if i < len(yearVals) {
			v, _ := strconv.Atoi(yearVals[i][1])
			data.CitationsPerYear[y[1]] = v
		}
	}

	// Extract paper list from profile
	papers := []ScholarPaper{}
	paperRows := rePaperRow.FindAllStringSubmatch(html, -1)
	for _, row := range paperRows {
		rowHTML := row[1]

		titleMatch := rePaperTitleProfile.FindStringSubmatch(rowHTML)
		if len(titleMatch) < 2 {
			continue
		}
		title := stripTags(titleMatch[1])
		if title == "" {
			continue
		}

		paper := ScholarPaper{Title: title}

		// Paper URL
		hrefMatch := rePaperHref.FindStringSubmatch(rowHTML)
		if len(hrefMatch) > 1 {
			paper.URL = "https://scholar.google.com" + strings.ReplaceAll(hrefMatch[1], "&amp;", "&")
		}

		// Citations
		citeMatch := rePaperCiteProfile.FindStringSubmatch(rowHTML)
		if len(citeMatch) > 1 {
			paper.Citations, _ = strconv.Atoi(citeMatch[1])
		}

		// Year
		yearMatch := rePaperYearProfile.FindStringSubmatch(rowHTML)
		if len(yearMatch) > 1 {
			paper.Year, _ = strconv.Atoi(yearMatch[1])
		}

		papers = append(papers, paper)
	}

	// Sort by citations desc
	sort.Slice(papers, func(i, j int) bool {
		return papers[i].Citations > papers[j].Citations
	})
	if len(papers) > 20 {
		papers = papers[:20]
	}
	data.Papers = papers

	return data
}

// ─── Core Scrape Function ─────────────────────────────────────────────────────

func cleanAcademicDegrees(name string) string {
	// 1. Split by comma to remove suffix degrees
	if parts := strings.Split(name, ","); len(parts) > 0 {
		name = parts[0]
	}

	// 2. Remove prefix titles
	words := strings.Fields(name)
	cleanedWords := []string{}

	prefixSet := map[string]bool{
		"prof": true, "professor": true, "dr": true, "drg": true,
		"apt": true, "ns": true, "drs": true, "dra": true, "ir": true,
		"hj": true, "h": true,
	}

	for _, word := range words {
		// Remove trailing dot for check
		checkWord := strings.ToLower(strings.TrimSuffix(word, "."))
		if prefixSet[checkWord] {
			continue
		}
		cleanedWords = append(cleanedWords, word)
	}

	if len(cleanedWords) == 0 {
		return name
	}

	return strings.Join(cleanedWords, " ")
}

// scrapeScholarByName searches Google Scholar for a lecturer by name and affiliation
func scrapeScholarByName(name string, affiliation string) (*ScholarData, error) {
	normalName := toTitleCase(cleanAcademicDegrees(name))
	if affiliation == "" {
		affiliation = "Universitas Perintis"
	}
	query := fmt.Sprintf(`"%s" "%s"`, normalName, affiliation)
	searchURL := fmt.Sprintf(
		"https://scholar.google.com/scholar?q=%s&hl=id",
		url.QueryEscape(query),
	)

	html, err := fetchURL(searchURL)
	if err != nil {
		return nil, fmt.Errorf("search fetch error: %w", err)
	}

	// Check if blocked
	if strings.Contains(html, "unusual traffic") || strings.Contains(html, "CAPTCHA") {
		return nil, fmt.Errorf("blocked by Google CAPTCHA")
	}

	scholarID, profileURL := parseScholarSearchHTML(html)

	// If not found via search, try direct author search
	if scholarID == "" {
		searchURL2 := fmt.Sprintf(
			"https://scholar.google.com/scholar?q=author:%s+%s&hl=id",
			url.QueryEscape(normalName),
			url.QueryEscape(affiliation),
		)
		html2, err2 := fetchURL(searchURL2)
		if err2 == nil {
			scholarID, profileURL = parseScholarSearchHTML(html2)
		}
	}

	if scholarID == "" {
		return nil, fmt.Errorf("profil Google Scholar tidak ditemukan untuk: %s (afiliasi: %s)", normalName, affiliation)
	}

	// Small delay before fetching profile
	time.Sleep(time.Duration(1500+rand.Intn(1500)) * time.Millisecond)

	// Fetch full profile
	if profileURL == "" {
		profileURL = fmt.Sprintf("https://scholar.google.com/citations?user=%s&hl=id&sortby=citedby", scholarID)
	} else if !strings.Contains(profileURL, "sortby") {
		profileURL += "&sortby=citedby"
	}
	profileHTML, err := fetchURL(profileURL)
	if err != nil {
		return nil, fmt.Errorf("profile fetch error: %w", err)
	}

	result := parseScholarProfile(profileHTML, scholarID)
	return &result, nil
}

// ─── DB Sync Functions ────────────────────────────────────────────────────────

type lecturerForScholar struct {
	ID              string
	Name            string
	ScholarFetchedAt *time.Time
}

func getLecturersForScholarSync(ctx context.Context, forceAll bool) ([]lecturerForScholar, error) {
	query := `
		SELECT id, name, scholar_fetched_at
		FROM faculty_lecturers
		WHERE active = true
		ORDER BY sort_order, name
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []lecturerForScholar
	sevenDaysAgo := time.Now().Add(-7 * 24 * time.Hour)

	for rows.Next() {
		var l lecturerForScholar
		var fetchedAt *time.Time
		if err := rows.Scan(&l.ID, &l.Name, &fetchedAt); err != nil {
			continue
		}
		l.ScholarFetchedAt = fetchedAt

		// Skip if fetched recently and not force mode
		if !forceAll && fetchedAt != nil && fetchedAt.After(sevenDaysAgo) {
			continue
		}
		result = append(result, l)
	}
	return result, nil
}

func saveScholarData(ctx context.Context, lecturerID string, data *ScholarData) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	now := time.Now()
	_, err = db.ExecContext(ctx,
		`UPDATE faculty_lecturers SET scholar_data = $1, scholar_fetched_at = $2, scholar_id = $3 WHERE id = $4`,
		string(jsonBytes), now, data.ScholarID, lecturerID,
	)
	return err
}

// ─── Batch Sync Job ───────────────────────────────────────────────────────────

func batchSyncAllScholars(ctx context.Context, forceAll bool, affiliation string) {
	syncState.mu.Lock()
	if syncState.running {
		syncState.mu.Unlock()
		return
	}
	syncState.running = true
	syncState.errors = []string{}
	syncState.done = 0
	syncState.lastStart = time.Now()
	syncState.mu.Unlock()

	defer func() {
		syncState.mu.Lock()
		syncState.running = false
		syncState.mu.Unlock()
	}()

	lecturers, err := getLecturersForScholarSync(ctx, forceAll)
	if err != nil {
		syncState.mu.Lock()
		syncState.errors = append(syncState.errors, "DB error: "+err.Error())
		syncState.mu.Unlock()
		return
	}

	syncState.mu.Lock()
	syncState.total = len(lecturers)
	syncState.mu.Unlock()

	log.Printf("[Scholar] Starting batch sync: %d lecturers, affiliation=%q", len(lecturers), affiliation)

	for _, l := range lecturers {
		select {
		case <-ctx.Done():
			return
		default:
		}

		data, err := scrapeScholarByName(l.Name, affiliation)
		if err != nil {
			log.Printf("[Scholar] %s: %v", l.Name, err)
			syncState.mu.Lock()
			syncState.errors = append(syncState.errors, fmt.Sprintf("%s: %v", l.Name, err))
			if len(syncState.errors) > 50 {
				syncState.errors = syncState.errors[len(syncState.errors)-50:]
			}
			syncState.done++
			syncState.mu.Unlock()
		} else {
			if err2 := saveScholarData(ctx, l.ID, data); err2 != nil {
				log.Printf("[Scholar] Save error for %s: %v", l.Name, err2)
			} else {
				log.Printf("[Scholar] ✓ %s → h=%d citations=%d papers=%d",
					l.Name, data.HIndex, data.TotalCitations, len(data.Papers))
			}
			syncState.mu.Lock()
			syncState.done++
			syncState.mu.Unlock()
		}

		// Random delay 3-8 seconds between requests
		delay := time.Duration(3000+rand.Intn(5000)) * time.Millisecond
		time.Sleep(delay)
	}

	log.Printf("[Scholar] Batch sync complete: %d processed, %d errors", len(lecturers), len(syncState.errors))
}

// ─── HTTP Handlers ────────────────────────────────────────────────────────────

// GET /api/v1/admin/scholar/status
func scholarSyncStatus(c *fiber.Ctx) error {
	return c.JSON(getSyncStatus())
}

// POST /api/v1/admin/scholar/sync
// Body: { "force": true/false, "affiliation": "Universitas Perintis Indonesia" }
func scholarSyncAll(c *fiber.Ctx) error {
	syncState.mu.Lock()
	running := syncState.running
	syncState.mu.Unlock()
	if running {
		return c.Status(409).JSON(fiber.Map{"error": "Sync sedang berjalan"})
	}

	var body struct {
		Force       bool   `json:"force"`
		Affiliation string `json:"affiliation"`
	}
	_ = c.BodyParser(&body)
	if body.Affiliation == "" {
		body.Affiliation = "Universitas Perintis"
	}

	// Run in background goroutine
	go batchSyncAllScholars(context.Background(), body.Force, body.Affiliation)

	return c.JSON(fiber.Map{"message": "Scholar sync dimulai di background", "force": body.Force, "affiliation": body.Affiliation})
}

// POST /api/v1/admin/scholar/sync/:id
func scholarSyncOne(c *fiber.Ctx) error {
	id := c.Params("id")
	var body struct {
		Affiliation string `json:"affiliation"`
	}
	_ = c.BodyParser(&body)
	if body.Affiliation == "" {
		body.Affiliation = "Universitas Perintis"
	}

	var name string
	err := db.QueryRowContext(context.Background(),
		"SELECT name FROM faculty_lecturers WHERE id = $1", id).Scan(&name)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dosen tidak ditemukan"})
	}

	data, err := scrapeScholarByName(name, body.Affiliation)
	if err != nil {
		return c.Status(422).JSON(fiber.Map{"error": err.Error()})
	}

	if err2 := saveScholarData(context.Background(), id, data); err2 != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal simpan: " + err2.Error()})
	}

	return c.JSON(fiber.Map{"message": "Scholar data berhasil diperbarui", "data": data})
}

// GET /api/v1/admin/scholar/sync — alias for status
func scholarSyncGet(c *fiber.Ctx) error {
	return scholarSyncStatus(c)
}

// ─── Unused var suppress ─────────────────────────────────────────────────────
var _ = reCitations
var _ = reYear
var _ = reSearchResult
var _ = rePaperTitle
var _ = rePaperLink
var _ = reAuthors
var _ = reCiteCount
