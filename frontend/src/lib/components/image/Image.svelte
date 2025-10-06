<script lang="ts">
	type imageProps = {
		height?: number;
		width?: number;
		src: string;
		alt?: string;
		additionalClass?: string;
		round?: 'sm' | 'full' | 'lg';
	};
	let { height = 32, width = 32, src, alt, round, additionalClass }: imageProps = $props();
	let containerStyle = `height: ${height}px; width: ${width}px;`;
	let imageClass = $state<string>('h-full w-full object-cover');
	let containerClass = $state<string>('');
	if (round) {
		containerClass = `rounded-${round} overflow-hidden`;
	}
	if (additionalClass) {
		imageClass += ` ${additionalClass}`;
	}
	let isLoading = $state<boolean>(true);
</script>

<div class={containerClass} style={containerStyle}>
	<img
		style:display={isLoading ? 'none' : 'inline'}
		class={imageClass}
		{src}
		{alt}
		{height}
		{width}
		onload={() => {
			isLoading = false;
		}}
	/>
	{#if isLoading}
		<div class="h-full w-full animate-pulse object-cover"></div>
	{/if}
</div>
