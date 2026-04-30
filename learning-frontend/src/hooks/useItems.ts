import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type { Item, PaginationMeta } from "@/types";

export function useItems(topicId: string, opts?: { category?: string; page?: number; limit?: number }) {
  const params = new URLSearchParams();
  if (opts?.category) params.set("category", opts.category);
  if (opts?.page) params.set("page", String(opts.page));
  if (opts?.limit) params.set("limit", String(opts.limit));
  const qs = params.toString() ? `?${params}` : "";

  return useQuery({
    queryKey: ["items", topicId, opts],
    queryFn: () =>
      api.get<{ data: Item[]; meta: PaginationMeta }>(`/topics/${topicId}/items${qs}`),
    enabled: !!topicId,
  });
}

export function useAllItems(topicId: string) {
  return useQuery({
    queryKey: ["items", topicId, "all"],
    queryFn: () =>
      api.get<{ data: Item[]; meta: PaginationMeta }>(`/topics/${topicId}/items?limit=500`).then((r) => r.data),
    enabled: !!topicId,
  });
}

export function useCreateItem(topicId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: { term: string; meaning: string; example: string; hint?: string; category: string; sort_order?: number }) =>
      api.post<Item>(`/topics/${topicId}/items`, body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["items", topicId] }),
  });
}

export function useUpdateItem(topicId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ itemId, body }: { itemId: string; body: Partial<{ term: string; meaning: string; example: string; hint: string; category: string; sort_order: number }> }) =>
      api.put<Item>(`/topics/${topicId}/items/${itemId}`, body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["items", topicId] }),
  });
}

export function useDeleteItem(topicId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (itemId: string) => api.delete(`/topics/${topicId}/items/${itemId}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["items", topicId] }),
  });
}
