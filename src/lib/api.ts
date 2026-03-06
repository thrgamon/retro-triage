import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { AnalysisResponse, RetroListItem, RetroResponse, CardResponse, ColumnType } from './types';

const BASE = '/api';

async function fetchJSON<T>(url: string, init?: RequestInit): Promise<T> {
	const res = await fetch(`${BASE}${url}`, {
		...init,
		headers: { 'Content-Type': 'application/json', ...init?.headers },
	});

	if (res.status === 204) return undefined as T;

	const body = await res.json().catch(() => null);
	if (!res.ok) {
		const msg = body?.error ?? `Request failed (${res.status})`;
		throw new Error(msg);
	}
	return body as T;
}

// --- Retros ---

export function useRetros() {
	return useQuery({
		queryKey: ['retros'],
		queryFn: () => fetchJSON<RetroListItem[]>('/retros'),
	});
}

export function useRetro(id: string) {
	return useQuery({
		queryKey: ['retro', id],
		queryFn: () => fetchJSON<RetroResponse>(`/retros/${id}`),
		enabled: !!id,
	});
}

export function useCreateRetro() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: (name: string) => fetchJSON<RetroListItem>('/retros', { method: 'POST', body: JSON.stringify({ name }) }),
		onSuccess: () => qc.invalidateQueries({ queryKey: ['retros'] }),
	});
}

// --- Cards ---

export function useCreateCard(retroId: string) {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: (data: { column_type: ColumnType; content: string; author_name?: string }) =>
			fetchJSON<CardResponse>(`/retros/${retroId}/cards`, { method: 'POST', body: JSON.stringify(data) }),
		onSuccess: () => qc.invalidateQueries({ queryKey: ['retro', retroId] }),
	});
}

export function useDeleteCard(retroId: string) {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: (cardId: string) => fetchJSON<void>(`/retros/${retroId}/cards/${cardId}`, { method: 'DELETE' }),
		onSuccess: () => qc.invalidateQueries({ queryKey: ['retro', retroId] }),
	});
}

// --- Analysis ---

export function useAnalysis(retroId: string) {
	return useQuery({
		queryKey: ['analysis', retroId],
		queryFn: () => fetchJSON<AnalysisResponse>(`/retros/${retroId}/analysis`),
		enabled: !!retroId,
		retry: false,
	});
}

export function useRunAnalysis(retroId: string) {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: () => fetchJSON<AnalysisResponse>(`/retros/${retroId}/analyse`, { method: 'POST' }),
		onSuccess: () => qc.invalidateQueries({ queryKey: ['analysis', retroId] }),
	});
}
