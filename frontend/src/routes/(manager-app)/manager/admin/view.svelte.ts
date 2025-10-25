export class adminLayoutView {
	menuOpen = $state<boolean>(false);
	handleMenu() {
		this.menuOpen = !this.menuOpen;
	}
}
