/**
 * Tag feature types
 */

/**
 * Tag interface
 */
export interface Tag {
	id: string
	name: string
	color?: string
	userId: string
	createdAt: string
	updatedAt: string
}

/**
 * Tag store state interface
 */
export interface TagState {
	tags: Tag[]
	currentTag: null | Tag
	isLoading: boolean
	error: null | string
}
