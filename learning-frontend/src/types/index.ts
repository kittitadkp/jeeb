export interface Topic {
  id: string;
  name: string;
  description: string;
  category: string;
  icon: string;
  created_at: string;
  updated_at: string;
}

export interface Item {
  id: string;
  topic_id: string;
  term: string;
  meaning: string;
  example: string;
  hint: string;
  category: string;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export interface UserProgress {
  id: string;
  user_id: string;
  topic_id: string;
  item_id: string;
  status: "learning" | "mastered";
  review_count: number;
  last_reviewed_at: string;
  created_at: string;
  updated_at: string;
}

export interface TopicStats {
  topic_id: string;
  name: string;
  icon: string;
  total: number;
  mastered: number;
  learning: number;
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}
