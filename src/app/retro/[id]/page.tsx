'use client';

import { use, useState } from 'react';
import Link from 'next/link';
import { useRetro, useCreateCard, useDeleteCard, useRunAnalysis, useAnalysis } from '@/lib/api';
import { COLUMN_TYPES } from '@/lib/types';
import type { ColumnType, CardResponse } from '@/lib/types';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { Card, CardContent } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { X, Loader2 } from 'lucide-react';
import { AnalysisView } from './analysis-view';

export default function RetroPage({ params }: { params: Promise<{ id: string }> }) {
	const { id } = use(params);
	const { data: retro, isLoading } = useRetro(id);
	const { data: analysis } = useAnalysis(id);
	const runAnalysis = useRunAnalysis(id);
	const [activeTab, setActiveTab] = useState('board');

	const analysisData = runAnalysis.data ?? analysis;
	const hasAnalysis = !!analysisData;

	if (isLoading) {
		return <p className="p-8 text-muted-foreground">Loading...</p>;
	}

	if (!retro) {
		return <p className="p-8 text-muted-foreground">Retro not found.</p>;
	}

	return (
		<main className="min-h-screen px-4 py-6">
			<div className="mb-6 flex items-center justify-between">
				<div>
					<Link href="/" className="text-sm text-muted-foreground hover:underline">
						&larr; All retros
					</Link>
					<h1 className="text-2xl font-bold">{retro.name}</h1>
				</div>
				<Button
					onClick={async () => {
						await runAnalysis.mutateAsync();
						setActiveTab('analysis');
					}}
					disabled={runAnalysis.isPending || retro.cards.length === 0}
					size="lg"
				>
					{runAnalysis.isPending ? (
						<>
							<Loader2 className="mr-2 h-4 w-4 animate-spin" />
							Analysing...
						</>
					) : (
						'Analyse'
					)}
				</Button>
			</div>

			{runAnalysis.isError && (
				<p className="mb-4 text-sm text-destructive-foreground">Analysis failed: {runAnalysis.error.message}</p>
			)}

			<Tabs value={activeTab} onValueChange={setActiveTab}>
				<TabsList>
					<TabsTrigger value="board">Board</TabsTrigger>
					<TabsTrigger value="analysis" disabled={!hasAnalysis}>
						Analysis
					</TabsTrigger>
				</TabsList>

				<TabsContent value="board" className="mt-4">
					<div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
						{COLUMN_TYPES.map((col) => (
							<RetroColumn
								key={col.type}
								retroId={id}
								columnType={col.type}
								label={col.label}
								colour={col.colour}
								cards={retro.cards.filter((c) => c.column_type === col.type)}
							/>
						))}
					</div>
				</TabsContent>

				<TabsContent value="analysis" className="mt-4">
					{hasAnalysis && <AnalysisView analysis={analysisData} cards={retro.cards} />}
				</TabsContent>
			</Tabs>
		</main>
	);
}

function RetroColumn({
	retroId,
	columnType,
	label,
	colour,
	cards,
}: {
	retroId: string;
	columnType: ColumnType;
	label: string;
	colour: string;
	cards: CardResponse[];
}) {
	const [content, setContent] = useState('');
	const createCard = useCreateCard(retroId);
	const deleteCard = useDeleteCard(retroId);

	const handleAdd = async () => {
		if (!content.trim()) return;
		await createCard.mutateAsync({ column_type: columnType, content: content.trim() });
		setContent('');
	};

	return (
		<div className={`rounded-lg border-2 ${colour} bg-card p-3`}>
			<h2 className="mb-3 text-sm font-semibold">{label}</h2>
			<div className="mb-3 space-y-2">
				{cards.map((card) => (
					<Card key={card.id} className="group relative">
						<CardContent className="p-2 text-sm">
							{card.content}
							<button
								type="button"
								onClick={() => deleteCard.mutate(card.id)}
								className="absolute right-1 top-1 hidden text-muted-foreground hover:text-destructive-foreground group-hover:block"
							>
								<X className="h-3 w-3" />
							</button>
						</CardContent>
					</Card>
				))}
			</div>
			<div className="flex gap-1">
				<Textarea
					placeholder="Add a card..."
					value={content}
					onChange={(e) => setContent(e.target.value)}
					onKeyDown={(e) => {
						if (e.key === 'Enter' && !e.shiftKey) {
							e.preventDefault();
							handleAdd();
						}
					}}
					className="min-h-[2.5rem] resize-none text-sm"
					rows={1}
				/>
				<Button size="sm" variant="ghost" onClick={handleAdd} disabled={createCard.isPending || !content.trim()}>
					+
				</Button>
			</div>
		</div>
	);
}
