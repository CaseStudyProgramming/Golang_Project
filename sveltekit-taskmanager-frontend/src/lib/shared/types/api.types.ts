/**
 * Shared API types and interfaces
 */

export interface ApiResponse<T> {
	data: T;
	error?: string;
	message?: string;
}

export interface PaginationParams {
	page?: number;
	limit?: number;
	sort?: string;
	order?: 'asc' | 'desc';
}

export interface PaginatedResponse<T> {
	data: T[];
	total: number;
	page: number;
	limit: number;
	totalPages: number;
}
