import { env } from '$env/dynamic/public';

export const API_BASE_URL = env.VITE_API_BASE_URL || 'http://localhost:8080/api';
export const APP_NAME = env.VITE_APP_NAME || 'Task Manager';
export const APP_VERSION = env.VITE_APP_VERSION || '1.0.0';
