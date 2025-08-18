import { PUBLIC_API_URL } from '$env/static/public';
import { createClient, type Config } from './api/client';
import { client } from './api/client.gen';

const config: Config = {
	baseUrl: PUBLIC_API_URL,
	throwOnError: true,
	credentials: 'include',
	querySerializer: {
		array: {
			explode: false,
			style: 'form'
		}
	}
};

client.setConfig(config);

export function getClient(
	fetch:
		| {
				(input: RequestInfo | URL, init?: RequestInit): Promise<Response>;
				(input: string | URL | globalThis.Request, init?: RequestInit): Promise<Response>;
		  }
		| undefined
) {
	if (typeof fetch === 'undefined') {
		return client;
	}

	return createClient({
		...config,
		fetch
	});
}

export default client;
