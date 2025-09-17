<script lang="ts">
    type imageProps = {
        height?: number
        width?: number
        src: string
        alt?: string
        rounded?: boolean
    }
    let {height=32, width=32, src, alt, rounded} : imageProps = $props()
    let containerClass = $state<string>(`h-${height} w-${width} `)
    if (rounded){
        containerClass += `rounded-full `
    }
    let imageClass = "h-full w-full"
</script>

<div class={containerClass}>
    <img class={imageClass} src={src} alt={alt}/>
</div>