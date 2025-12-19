import type { Fetch, ServerResponse } from '$lib/types';
import type {
	CourseDetail,
	MentorPaymentInfo,
	ScheduleSlot,
	MentorScheduleInfo,
	DowRes,
	OperationalCostRes,
	DiscountRes
} from './model';
import { FetchData } from '$lib/utils';

class CreateRequestController {
	async getCourseDetail(fetch: Fetch, id: string) {
		const url = `/api/v1/courses/${id}`;
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<CourseDetail> = await res?.json();
		return { success, message, status, resBody };
	}
	async getAvailableDayOfWeek(fetch: Fetch, id: string) {
		const url = `/api/v1/courses/${id}/mentor-availability`;
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<DowRes> = await res?.json();
		return { success, message, status, resBody };
	}
	async getAdditionalCost(fetch: Fetch) {
		const url = '/api/v1/additional-cost/operational';
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<OperationalCostRes> = await res?.json();
		return { success, message, status, resBody };
	}
	async getDiscount(fetch: Fetch, req?: Request) {
		let participant: number = 1;
		if (req) {
			const formData = await req.formData();
			const participantInput = formData.get('participant');
			if (participantInput) {
				participant = parseInt(participantInput as string);
			}
		}
		const url = `/api/v1/discounts/final-discount/${participant}`;
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<DiscountRes> = await res?.json();
		return { success, message, status, resBody };
	}
	async createRequest(fetch: Fetch, id: string, req: Request) {
		const url = `/api/v1/courses/${id}/course-requests`;
		const formData = await req.formData();
		const schedules = formData.get('schedules');
		if (!schedules) {
			return { success: false, message: 'please insert at least one schedule', status: 400 };
		}
		const parsedSchedules: ScheduleSlot[] = JSON.parse(schedules as string);
		const payment = formData.get('payment');
		if (!payment) {
			return { success: false, message: 'please insert payment method', status: 400 };
		}
		const participant = formData.get('participant');
		if (!participant) {
			return { success: false, message: 'please insert number of participant', status: 400 };
		}

		const reqBody = JSON.stringify({
			preferred_slots: parsedSchedules,
			payment_method_id: parseInt(payment as string),
			number_of_participant: parseInt(participant as string)
		});
		const { success, status, message, res } = await FetchData(fetch, url, 'POST', reqBody);
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<{ id: number }> = await res?.json();
		return { success, message, status, resBody };
	}
	async getMentorPayments(fetch: Fetch, id: string) {
		const url = `/api/v1/mentors/${id}/payment-methods`;
		const { success, res, status, message } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, status, message };
		}
		const resBody: ServerResponse<MentorPaymentInfo[]> = await res?.json();
		return { success, status, message, resBody };
	}
	async getMentorSchedules(fetch: Fetch, id: string) {
		const url = `/api/v1/mentors/${id}/availability`;
		const { success, res, status, message } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, status, message };
		}
		const resBody: ServerResponse<MentorScheduleInfo[]> = await res?.json();
		return { success, status, message, resBody };
	}
}

export const controller = new CreateRequestController();
