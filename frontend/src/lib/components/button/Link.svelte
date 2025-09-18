<script lang="ts">
    import { Button } from "bits-ui";
	import type { Snippet } from "svelte";
    type anchorProps = {
        href: string
        className: string
        children: Snippet
    }
    let {href, children, className} : anchorProps = $props()
</script>


<Button.Root href={href} class={className}>
    {@render children()}
</Button.Root>
