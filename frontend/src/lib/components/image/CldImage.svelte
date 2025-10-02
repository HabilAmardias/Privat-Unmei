<script lang="ts">
	import { CldImage } from 'svelte-cloudinary';
	type cldImageProps = {
		src: string;
		alt?: string;
		width?: number;
		height?: number;
		className?: string;
	};
	let { src, alt, width = 32, height = 32, className }: cldImageProps = $props();
</script>

<CldImage {src} {alt} {width} {height} class={className} />
