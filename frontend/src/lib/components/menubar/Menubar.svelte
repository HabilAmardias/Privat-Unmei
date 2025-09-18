<script lang="ts">
    import { Menubar } from "bits-ui";
	import type { Snippet } from "svelte";
    type menubarProps = {
        children: Snippet
        className: string
    }
    let {children, className} : menubarProps = $props()
</script>

<Menubar.Root class={className}>
    {@render children()}
</Menubar.Root>