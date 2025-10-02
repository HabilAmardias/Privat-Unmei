<script lang="ts">
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { View } from './view.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import { Pencil } from '@lucide/svelte';
	import Textarea from '$lib/components/form/Textarea.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import FileInput from '$lib/components/form/FileInput.svelte';
	import Select from '$lib/components/select/Select.svelte';
	import { statusOptions } from './model';
	import Search from '$lib/components/search/Search.svelte';

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

{#if View.isEdit}
	<form action="?/updateProfile" class="flex h-full flex-col justify-center gap-4 p-8">
		<div class="flex items-center gap-4">
			<FileInput bind:files={View.profileImage} id="profile_image" name="profile_image">
				<div class="group relative inline-block overflow-hidden rounded-full">
					<CldImage
						src={data.resBody.data.profile_image}
						width={View.size}
						height={View.size}
						className="rounded-full shadow-2xl border-gray-400 brightness-60 md:brightness-100 md:border-none md:shadow-none md:hover:shadow-2xl md:group-hover:border-gray-400 md:transition-all md:duration-300 md:group-hover:brightness-60"
					/>
					<div
						class="absolute inset-0 flex items-center justify-center bg-opacity-0 transition-all duration-300 group-hover:bg-opacity-30"
					>
						<Pencil
							class="text-white md:scale-50 md:transform md:opacity-0 md:transition-all md:duration-300 md:group-hover:scale-100 md:group-hover:opacity-100"
						/>
					</div>
				</div>
			</FileInput>
			<div class="flex flex-col gap-1">
				<Input id="name" name="name" type="text" bind:value={View.name} />
				<p class="text-md">{data.resBody.data.email}</p>
			</div>
		</div>
		<div class="flex flex-col gap-2">
			<Textarea id="bio" name="bio" bind:value={View.bio}>Bio:</Textarea>
		</div>
		<div class="flex gap-1">
			<Button onClick={() => View.setIsEdit()}>Cancel</Button>
			<Button formAction="?/updateProfile" type="submit">Submit</Button>
		</div>
	</form>
{:else}
	<div class="flex h-full flex-col gap-4 p-8">
		<div class="flex items-center gap-4">
			<CldImage
				src={data.resBody.data.profile_image}
				width={View.size}
				height={View.size}
				className="rounded-full"
			/>
			<div class="flex flex-col gap-1">
				<div class="flex gap-1">
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
				</div>
				<p class="text-md">{data.resBody.data.email}</p>
			</div>
		</div>
		<div class="flex flex-col gap-2">
			<p>Bio:</p>
			<p class="text-justify">{data.resBody.data.bio}</p>
		</div>
		<div>
			<h3 class="text-xl font-bold text-[var(--tertiary-color)]">Orders</h3>
			<form class="flex items-center gap-4" action="?/myOrders">
				<Input width={300} placeholder="Enter a Keyword" id="search" name="search" type="text" />
				<Select defaultLable="Select status" options={statusOptions} bind:value={View.status} />
				<Button type="submit" formAction="?/myOrders">Search</Button>
			</form>
		</div>
	</div>
{/if}
