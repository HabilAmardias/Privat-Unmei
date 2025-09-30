class LoadingStore {
	logOutLoadID = $state<string | undefined>();
	setLogOutLoadID(id: string) {
		this.logOutLoadID = id;
	}
	removeLogOutLoadID() {
		this.logOutLoadID = undefined;
	}
}
export const loadingStore = new LoadingStore();
