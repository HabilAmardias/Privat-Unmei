export class MentorLayoutView {
	menuOpen = $state<boolean>(false);
	handleMenu() {
		this.menuOpen = !this.menuOpen;
	}
}
