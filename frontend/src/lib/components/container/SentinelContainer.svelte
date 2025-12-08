<script lang="ts">
	import { onMount } from 'svelte';
	let { onIntersect }: { onIntersect: () => void } = $props();

	let container = $state<HTMLDivElement>();

	onMount(() => {
		const observer = new IntersectionObserver(
			(entries) => {
				entries.forEach((el) => {
					if (el.isIntersecting) {
						onIntersect();
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

<div bind:this={container} class="h-[1px]"></div>
