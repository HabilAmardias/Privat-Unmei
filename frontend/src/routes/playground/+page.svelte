<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import Link from '$lib/components/button/Link.svelte';
	import toast from 'svelte-french-toast';
	import Image from '$lib/components/image/Image.svelte';
	import Menubar from '$lib/components/menubar/Menubar.svelte';
	import Menu from '$lib/components/menubar/Menu.svelte';
	import imageSrc from '$lib/images/svelte-welcome.png';
	import Datepicker from '$lib/components/calendar/Datepicker.svelte';
	import RatingGroup from '$lib/components/rating/RatingGroup.svelte';
	import Select from '$lib/components/select/Select.svelte';
	import Combobox from '$lib/components/search/Combobox.svelte';

	let selectedValue = $state<string>('');
    let keyword = $state<string>("");
    let searchValue = $state<string>("");

	function openToastSuccess() {
		toast.success('Success!', {
			position: 'top-right'
		});
	}

	function selectOnChange(val: string) {
		selectedValue = val;
		console.log($state.snapshot(selectedValue));
	}

    function searchOnChange(val :string){
        searchValue = val
        console.log($state.snapshot(searchValue))
    }

    const searchLabel = "Search"
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

<div>
	<Button>Button Example</Button>
	<Button onClick={openToastSuccess}>Click for success toast</Button>
	<span class="bg[var(--secondary-color)] rounded-md p-2">
		<Link theme="dark" href="/">Home</Link>
	</span>
	<span class="rounded-md bg-[var(--tertiary-color)] p-2">
		<Link href="/">Home</Link>
	</span>
	<Image src={imageSrc} alt="test-image" height={256} width={256} round="full" />
	<Menubar>
		<Menu separator={true}>
			<Link href="/">Home</Link>
		</Menu>
		<Menu>
			<Link href="/playground">Playground</Link>
		</Menu>
	</Menubar>
	<Datepicker
		onChange={(date) => console.log(date?.toString())}
		dows={[0, 1, 2]}
		label="Pick Your Reservation Date"
	/>
	<RatingGroup onChange={(rate) => console.log(rate)} />
	<Select 
    options={themes} 
    value={selectedValue} 
    onChange={selectOnChange} />
	<Combobox 
    keyword={keyword} 
    label={searchLabel} 
    items={themes}
    onValueChange={searchOnChange}/>
</div>
