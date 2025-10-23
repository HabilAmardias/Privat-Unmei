export type paymentMethod = {
	payment_method_id: number;
	payment_method_name: string;
};

export type paymentMethodOpts = {
	value: string;
	label: string;
};

export type mentorPaymentMethods = {
	payment_method_id: number;
	account_number: string;
};

export type generatedPassword = {
	password: string;
};

export type TimeOnly = {
	hour: number;
	minute: number;
	second: number;
};

export function stringToTimeOnly(s: string): TimeOnly {
	const [hour, minute, second] = s.split(':');
	return {
		hour: parseInt(hour),
		minute: parseInt(minute),
		second: second ? parseInt(second) : 0
	};
}

export function TimeOnlyToString(t: TimeOnly): string {
	let res: string = '';
	res += t.hour < 10 ? `0${t.hour}:` : `${t.hour}:`;
	res += t.minute < 10 ? `0${t.minute}:` : `${t.minute}:`;
	res += t.second < 10 ? `0${t.second}` : `${t.second}`;
	return res;
}

export type MentorSchedule = {
	day_of_week: number;
	start_time: TimeOnly;
	end_time: TimeOnly;
};
