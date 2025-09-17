<script lang="ts">
    import { Button } from "bits-ui";
	import type { Snippet } from "svelte";
    type anchorProps = {
        href: string
        children: Snippet
    }
    let {href, children} : anchorProps = $props()
</script>


<Button.Root href={href} class="text-[var(--dark-color)] hover:text-[var(--primary-color)]">
    {@render children()}
</Button.Root>
