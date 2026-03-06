'use client';

import type { AnalysisResponse, CardResponse } from '@/lib/types';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export function AnalysisView({ analysis, cards }: { analysis: AnalysisResponse; cards: CardResponse[] }) {
	const result = analysis.result;
	const cardMap = new Map(cards.map((c) => [c.id, c]));

	return (
		<div>
			<Card className="mb-6">
				<CardHeader>
					<CardTitle className="text-base">Summary</CardTitle>
				</CardHeader>
				<CardContent>
					<p className="text-sm">{result.overall_summary}</p>
				</CardContent>
			</Card>

			<div className="space-y-6">
				{result.groups.map((group, i) => (
					<Card key={`group-${i}`}>
						<CardHeader>
							<CardTitle className="text-base">{group.theme}</CardTitle>
						</CardHeader>
						<CardContent className="space-y-4">
							<div>
								<h4 className="mb-1 text-xs font-semibold uppercase text-muted-foreground">Cards in this group</h4>
								<ul className="list-inside list-disc space-y-1 text-sm">
									{group.card_ids.map((cid) => {
										const card = cardMap.get(cid);
										return (
											<li key={cid} className="text-muted-foreground">
												{card ? card.content : cid}
											</li>
										);
									})}
								</ul>
							</div>

							<div>
								<h4 className="mb-1 text-xs font-semibold uppercase text-muted-foreground">Synthesis</h4>
								<p className="text-sm">{group.synthesis}</p>
							</div>

							<div>
								<h4 className="mb-1 text-xs font-semibold uppercase text-muted-foreground">5 Whys</h4>
								<div className="space-y-2 text-sm">
									{group.five_whys.map((why, j) => {
										const isQuestion = why.toLowerCase().startsWith('why');
										return (
											<div key={`why-${j}`} style={{ marginLeft: `${Math.floor(j / 2) * 16}px` }}>
												{isQuestion ? (
													<p className="font-medium text-foreground">{why}</p>
												) : (
													<p className="text-muted-foreground">{why}</p>
												)}
											</div>
										);
									})}
								</div>
							</div>

							<div>
								<h4 className="mb-1 text-xs font-semibold uppercase text-muted-foreground">Hypothesised Causes</h4>
								<ul className="list-inside list-disc space-y-1 text-sm">
									{group.hypothesised_causes.map((cause, j) => (
										<li key={`cause-${j}`}>{cause}</li>
									))}
								</ul>
							</div>
						</CardContent>
					</Card>
				))}
			</div>
		</div>
	);
}
