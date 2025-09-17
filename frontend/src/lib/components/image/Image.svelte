<script lang="ts">
    type imageProps = {
        height?: number
        width?: number
        src: string
        alt?: string
        rounded?: boolean
    }
    let {height=32, width=32, src, alt, rounded} : imageProps = $props()
    let containerStyle = `height: ${height}px; width: ${width}px;`
    let containerClass = $state<string>(`inline-block overflow-hidden flex-none`)
    if (rounded){
        containerClass += ` rounded-full `
    }
    let imageClass = "h-full w-full object-cover"
</script>

<div class={containerClass} style={containerStyle}>
    <img width={width} height={height} class={imageClass} src={src} alt={alt}/>
</div>