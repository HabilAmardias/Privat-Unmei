<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import InputSecret from '$lib/components/form/InputSecret.svelte';
	import { ManagerAuthView } from './view.svelte';
	import { enhance } from '$app/forms';
	import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
	import toast from 'svelte-french-toast';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { loadingStore } from '$lib/stores/LoadingStore.svelte';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { PrivatUnmeiLogo } from '$lib/utils/constants';
	import { adminLogin, mentorLogin } from './constants';
	import { adminRole, mentorRole } from '$lib/utils/constants';
	import NavigationButton from '$lib/components/button/NavigationButton.svelte';
	import type { PageProps } from './$types';

	const View = new ManagerAuthView();

	let { data }: PageProps = $props();

	onMount(() => {
		if (data.authToken) {
			if (data.role === adminRole) {
				goto('/manager/admin', { replaceState: true });
				return;
			}
			if (data.role === mentorRole) {
				goto('/manager/mentor', { replaceState: true });
				return;
			}
		}
		View.setIsDesktop(window.innerWidth >= 768);
		function isDesktop() {
			View.setIsDesktop(window.innerWidth >= 768);
		}
		if (loadingStore.logOutLoadID) {
			toast.dismiss(loadingStore.logOutLoadID);
			loadingStore.removeLogOutLoadID();
			toast.success('log out success', { position: 'top-right' });
		}
		window.addEventListener('resize', isDesktop);
		return () => {
			window.removeEventListener('resize', isDesktop);
		};
	});

	function onAdminLoginSubmit(args: EnhancementArgs) {
		View.setIsLoading(true);
		const loadID = toast.loading('logging in.....', { position: 'top-right' });
		return async ({ result, update }: EnhancementReturn) => {
			if (result.type === 'success') {
				await goto('/manager/admin', { replaceState: true });
				View.setIsLoading(false);
				toast.dismiss(loadID);
				toast.success('login success', {
					position: 'top-right'
				});
			}
			if (result.type === 'failure') {
				View.setIsLoading(false);
				toast.dismiss(loadID);
				toast.error(result.data?.message, {
					position: 'top-right'
				});
			}
			update();
		};
	}
	function onMentorLoginSubmit(args: EnhancementArgs) {
		View.setIsLoading(true);
		const loadID = toast.loading('logging in.....', { position: 'top-right' });
		return async ({ result, update }: EnhancementReturn) => {
			if (result.type === 'success') {
				await goto('/manager/mentor', { replaceState: true });
				View.setIsLoading(false);
				toast.dismiss(loadID);
				toast.success('login success', {
					position: 'top-right'
				});
			}
			if (result.type === 'failure') {
				View.setIsLoading(false);
				toast.dismiss(loadID);
				toast.error(result.data?.message, {
					position: 'top-right'
				});
			}
			update();
		};
	}
</script>

<svelte:head>
	<title>Manager Login - Privat Unmei</title>
	<meta name="description" content="Manager Login - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex h-screen w-full flex-col items-center justify-center gap-8 md:flex-row md:gap-0">
	<div class="hidden md:flex md:flex-1">
		<CldImage src={PrivatUnmeiLogo} width={400} height={125} />
	</div>
	<div class="block md:hidden">
		<CldImage src={PrivatUnmeiLogo} width={200} height={60} />
	</div>

	<div class="flex flex-col">
		<NavigationButton
			menus={[
				{
					header: 'Admin',
					onClick: () => View.switchForm(adminLogin)
				},
				{
					header: 'Mentor',
					onClick: () => View.switchForm(mentorLogin)
				}
			]}
		/>
		{#if View.loginMenu === mentorLogin}
			<form
				use:enhance={onMentorLoginSubmit}
				action="?/loginMentor"
				method="post"
				class="border-1 flex flex-col gap-4 rounded-bl-lg rounded-br-lg rounded-tr-lg border-[var(--tertiary-color)] p-4"
			>
				<h2 class="mb-3 text-2xl text-[var(--tertiary-color)]">Mentor Login</h2>
				<Input
					err={View.emailError}
					onBlur={() => View.emailOnBlur()}
					type="email"
					name="email"
					placeholder="Email"
					id="email"
					bind:value={View.email}
				/>
				<InputSecret
					err={View.passwordError}
					onBlur={() => View.passwordOnBlur()}
					id="password"
					placeholder="Password"
					name="password"
					bind:value={View.password}
				/>
				<Button disabled={View.loginDisabled} full={true} type="submit" formAction="?/loginMentor"
					>Login</Button
				>
			</form>
		{:else}
			<form
				use:enhance={onAdminLoginSubmit}
				action="?/loginAdmin"
				method="post"
				class="border-1 flex flex-col gap-4 rounded-bl-lg rounded-br-lg rounded-tr-lg border-[var(--tertiary-color)] p-4"
			>
				<h2 class="mb-3 text-2xl text-[var(--tertiary-color)]">Admin Login</h2>
				<Input
					err={View.emailError}
					onBlur={() => View.emailOnBlur()}
					type="email"
					name="email"
					placeholder="Email"
					id="email"
					bind:value={View.email}
				/>
				<InputSecret
					err={View.passwordError}
					onBlur={() => View.passwordOnBlur()}
					id="password"
					placeholder="Password"
					name="password"
					bind:value={View.password}
				/>
				<Button disabled={View.loginDisabled} full={true} type="submit" formAction="?/loginAdmin"
					>Login</Button
				>
			</form>
		{/if}
	</div>
</div>
