<script lang="ts">
	import { onMount, type Snippet } from 'svelte';
	let { children }: { children: Snippet } = $props();

	let isVisible = $state<boolean>(false);
	let container = $state<Element | null>(null);

	onMount(() => {
		const observer = new IntersectionObserver(
			(entries) => {
				entries.forEach((el) => {
					if (el.isIntersecting) {
						isVisible = true;
					}
				});
			},
			{ threshold: 0.5 }
		);
		if (container) {
			observer.observe(container);
		}
		return () => observer.disconnect();
	});
</script>

<div bind:this={container} class={`content ${isVisible ? 'inview' : ''}`}>
	{@render children()}
</div>

<style>
	.content {
		opacity: 0;
		transform: translateY(20px);
		transition:
			opacity 0.6s ease-out,
			transform 0.6s ease-out;
	}
	.content.inview {
		opacity: 1;
		transform: translateY(0);
	}
</style>
