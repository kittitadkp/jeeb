import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type { Topic } from "@/types";

export function useTopics() {
  return useQuery({
    queryKey: ["topics"],
    queryFn: () => api.get<{ data: Topic[] }>("/topics").then((r) => r.data),
  });
}

export function useTopic(id: string) {
  return useQuery({
    queryKey: ["topics", id],
    queryFn: () => api.get<Topic>(`/topics/${id}`),
    enabled: !!id,
  });
}

export function useCreateTopic() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: { name: string; description: string; category: string; icon: string }) =>
      api.post<Topic>("/topics", body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["topics"] }),
  });
}

export function useUpdateTopic(id: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: Partial<{ name: string; description: string; category: string; icon: string }>) =>
      api.put<Topic>(`/topics/${id}`, body),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["topics"] });
      qc.invalidateQueries({ queryKey: ["topics", id] });
    },
  });
}

export function useDeleteTopic() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`/topics/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["topics"] }),
  });
}
