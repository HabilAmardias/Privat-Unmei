import type { ActionResult } from '@sveltejs/kit';

export type SameSiteType = 'none' | 'lax' | 'strict';

export type LoginResponse = {
	success: boolean;
	data: {
		status: 'verified' | 'unverified';
	};
};

export type CookiesData = {
	key: string;
	value: string;
	maxAge?: number;
	path: string;
	domain?: string;
	sameSite?: SameSiteType;
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
