<script lang="ts">
    import { Menubar } from "bits-ui";
	import type { Snippet } from "svelte";
    type menubarProps = {
        children: Snippet
    }
    let {children} : menubarProps = $props()
</script>

<Menubar.Root class="bg-[var(--tertiary-color)] p-2 flex gap-1">
    {@render children()}
</Menubar.Root>