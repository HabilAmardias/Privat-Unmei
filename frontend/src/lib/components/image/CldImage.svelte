<script lang="ts">
    import { CldImage } from 'svelte-cloudinary';
    type cldImageProps = {
        src: string
        alt?: string
        width?: number
        height?: number
    }
    let {src, alt, width = 32, height = 32} : cldImageProps = $props()
</script>

<CldImage
    src={src}
    alt={alt}
    width={width}
    height={height}
/>