/**
 * Tag feature types
 */

/**
 * Tag interface
 */
export interface Tag {
	id: string;
	name: string;
	color?: string;
	userId: string;
	createdAt: string;
	updatedAt: string;
}

/**
 * Tag creation payload
 */
export interface CreateTagPayload {
	name: string;
	color?: string;
}

/**
 * Tag update payload
 */
export interface UpdateTagPayload {
	name?: string;
	color?: string;
}

/**
 * Tag store state interface
 */
export interface TagState {
	tags: Tag[];
	currentTag: Tag | null;
	isLoading: boolean;
	error: string | null;
}
