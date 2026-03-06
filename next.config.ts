import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
	output: 'standalone',
	experimental: {
		proxyTimeout: 120_000,
	},
	async rewrites() {
		return [
			{
				source: '/api/:path*',
				destination: `${process.env.API_URL || 'http://localhost:8080'}/api/:path*`,
			},
		];
	},
};

export default nextConfig;
