
/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />
/// <reference types="@sveltejs/kit" />
/// <reference types="../.svelte-kit/ambient.d.ts" />

import { files, version } from '$service-worker';

const self = globalThis.self as unknown as ServiceWorkerGlobalScope;
const CACHE = `cache-${version}`;
const ASSETS = [
	...files
];

self.addEventListener('install', (event) => {
	event.waitUntil((async () => {
		const cache = await caches.open(CACHE);
		await cache.add(new Request("offline.html", { cache: 'reload' }));
		await cache.addAll(files)
	})());
});

self.addEventListener('activate', (event) => {
	event.waitUntil((async () => {
		for (const key of await caches.keys()) {
			if (key !== CACHE) await caches.delete(key);
		}
	})());
	self.clients.claim();
});

self.addEventListener('fetch', (event) => {

	if (event.request.method !== 'GET') return;

	async function respond() {
		const url = new URL(event.request.url);
		const cache = await caches.open(CACHE);

		if (ASSETS.includes(url.pathname)) {
			const response = await cache.match(url.pathname);

			if (response) {
				return response;
			}
		}

		try {
			const response = await fetch(event.request);
			if (!(response instanceof Response)) {
				throw new Error('invalid response from fetch');
			}

			return response;
		} catch (error) {
			if (ASSETS.includes(url.pathname)) {
				const response = await cache.match(url.pathname);

				if (response) {
					return response;
				}
			}
			if (url.pathname.startsWith("/app") && url.search) {
				const response = await fetch(url, { redirect: 'follow' });

				if (response) {
					return response
				}
				throw error;
			}
			console.log(error);
			const cachedResponse = await cache.match("offline.html");
			if (cachedResponse) {
				return cachedResponse;
			}
			throw error;

		}
	}

	event.respondWith(respond());
});