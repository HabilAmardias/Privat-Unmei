<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import Select from '$lib/components/select/Select.svelte';
	import { degreeOpts } from './constants';
	import { CreateMentorView } from './view.svelte';
	const View = new CreateMentorView();
</script>

<div class="flex h-full flex-col gap-4 p-4">
	<h3 class="text-xl font-bold text-[var(--tertiary-color)]">Create New Mentor</h3>
	<form action="?/createMentor" method="POST">
		<Input type="text" placeholder="Input mentor name" name="name" id="name" />
		<Input type="email" placeholder="Input mentor email" name="email" id="name" />
		<div class="flex items-center gap-4">
			<p>test password</p>
			<Button formAction="?/generatePassword">Generate Password</Button>
		</div>
		<Select
			options={degreeOpts}
			defaultLable="Choose mentor degree"
			name="degree"
			value={View.degree}
		/>
		<Input type="text" name="major" id="major" placeholder="Insert mentor major" />
	</form>
</div>
