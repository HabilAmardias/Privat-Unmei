<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import Link from '$lib/components/button/Link.svelte';
	import toast from 'svelte-french-toast';
	import Image from '$lib/components/image/Image.svelte';
	import imageSrc from '$lib/images/svelte-welcome.png';
	import Datepicker from '$lib/components/calendar/Datepicker.svelte';
	import RatingGroup from '$lib/components/rating/RatingGroup.svelte';
	import Select from '$lib/components/select/Select.svelte';
	import Search from '$lib/components/search/Search.svelte';
	import Pagination from '$lib/components/pagination/Pagination.svelte';
	import AlertDialog from '$lib/components/dialog/AlertDialog.svelte';
	import Dialog from '$lib/components/dialog/Dialog.svelte';
	import InputSecret from '$lib/components/form/InputSecret.svelte';
	import Card from '$lib/components/card/Card.svelte';
	import Input from '$lib/components/form/Input.svelte';

	let selectedValue = $state<string>('');
	let keyword = $state<string>('');
	let searchValue = $state<string>('');
	let openAlert = $state<boolean>(false);
	let password = $state<string>('');

	function openToastSuccess() {
		toast.success('Success!', {
			position: 'top-right'
		});
	}

	function rateOnChange(num: number) {
		console.log(num);
	}

	function selectOnChange(val: string) {
		selectedValue = val;
		console.log($state.snapshot(selectedValue));
	}

	function searchOnChange(val: string) {
		searchValue = val;
		console.log($state.snapshot(searchValue));
	}

	function alertOnSubmit(e: SubmitEvent & { currentTarget: EventTarget & HTMLFormElement }) {
		e.preventDefault();
		openAlert = false;
	}

	function pageOnChange(num: number) {
		console.log(num);
	}

	const themes = [
		{ value: 'light-monochrome', label: 'Light Monochrome' },
		{ value: 'dark-green', label: 'Dark Green' },
		{ value: 'svelte-orange', label: 'Svelte Orange' },
		{ value: 'punk-pink', label: 'Punk Pink' },
		{ value: 'ocean-blue', label: 'Ocean Blue' },
		{ value: 'sunset-orange', label: 'Sunset Orange' },
		{ value: 'sunset-red', label: 'Sunset Red' },
		{ value: 'forest-green', label: 'Forest Green' },
		{ value: 'lavender-purple', label: 'Lavender Purple' },
		{ value: 'mustard-yellow', label: 'Mustard Yellow' },
		{ value: 'slate-gray', label: 'Slate Gray' },
		{ value: 'neon-green', label: 'Neon Green' },
		{ value: 'coral-reef', label: 'Coral Reef' },
		{ value: 'midnight-blue', label: 'Midnight Blue' },
		{ value: 'crimson-red', label: 'Crimson Red' },
		{ value: 'mint-green', label: 'Mint Green' },
		{ value: 'pastel-pink', label: 'Pastel Pink' },
		{ value: 'golden-yellow', label: 'Golden Yellow' },
		{ value: 'deep-purple', label: 'Deep Purple' },
		{ value: 'turquoise-blue', label: 'Turquoise Blue' },
		{ value: 'burnt-orange', label: 'Burnt Orange' }
	];
</script>

{#snippet alertTitle()}
	Confirm your action
{/snippet}

{#snippet alertDescription()}
	Are you sure you want to do this?
{/snippet}

{#snippet dialogTitle()}
	Title Test
{/snippet}

{#snippet dialogDescription()}
	Test description
{/snippet}

<svelte:head>
	<title>About</title>
	<meta name="playground" content="playground for components" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<Button>Button Example</Button>
<Button onClick={openToastSuccess}>Click for success toast</Button>
<span class="bg[var(--secondary-color)] rounded-md p-2">
	<Link theme="dark" href="/">Home</Link>
</span>
<span class="rounded-md bg-[var(--tertiary-color)] p-2">
	<Link href="/">Home</Link>
</span>
<Image src={imageSrc} alt="test-image" height={256} width={256} round="full" />

<RatingGroup onChange={rateOnChange} />

<Select options={themes} value={selectedValue} onChange={selectOnChange} />

<Search
	{keyword}
	label="Search"
	items={themes}
	onValueChange={searchOnChange}
	onKeywordChange={(e) => console.log(e.currentTarget.value)}
/>

<Search
	{keyword}
	label="Search without Dropdown"
	onKeywordChange={(e) => console.log(e.currentTarget.value)}
/>

<Pagination onPageChange={pageOnChange} count={100} perPage={15} />

<AlertDialog
	bind:open={openAlert}
	onSubmit={alertOnSubmit}
	description={alertDescription}
	title={alertTitle}
>
	Submit
</AlertDialog>
<Dialog {dialogTitle} dialogContent={dialogDescription}>Test</Dialog>

<Card>
	<Input placeholder="Email" type="email" name="email" id="email" />
	<InputSecret id="password" placeholder="Password" name="password" bind:value={password} />
	<Datepicker onChange={(date) => console.log(date?.toString())} dows={[0, 1, 2]} />
</Card>
