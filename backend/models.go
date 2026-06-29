package main

import "encoding/json"

type User struct {
	ID         string   `json:"id"`
	Email      string   `json:"email"`
	FullName   string   `json:"full_name"`
	Role       string   `json:"role"`
	Status     string   `json:"status"`
	CreatedAt  string   `json:"created_at"`
	FacultyIDs []string `json:"faculty_ids,omitempty"`
}

type HeroSlide struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Eyebrow     string `json:"eyebrow"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	CTA1Label   string `json:"cta_primary_label"`
	CTA1Href    string `json:"cta_primary_href"`
	CTA2Label   string `json:"cta_secondary_label"`
	CTA2Href    string `json:"cta_secondary_href"`
	VideoURL    string `json:"video_url"`
	Active      bool   `json:"active"`
	SortOrder   int    `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type Faculty struct {
	ID              string      `json:"id"`
	Code            string      `json:"code"`
	Name            string      `json:"name"`
	Slug            string      `json:"slug"`
	Description     string      `json:"description"`
	Vision          string      `json:"vision"`
	Mission         string      `json:"mission"`
	AboutContent    string      `json:"about_content"`
	HeroTitle       string      `json:"hero_title"`
	HeroEyebrow     string      `json:"hero_eyebrow"`
	HeroDescription string      `json:"hero_description"`
	CoverImageURL   string      `json:"cover_image_url"`
	CoverImages     *json.RawMessage `json:"cover_images"`
	AboutImageURL   string      `json:"about_image_url"`
	AboutImages     *json.RawMessage `json:"about_images"`
	FacilityImages  *json.RawMessage `json:"facility_images"`
	Accent          string      `json:"accent"`
	Active          bool        `json:"active"`
	Programs        int         `json:"programs"`
	SortOrder       int         `json:"sort_order"`
	Facilities      *json.RawMessage `json:"facilities"`
	ContactInfo     *json.RawMessage `json:"contact_info"`
	Tujuan          string           `json:"tujuan"`
	StrukturOrganisasiURL string `json:"struktur_organisasi_url"`
	Kerjasama       *json.RawMessage `json:"kerjasama"`
	CreatedAt       string      `json:"created_at"`
	UpdatedAt       string      `json:"updated_at"`
}

type Program struct {
	ID            string `json:"id"`
	FacultyID     string `json:"faculty_id"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Level         string `json:"level"`
	Description   string `json:"description"`
	Accreditation string `json:"accreditation"`
	DurationYears int    `json:"duration_years"`
	Active        bool             `json:"active"`
	SortOrder     int              `json:"sort_order"`
	Syllabus      *json.RawMessage `json:"syllabus,omitempty"`
	CareerPaths   *json.RawMessage `json:"career_paths,omitempty"`
	StrukturOrganisasi *json.RawMessage `json:"struktur_organisasi,omitempty"`
	Visi                  string           `json:"visi"`
	Misi                  string           `json:"misi"`
	Tujuan                string           `json:"tujuan"`
	GelarAkademik         string           `json:"gelar_akademik"`
	KompetensiLulusan     string           `json:"kompetensi_lulusan"`
	FasilitasLaboratorium string           `json:"fasilitas_laboratorium"`
	FasilitasLaboratoriumImage string           `json:"fasilitas_laboratorium_image"`
	FasilitasLaboratoriumImages *json.RawMessage `json:"fasilitas_laboratorium_images,omitempty"`
	CreatedAt     string           `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type Lecturer struct {
	ID              string           `json:"id"`
	FacultyID       string           `json:"faculty_id"`
	Slug            string           `json:"slug"`
	Name            string           `json:"name"`
	Gelar           string           `json:"gelar"`
	Position        string           `json:"position"`
	Education       string           `json:"education"`
	Expertise       string           `json:"expertise"`
	PhotoURL        string           `json:"photo_url"`
	LeadershipGroup string           `json:"leadership_group"`
	Active          bool             `json:"active"`
	SortOrder       int              `json:"sort_order"`
	Source          string           `json:"source"`
	PDDIKTIID       string           `json:"pddikti_id"`
	PDDIKTIUUID     string           `json:"pddikti_uuid,omitempty"`
	ScholarID       string           `json:"scholar_id,omitempty"`
	ScholarData     *json.RawMessage `json:"scholar_data,omitempty"`
	CreatedAt       string           `json:"created_at"`
	UpdatedAt       string           `json:"updated_at"`
}


type News struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Slug           string `json:"slug"`
	Excerpt        string `json:"excerpt"`
	Content        string `json:"content"`
	ImageURL       string `json:"image_url"`
	Category       string `json:"category"`
	Date           string `json:"date"`
	Featured       bool   `json:"featured"`
	Active         bool   `json:"active"`
	AuthorName     string `json:"author_name"`
	SeoTitle       string `json:"seo_title"`
	SeoDescription string `json:"seo_description"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type Testimonial struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	Quote        string `json:"quote"`
	ImageURL     string `json:"image_url"`
	Year         string `json:"year"`
	VideoURL     string `json:"video_url"`
	YoutubeURL   string `json:"youtube_url"`
	TiktokURL    string `json:"tiktok_url"`
	InstagramURL string `json:"instagram_url"`
	Active       bool   `json:"active"`
	SortOrder    int    `json:"sort_order"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type BlogPost struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Slug           string   `json:"slug"`
	Excerpt        string   `json:"excerpt"`
	Content        string   `json:"content"`
	CoverImageURL  string   `json:"cover_image_url"`
	AuthorName     string   `json:"author_name"`
	Category       string   `json:"category"`
	Tags           []string `json:"tags"`
	FacultyID      string   `json:"faculty_id"`
	Status         string   `json:"status"`
	PublishedAt    string   `json:"published_at"`
	SeoTitle       string   `json:"seo_title"`
	SeoDescription string   `json:"seo_description"`
	Views          int      `json:"views"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

type Page struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Subtitle       string `json:"subtitle"`
	Slug           string `json:"slug"`
	Content        string `json:"content"`
	CoverImageURL  string `json:"cover_image_url"`
	Published      bool   `json:"published"`
	SortOrder      int    `json:"sort_order"`
	SeoTitle       string `json:"seo_title"`
	SeoDescription string `json:"seo_description"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type ContactMessageReply struct {
	ID         string `json:"id"`
	MessageID  string `json:"message_id"`
	SenderType string `json:"sender_type"` // "admin" or "user"
	SenderName string `json:"sender_name"`
	Message    string `json:"message"`
	CreatedAt  string `json:"created_at"`
}

type ContactMessage struct {
	ID           string                `json:"id"`
	Name         string                `json:"name"`
	Email        string                `json:"email"`
	Phone        string                `json:"phone"`
	Subject      string                `json:"subject"`
	Message      string                `json:"message"`
	IsRead       bool                  `json:"is_read"`
	IsReadUser   bool                  `json:"is_read_user"`
	ReplyMessage string                `json:"reply_message"`
	RepliedAt    string                `json:"replied_at"`
	Replies      []ContactMessageReply `json:"replies"`
	CreatedAt    string                `json:"created_at"`
}

type SiteSetting struct {
	Key       string          `json:"key"`
	Value     *json.RawMessage `json:"value"`
	UpdatedAt string          `json:"updated_at"`
}

type CampusEvent struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	EventDate   string `json:"event_date"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Location       string `json:"location"`
	MapCoordinates string `json:"map_coordinates"`
	Active         bool   `json:"active"`
	SortOrder      int    `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
