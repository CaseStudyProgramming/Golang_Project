import { env } from '$env/dynamic/private';

export const API_BASE_URL = env.API_BASE_URL || 'http://localhost:8080/api';
export const APP_NAME = env.APP_NAME || 'Task Manager';
export const APP_VERSION = env.APP_VERSION || '1.0.0';
