<script lang="ts">
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { View } from './view.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import { Pencil } from '@lucide/svelte';
	import Textarea from '$lib/components/form/Textarea.svelte';
	import Input from '$lib/components/form/Input.svelte';

	let { data }: PageProps = $props();
	onMount(() => {
		View.setBio(data.resBody.data.bio);
		View.setName(data.resBody.data.name);
		View.setIsDesktop(window.innerWidth >= 768);
		function isDesktop() {
			View.setIsDesktop(window.innerWidth >= 768);
		}
		window.addEventListener('resize', isDesktop);

		return () => {
			window.removeEventListener('resize', isDesktop);
		};
	});
</script>

<div class="flex h-full flex-col gap-4 p-8">
	<div class="flex items-center gap-4">
		<CldImage
			src={data.resBody.data.profile_image}
			width={View.size}
			height={View.size}
			round="full"
		/>
		<div class="flex flex-col gap-1">
			<div class="flex gap-1">
				{#if !View.isEdit}
					<b class="text-xl text-[var(--tertiary-color)]">{data.resBody.data.name}</b>
					<Button
						onClick={() => View.setIsEdit()}
						type="button"
						withPadding={false}
						withBg={false}
						textColor="dark"
					>
						<Pencil width={24} height={24} />
					</Button>
				{:else}
					<Input id="name" name="name" type="text" bind:value={View.name} />
				{/if}
			</div>
			<p class="text-md">{data.resBody.data.email}</p>
		</div>
	</div>
	<div class="flex flex-col gap-2">
		{#if !View.isEdit}
			<p>Bio:</p>
			<p class="text-justify">{data.resBody.data.bio}</p>
		{:else}
			<Textarea id="bio" name="bio" bind:value={View.bio}>Bio:</Textarea>
		{/if}
	</div>
	{#if View.isEdit}
		<div class="flex gap-1">
			<Button onClick={() => View.setIsEdit()}>Cancel</Button>
			<Button formAction="?/updateProfile" type="submit">Submit</Button>
		</div>
	{/if}
</div>
