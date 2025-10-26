<script lang="ts">
	import type { Snippet } from 'svelte';

	type DownloadLinkProps = {
		children: Snippet;
		href: string;
		className?: string;
		download: string;
	};
	let { children, href, className, download }: DownloadLinkProps = $props();
</script>

<a class={className} {href} {download}>{@render children()}</a>
