<script lang="ts">
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { profileView } from './view.svelte';
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
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import Loading from '$lib/components/loader/Loading.svelte';
	import Image from '$lib/components/image/Image.svelte';
	import { CreateToast, DismissToast } from '$lib/utils/helper';

	let { data }: PageProps = $props();
	const View = new profileView(data.profile, data.orders);

	onMount(() => {
		View.setIsDesktop(window.innerWidth >= 768);
		function isDesktop() {
			View.setIsDesktop(window.innerWidth >= 768);
		}
		window.addEventListener('resize', isDesktop);

		return () => {
			window.removeEventListener('resize', isDesktop);
		};
	});

	function onCancel() {
		View.setProfileImage(undefined);
		View.setBio(data.profile.bio);
		View.setName(data.profile.name);
		View.setNameError(undefined);
		View.setBioError(undefined);
		View.setIsEdit();
	}

	function onVerifySubmit(args: EnhancementArgs) {
		View.setVerifyIsLoading(true);
		const loadID = CreateToast('loading', 'sending....');
		return async ({ result }: EnhancementReturn) => {
			View.setVerifyIsLoading(false);
			DismissToast(loadID);
			if (result.type === 'success') {
				CreateToast('success', result.data?.message);
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	}

	function onUpdateOrders(args: EnhancementArgs) {
		View.setOrdersIsLoading(true);
		args.formData.append('last_id', `${View.lastID}`);
		return async ({ result }: EnhancementReturn) => {
			View.setOrdersIsLoading(false);
			if (result.type === 'success') {
				View.setOrders(result.data?.orders);
				View.setTotalRow(result.data?.totalRow);
				CreateToast('success', result.data?.message);
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	}

	function onUpdateProfile(args: EnhancementArgs) {
		const loadID = CreateToast('loading', 'updating....');
		return async ({ result, update }: EnhancementReturn) => {
			View.setIsEdit();
			View.setProfileIsLoading(true);
			await update({ reset: false });
			View.setProfileImage(undefined);
			View.setBio(data.profile.bio);
			View.setName(data.profile.name);
			View.setNameError(undefined);
			View.setBioError(undefined);
			View.setProfileIsLoading(false);
			DismissToast(loadID);
			if (result.type === 'success') {
				CreateToast('success', 'update profile success');
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	}
</script>

<svelte:head>
	<title>My Profile - Privat Unmei</title>
	<meta name="description" content="My Profile - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

{#if View.isEdit}
	<form
		use:enhance={onUpdateProfile}
		action="?/updateProfile"
		method="POST"
		enctype="multipart/form-data"
		class="flex h-full flex-col justify-center gap-4 p-4"
	>
		<div class="flex items-center gap-4">
			<FileInput accept="image/png" bind:files={View.profileImage} id="profile_image" name="file">
				<div class="group relative inline-block overflow-hidden rounded-full">
					{#if View.profileImage}
						<Image
							src={URL.createObjectURL(View.profileImage[0])}
							width={View.size}
							height={View.size}
							className="rounded-full shadow-2xl border-gray-400 brightness-60 md:brightness-100 md:border-none md:shadow-none md:hover:shadow-2xl md:group-hover:border-gray-400 md:transition-all md:duration-300 md:group-hover:brightness-60"
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
				<Input
					err={View.nameError}
					onBlur={() => View.nameOnBlur()}
					id="name"
					name="name"
					type="text"
					bind:value={View.name}
				/>
				<p class="text-md">{data.profile.email}</p>
			</div>
		</div>
		<div class="flex flex-col gap-2">
			<Textarea
				err={View.bioError}
				onBlur={() => View.bioOnBlur()}
				id="bio"
				name="bio"
				bind:value={View.bio}>Bio:</Textarea
			>
		</div>
		<div class="flex gap-1">
			<Button type="button" onClick={onCancel}>Cancel</Button>
			<Button disabled={View.updateProfileDisable} formAction="?/updateProfile" type="submit"
				>Submit</Button
			>
		</div>
	</form>
{:else}
	<div class="flex h-full flex-col gap-4 p-4">
		<div class="flex flex-col gap-4">
			{#if View.profileIsLoading}
				<Loading />
			{:else}
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
					</div>
				</div>
				<div class="flex flex-col gap-2">
					<b class="text-xl text-[var(--tertiary-color)]">Bio:</b>
					<p class="text-justify">{data.profile.bio}</p>
				</div>
			{/if}
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
						<b class="mx-auto self-center text-[var(--tertiary-color)]">No orders found</b>
					{:else}
						<ScrollArea class="flex-1" orientation="horizontal" viewportClasses="max-h-full">
							{#each View.orders as order (order.id)}
								<div>
									<p>{order.course_name}</p>
									<p>{order.mentor_name}</p>
									<p>{order.mentor_email}</p>
									<p>{order.total_price}</p>
									<p>{order.status}</p>
								</div>
							{/each}
						</ScrollArea>
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
		{:else}
			<form
				use:enhance={onVerifySubmit}
				method="POST"
				action="?/sendVerification"
				class="flex w-full flex-1 items-center justify-center"
			>
				<Button disabled={View.verifyIsLoading} type="submit" formAction="?/sendVerification"
					>Send Verification Link</Button
				>
			</form>
		{/if}
	</div>
{/if}
