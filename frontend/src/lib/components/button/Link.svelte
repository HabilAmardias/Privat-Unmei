<script lang="ts">
    import { Button } from "bits-ui";
	import type { Snippet } from "svelte";
    type anchorProps = {
        href: string
        inverted?: boolean
        children: Snippet
    }
    let {href, children, inverted = false} : anchorProps = $props()
    
    let anchorClass = $state<string>("text-[var(--tertiary-color)] hover:text-[var(--secondary-color)]")
    if (inverted){
        anchorClass = "text-[var(--secondary-color)] hover:text-[var(--tertiary-color)]"
    }
</script>


<Button.Root href={href} class={anchorClass}>
    {@render children()}
</Button.Root>
