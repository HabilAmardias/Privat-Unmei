export type ChatroomInfo = {
	id: number;
	user_id: string;
	username: string;
	public_id: string;
	profile_image: string;
};

export type MessageInfo = {
	id: number;
	sender_id: string;
	sender_name: string;
	sender_public_id: string;
	chatroom_id: number;
	content: string;
};

export type StudentProfile = {
	id: string;
	name: string;
	bio: string;
	profile_image: string;
	public_id: string;
};
