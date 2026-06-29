export interface Post {
  id: string;
  title: string;
  slug: string;
  excerpt: string;
  content?: string;
  cover_image_url?: string;
  featured_image?: string;
  category?: string;
  created_at: string;
}

export interface PageData {
  id: string;
  title: string;
  subtitle?: string;
  content: string;
  slug: string;
  cover_image_url?: string;
}

export interface DownloadItem {
  id: string;
  title: string;
  description: string;
  file_url: string;
  category: string;
  downloads_count: number;
}

export interface GalleryItem {
  id: string;
  title: string;
  description?: string;
  media_url: string;
  media_type: string;
  created_at: string;
}

export interface PermohonanTicket {
  id: string;
  ticket_number: string;
  applicant_type: string;
  name: string;
  identity_number: string;
  email: string;
  phone: string;
  address: string;
  details: string;
  purpose: string;
  obtain_method: string;
  delivery_method: string;
  attachment_url: string;
  status: string;
  admin_response?: string;
  created_at: string;
}

export interface UserProfile {
  id: string;
  email: string;
  full_name: string;
  role: string;
}
