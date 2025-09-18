<script lang="ts">
    type imageProps = {
        height?: number
        width?: number
        src: string
        alt?: string
        className: string
    }
    let {height=32, width=32, src, alt, className} : imageProps = $props()
    let containerStyle = `height: ${height}px; width: ${width}px;`
    let isLoading = $state<boolean>(true)
</script>

<div class={className} style={containerStyle}>
    <img style:display={isLoading? "none" : "inline"} class="h-full w-full object-cover" src={src} alt={alt} onload={() =>{
            isLoading = false
        }}/>
    {#if isLoading}
        <div class="h-full w-full object-cover bg-gradient-to-r from-transparent via-gray-300 to-transparent animate-shimmer"></div>
    {/if}
</div>

<style>
    @keyframes shimmer {
        0% { background-position: 0%; }
        100% { background-position: 100%; }
    }

    .animate-shimmer {
        animation: shimmer 1.5s infinite;
    }
</style>