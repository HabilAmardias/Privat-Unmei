<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import InputSecret from '$lib/components/form/InputSecret.svelte';
	import { ManagerAuthView } from './view.svelte';
	import { enhance } from '$app/forms';
	import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { loadingStore } from '$lib/stores/LoadingStore.svelte';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { PrivatUnmeiLogo } from '$lib/utils/constants';
	import { adminLogin, mentorLogin } from './constants';
	import { adminRole, mentorRole } from '$lib/utils/constants';
	import NavigationButton from '$lib/components/button/NavigationButton.svelte';
	import type { PageProps } from './$types';
	import { CreateToast, DismissToast } from '$lib/utils/helper';

	const View = new ManagerAuthView();

	let { data }: PageProps = $props();

	onMount(() => {
		if (data.authToken) {
			switch (data.role) {
				case adminRole:
					goto('/manager/admin', { replaceState: true });
					break;
				case mentorRole:
					goto('/manager/mentor', { replaceState: true });
					break;
				default:
					goto('/courses', { replaceState: true });
					break;
			}
		}
		View.setIsDesktop(window.innerWidth >= 768);
		function isDesktop() {
			View.setIsDesktop(window.innerWidth >= 768);
		}
		if (loadingStore.logOutLoadID) {
			DismissToast(loadingStore.logOutLoadID);
			loadingStore.removeLogOutLoadID();
			CreateToast('success', 'logout success');
		}
		window.addEventListener('resize', isDesktop);
		return () => {
			window.removeEventListener('resize', isDesktop);
		};
	});

	function onAdminLoginSubmit(args: EnhancementArgs) {
		View.setIsLoading(true);
		const loadID = CreateToast('loading', 'logging in....');
		return async ({ result }: EnhancementReturn) => {
			if (result.type === 'success') {
				await goto('/manager/admin', { replaceState: true });
				View.setIsLoading(false);
				DismissToast(loadID);
				CreateToast('success', 'login success');
			}
			if (result.type === 'failure') {
				View.setIsLoading(false);
				DismissToast(loadID);
				CreateToast('error', result.data?.message);
			}
		};
	}
	function onMentorLoginSubmit(args: EnhancementArgs) {
		View.setIsLoading(true);
		const loadID = CreateToast('loading', 'logging in....');
		return async ({ result }: EnhancementReturn) => {
			if (result.type === 'success') {
				await goto('/manager/mentor', { replaceState: true });
				View.setIsLoading(false);
				DismissToast(loadID);
				CreateToast('success', 'login success');
			}
			if (result.type === 'failure') {
				View.setIsLoading(false);
				DismissToast(loadID);
				CreateToast('error', result.data?.message);
			}
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
