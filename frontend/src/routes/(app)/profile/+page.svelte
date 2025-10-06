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
	import Pagination from '$lib/components/pagination/Pagination.svelte';
	import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
	import { enhance } from '$app/forms';
	import toast from 'svelte-french-toast';
	import { ScrollArea } from 'bits-ui';
	import Loading from '$lib/components/loader/Loading.svelte';
	import Image from '$lib/components/image/Image.svelte';

	let { data }: PageProps = $props();

	onMount(() => {
		View.setBio(data.profile.bio);
		View.setName(data.profile.name);
		View.setIsDesktop(window.innerWidth >= 768);
		if (data.orders) {
			View.setOrders(data.orders.entries);
			View.setTotalRow(data.orders.page_info.total_row);
			View.setLastID(data.orders.page_info.last_id);
		}
		function isDesktop() {
			View.setIsDesktop(window.innerWidth >= 768);
		}
		window.addEventListener('resize', isDesktop);

		return () => {
			window.removeEventListener('resize', isDesktop);
		};
	});

	function onVerifySubmit(args: EnhancementArgs) {
		View.setVerifyIsLoading(true);
		const loadID = toast.loading('sending....', { position: 'top-right' });
		return async ({ result, update }: EnhancementReturn) => {
			View.setVerifyIsLoading(false);
			toast.dismiss(loadID);
			if (result.type === 'success') {
				toast.success(result.data?.message, { position: 'top-right' });
			}
			if (result.type === 'failure') {
				toast.error(result.data?.message, { position: 'top-right' });
			}
			update();
		};
	}

	function onUpdateOrders(args: EnhancementArgs) {
		View.setOrdersIsLoading(true);
		args.formData.append('last_id', `${View.lastID}`);
		return async ({ result, update }: EnhancementReturn) => {
			View.setOrdersIsLoading(false);
			if (result.type === 'success') {
				View.setOrders(result.data?.orders);
				View.setTotalRow(result.data?.totalRow);
			}
			if (result.type === 'failure') {
				toast.error(result.data?.message, { position: 'top-right' });
			}
			update({ reset: false });
		};
	}

	function onUpdateProfile(args: EnhancementArgs) {
		View.setProfileIsLoading(true);
		const loadID = toast.loading('updating...', { position: 'top-right' });
		return async ({ result, update }: EnhancementReturn) => {
			View.setProfileIsLoading(false);
			toast.dismiss(loadID);
			if (result.type === 'success') {
				toast.success('update profile success', { position: 'top-right' });
			}
			if (result.type === 'failure') {
				toast.error(result.data?.message, { position: 'top-right' });
			}
			View.setIsEdit();
			update({ reset: false });
		};
	}
</script>

{#if View.isEdit}
	<form
		use:enhance={onUpdateProfile}
		action="?/updateProfile"
		method="POST"
		enctype="multipart/form-data"
		class="flex h-full flex-col justify-center gap-4 p-4"
	>
		<div class="flex items-center gap-4">
			<FileInput bind:files={View.profileImage} id="profile_image" name="file">
				<div class="group relative inline-block overflow-hidden rounded-full">
					{#if View.profileImage}
						<Image
							src={URL.createObjectURL(View.profileImage[0])}
							width={View.size}
							height={View.size}
							round="full"
						/>
					{:else}
						<CldImage
							src={data.profile.profile_image}
							width={View.size}
							height={View.size}
							className="rounded-full shadow-2xl border-gray-400 brightness-60 md:brightness-100 md:border-none md:shadow-none md:hover:shadow-2xl md:group-hover:border-gray-400 md:transition-all md:duration-300 md:group-hover:brightness-60"
						/>
					{/if}
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
				<p class="text-md">{data.profile.email}</p>
			</div>
		</div>
		<div class="flex flex-col gap-2">
			<Textarea id="bio" name="bio" bind:value={View.bio}>Bio:</Textarea>
		</div>
		<div class="flex gap-1">
			<Button type="button" onClick={() => View.setIsEdit()}>Cancel</Button>
			<Button disabled={View.profileIsLoading} formAction="?/updateProfile" type="submit"
				>Submit</Button
			>
		</div>
	</form>
{:else}
	<div class="flex h-full flex-col gap-4 p-4">
		<div class="flex items-center gap-4">
			<CldImage
				src={data.profile.profile_image}
				width={View.size}
				height={View.size}
				className="rounded-full"
			/>
			<div class="flex flex-col gap-1">
				<div class="flex gap-1">
					<b class="text-xl text-[var(--tertiary-color)]">{data.profile.name}</b>
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
				<p class="text-md">{data.profile.email}</p>
				{#if data.userStatus !== 'verified'}
					<form use:enhance={onVerifySubmit} method="POST" action="?/sendVerification">
						<Button disabled={View.verifyIsLoading} type="submit" formAction="?/sendVerification"
							>Send Verification Link</Button
						>
					</form>
				{/if}
			</div>
		</div>
		<div class="flex flex-col gap-2">
			<b class="text-xl text-[var(--tertiary-color)]">Bio:</b>
			<p class="text-justify">{data.profile.bio}</p>
		</div>
		{#if data.userStatus === 'verified'}
			<div class="flex flex-1 flex-col gap-4">
				<h3 class="text-xl font-bold text-[var(--tertiary-color)]">Orders</h3>
				<form
					use:enhance={onUpdateOrders}
					class="grid grid-cols-3 gap-4"
					action="?/myOrders"
					method="POST"
				>
					<Input width="full" placeholder="Search" id="search" name="search" type="text" />
					<Select
						defaultLable="Status"
						name="status"
						options={statusOptions}
						bind:value={View.status}
					/>
					<Button disabled={View.ordersIsLoading} type="submit" full formAction="?/myOrders"
						>Search</Button
					>
				</form>
				<div class="flex flex-1">
					{#if View.ordersIsLoading}
						<Loading />
					{:else if !View.orders || View.orders.length === 0}
						<p>No orders found</p>
					{:else}
						<ScrollArea.Root class="h-full">
							<ScrollArea.Viewport class="h-full">
								{#each View.orders as order (order.id)}
									<div>
										<p>{order.course_name}</p>
										<p>{order.mentor_name}</p>
										<p>{order.mentor_email}</p>
										<p>{order.total_price}</p>
										<p>{order.status}</p>
									</div>
								{/each}
							</ScrollArea.Viewport>
						</ScrollArea.Root>
					{/if}
				</div>

				<form
					bind:this={View.paginationForm}
					action="?/myOrders"
					method="POST"
					class="flex items-center justify-center"
					use:enhance={onUpdateOrders}
				>
					<Pagination
						onPageChange={(num) => {
							View.onPageChange(num);
						}}
						pageNumber={View.pageNumber}
						count={View.totalRow}
						perPage={View.limit}
					/>
				</form>
			</div>
		{/if}
	</div>
{/if}
