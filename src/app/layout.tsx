import type { Metadata } from 'next';
import './globals.css';
import { QueryProvider } from '@/lib/query-provider';

export const metadata: Metadata = {
	title: 'Retro Triage',
	description: 'AI-powered retrospective analysis',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
	return (
		<html lang="en" className="dark">
			<body className="min-h-screen bg-background font-sans antialiased">
				<QueryProvider>{children}</QueryProvider>
			</body>
		</html>
	);
}
