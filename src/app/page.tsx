'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useRetros, useCreateRetro } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export default function Home() {
	const router = useRouter();
	const [name, setName] = useState('');
	const { data: retros, isLoading } = useRetros();
	const createRetro = useCreateRetro();

	const handleCreate = async () => {
		if (!name.trim()) return;
		const retro = await createRetro.mutateAsync(name.trim());
		router.push(`/retro/${retro.id}`);
	};

	return (
		<main className="mx-auto max-w-2xl px-4 py-12">
			<h1 className="mb-8 text-3xl font-bold">Retro Triage</h1>

			<div className="mb-8 flex gap-2">
				<Input
					placeholder="Retro name (e.g. Sprint 42)"
					value={name}
					onChange={(e) => setName(e.target.value)}
					onKeyDown={(e) => e.key === 'Enter' && handleCreate()}
				/>
				<Button onClick={handleCreate} disabled={createRetro.isPending || !name.trim()}>
					{createRetro.isPending ? 'Creating...' : 'New Retro'}
				</Button>
			</div>

			{isLoading && <p className="text-muted-foreground">Loading...</p>}

			{retros && retros.length > 0 && (
				<div className="space-y-2">
					<h2 className="text-lg font-semibold">Recent Retros</h2>
					{retros.map((r) => (
						<Card
							key={r.id}
							className="cursor-pointer transition-colors hover:bg-accent"
							onClick={() => router.push(`/retro/${r.id}`)}
						>
							<CardHeader className="py-3">
								<CardTitle className="text-base">{r.name}</CardTitle>
							</CardHeader>
							<CardContent className="pb-3 pt-0">
								<p className="text-xs text-muted-foreground">{new Date(r.created_at).toLocaleDateString()}</p>
							</CardContent>
						</Card>
					))}
				</div>
			)}
		</main>
	);
}
