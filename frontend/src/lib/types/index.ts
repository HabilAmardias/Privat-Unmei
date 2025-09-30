import { type ActionResult } from '@sveltejs/kit';

export type HTTPMethod = 'GET' | 'POST' | 'PATCH' | 'DELETE';

export type Fetch = (input: string | URL, init?: RequestInit) => Promise<Response>;

export type SameSite = 'none' | 'lax' | 'strict';

export type MessageResponse = {
	message: string;
};

export type ServerResponse<T> = {
	success: boolean;
	data: T;
};

export type FetchReturn = {
	success: boolean;
	message: string;
	status: number;
	res?: Response;
};

export type CookiesData = {
	key: string;
	value: string;
	maxAge?: number;
	path: string;
	domain?: string;
	sameSite?: SameSite;
	httpOnly?: boolean;
};

export type EnhancementArgs = {
	action: URL;
	formData: FormData;
	formElement: HTMLFormElement;
	controller: AbortController;
	submitter: HTMLElement | null;
	cancel: () => void;
};

type updateOps = { reset?: boolean; invalidateAll?: boolean };

export type EnhancementReturn = {
	result: ActionResult;
	update: (opts?: updateOps) => Promise<void>;
};
