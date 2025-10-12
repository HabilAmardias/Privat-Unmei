export class VerifyView {
	openDialog = $state<boolean>(false);
	setOpenDialog(b: boolean) {
		this.openDialog = b;
	}
}
