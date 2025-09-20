<script lang="ts">
    import { Menubar } from "bits-ui";
	import type { Snippet } from "svelte";
    type menubarProps = {
        children: Snippet
    }
    let {children} : menubarProps = $props()
</script>

<Menubar.Root class="bg-[var(--tertiary-color)] p-3 rounded-lg flex gap-1 mx-auto">
    {@render children()}
</Menubar.Root>