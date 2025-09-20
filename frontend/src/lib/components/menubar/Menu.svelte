<script lang="ts">
    import { Menubar } from "bits-ui";
	import type { Snippet } from "svelte";
    type menuItemProps = {
        children: Snippet
        separator?: boolean
    }
    let {children, separator = true} : menuItemProps = $props()
</script>

<Menubar.Menu>
    <Menubar.Trigger>
        {@render children()}
    </Menubar.Trigger>
    {#if separator}
        <Menubar.Separator/>
    {/if}
</Menubar.Menu>