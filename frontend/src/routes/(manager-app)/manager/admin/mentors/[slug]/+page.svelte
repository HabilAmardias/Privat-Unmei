<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { goto } from '$app/navigation';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { MentorDetailView } from './view.svelte';
	import DownloadLink from '$lib/components/button/DownloadLink.svelte';
	import AlertDialog from '$lib/components/dialog/AlertDialog.svelte';

	const View = new MentorDetailView();
	let { data }: PageProps = $props();

	onMount(() => {
		if (!data.isVerified) {
			goto('/admin/verify', { replaceState: true });
		}
		View.setIsDesktop(window.innerWidth >= 768);
	});
</script>

{#snippet dialogTitle()}
	Delete Mentor Confirmation
{/snippet}

{#snippet dialogDescription()}
	Irreversible action, are you sure want to proceed?
{/snippet}

<div class="flex flex-col gap-4 p-4">
	<div class="flex gap-4">
		<CldImage
			className="rounded-full"
			height={View.size}
			width={View.size}
			src={data.profile.profile_image}
		/>
		<div class="flex w-full items-center justify-between">
			<div class="flex flex-col gap-1">
				<p class="font-bold text-[var(--tertiary-color)]">{data.profile.name}</p>
				<p>{data.profile.email}</p>
				<DownloadLink
					download="resume"
					className="cursor-pointer text-[var(--tertiary-color)] hover:text-[var(--primary-color)]"
					href={data.profile.resume_file}>Download Resume</DownloadLink
				>
			</div>
			<AlertDialog
				action="?/deleteMentor"
				bind:open={View.alertOpen}
				enhancement={View.onDeleteMentor}
				title={dialogTitle}
				description={dialogDescription}>Delete Account</AlertDialog
			>
		</div>
	</div>
	<div class="flex justify-between text-center">
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">YoE:</p>
			<p>{data.profile.years_of_experience}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Campus:</p>
			<p>{data.profile.campus}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Degree:</p>
			<p>{data.profile.degree}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Major:</p>
			<p>{data.profile.major}</p>
		</div>
	</div>
	<div class="flex flex-col">
		<p class="font-bold text-[var(--tertiary-color)]">Bio:</p>
		<p>{data.profile.bio}</p>
	</div>
</div>
