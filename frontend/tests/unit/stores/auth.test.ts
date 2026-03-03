import { describe, it, expect, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { mockNuxtImport } from '@nuxt/test-utils/runtime';
import { useAuthStore } from '../../../stores/auth';

const mockCookies: Record<string, { value: any }> = {};

mockNuxtImport('useCookie', () => {
    return (name: string) => {
        if (!mockCookies[name]) {
            mockCookies[name] = { value: '' };
        }
        return mockCookies[name];
    };
});

describe('Auth Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
        mockCookies['auth_token'] = { value: '' };
        mockCookies['auth_user'] = { value: '' };
    });

    it('initializes with default state', () => {
        const store = useAuthStore();
        expect(store.user).toBeNull();
        expect(store.token).toBeNull();
        expect(store.isLoading).toBe(false);
        expect(store.error).toBeNull();
    });

    it('isAuthenticated getter works', () => {
        const store = useAuthStore();
        expect(store.isAuthenticated).toBe(false);

        store.token = 'test-token';
        expect(store.isAuthenticated).toBe(false);

        store.user = { email: 'test@dpay.com', role: 'admin' };
        expect(store.isAuthenticated).toBe(true);
    });

    it('setAuth updates state and cookies', () => {
        const store = useAuthStore();
        const mockAuth = {
            email: 'test@dpay.com',
            role: 'user',
            token: 'fake-jwt-token'
        };

        store.setAuth(mockAuth);

        expect(store.user).toEqual({ email: 'test@dpay.com', role: 'user' });
        expect(store.token).toBe('fake-jwt-token');
        expect(store.error).toBeNull();

        expect(mockCookies['auth_token'].value).toBe('fake-jwt-token');
        expect(JSON.parse(mockCookies['auth_user'].value)).toEqual({
            email: 'test@dpay.com',
            role: 'user'
        });
    });

    it('clearAuth resets state and cookies', () => {
        const store = useAuthStore();
        store.setAuth({ email: 'test@dpay.com', role: 'user', token: 'token' });

        store.clearAuth();

        expect(store.user).toBeNull();
        expect(store.token).toBeNull();
        expect(store.error).toBeNull();
        expect(store.isLoading).toBe(false);

        expect(mockCookies['auth_token'].value).toBe('');
        expect(mockCookies['auth_user'].value).toBe('');
    });

    it('initialize populates state from cookies if present correctly', () => {
        const store = useAuthStore();
        mockCookies['auth_token'] = { value: 'saved-token' };
        mockCookies['auth_user'] = { value: JSON.stringify({ email: 'saved@dpay.com', role: 'admin' }) };

        store.initialize();

        expect(store.token).toBe('saved-token');
        expect(store.user).toEqual({ email: 'saved@dpay.com', role: 'admin' });
    });

    it('setLoading works', () => {
        const store = useAuthStore();
        store.setLoading(true);
        expect(store.isLoading).toBe(true);
    });

    it('setError and clearError work', () => {
        const store = useAuthStore();
        store.setError('Failed login');
        expect(store.error).toBe('Failed login');

        store.clearError();
        expect(store.error).toBeNull();
    });
});
