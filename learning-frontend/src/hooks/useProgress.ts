import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type { TopicStats } from "@/types";

export function useTopicProgress(topicId: string) {
  return useQuery({
    queryKey: ["progress", topicId],
    queryFn: () =>
      api.get<{ data: Record<string, string> }>(`/topics/${topicId}/progress`).then((r) => r.data),
    enabled: !!topicId,
  });
}

export function useUpsertProgress() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ itemId, topicId, status }: { itemId: string; topicId: string; status: string }) =>
      api.put(`/progress/${itemId}`, { topic_id: topicId, status }),
    onSuccess: (_data, vars) => {
      qc.invalidateQueries({ queryKey: ["progress", vars.topicId] });
      qc.invalidateQueries({ queryKey: ["stats"] });
    },
  });
}

export function useResetProgress(topicId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: () => api.delete(`/topics/${topicId}/progress`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["progress", topicId] });
      qc.invalidateQueries({ queryKey: ["stats"] });
    },
  });
}

export function useStats() {
  return useQuery({
    queryKey: ["stats"],
    queryFn: () =>
      api.get<{ data: TopicStats[] }>("/stats").then((r) => r.data),
  });
}
