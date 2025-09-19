<script lang="ts">
    import { CldImage } from 'svelte-cloudinary';
    type cldImageProps = {
        src: string
        alt?: string
        width?: number
        height?: number
        round? :"sm" | "lg" | "full"
    }
    let {src, alt, width = 32, height = 32, round} : cldImageProps = $props()
    let className = $state<string>("")
    if (round){
        className = `rounded-${round}`
    }
</script>

<CldImage
    src={src}
    alt={alt}
    width={width}
    height={height}
    class={className}
/>