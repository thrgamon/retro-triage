export type ColumnType = 'went_well' | 'didnt_go_well' | 'puzzling' | 'action_item';

export const COLUMN_TYPES: { type: ColumnType; label: string; colour: string }[] = [
	{ type: 'went_well', label: 'Went Well', colour: 'border-green-500/40' },
	{ type: 'didnt_go_well', label: "Didn't Go Well", colour: 'border-red-500/40' },
	{ type: 'puzzling', label: 'Puzzling', colour: 'border-yellow-500/40' },
	{ type: 'action_item', label: 'Action Items', colour: 'border-blue-500/40' },
];

export interface RetroListItem {
	id: string;
	name: string;
	created_at: string;
}

export interface CardResponse {
	id: string;
	column_type: ColumnType;
	content: string;
	author_name: string;
	created_at: string;
}

export interface RetroResponse {
	id: string;
	name: string;
	cards: CardResponse[];
	created_at: string;
}

export interface AnalysisGroup {
	theme: string;
	card_ids: string[];
	synthesis: string;
	five_whys: string[];
	hypothesised_causes: string[];
}

export interface AnalysisResult {
	groups: AnalysisGroup[];
	overall_summary: string;
}

export interface AnalysisResponse {
	id: string;
	retro_id: string;
	result: AnalysisResult;
	created_at: string;
}
