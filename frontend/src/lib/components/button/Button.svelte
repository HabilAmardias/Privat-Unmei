<script lang="ts">
    import { Button } from "bits-ui";
	import type { Snippet } from "svelte";
    export type ButtonProps = {
        children: Snippet
        onClick?: () => void
        className: string
        formAction? : string
    }

    let {children, onClick, formAction, className} : ButtonProps = $props()
</script>

<Button.Root formaction={formAction} onclick={onClick} class={className}>
    {@render children()}
</Button.Root>