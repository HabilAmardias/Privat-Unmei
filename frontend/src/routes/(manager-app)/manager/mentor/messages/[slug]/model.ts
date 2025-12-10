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

export type MentorProfile = {
	id: string;
	profile_image: string;
	name: string;
	bio: string;
	years_of_experience: number;
	degree: string;
	major: string;
	campus: string;
	public_id: string;
};
