import { env } from '$env/dynamic/public';

export const API_BASE_URL = env.PUBLIC_API_BASE_URL || 'http://localhost:8080/api';
export const APP_NAME = env.PUBLIC_APP_NAME || 'Task Manager';
export const APP_VERSION = env.PUBLIC_APP_VERSION || '1.0.0';
