import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mountSuspended, mockComponent, mockNuxtImport } from '@nuxt/test-utils/runtime';
import { flushPromises } from '@vue/test-utils';
import { useAuthStore } from '../../stores/auth';
import LoginForm from '../../components/LoginForm.vue';

const { mockNavigateTo, mockLoginApi } = vi.hoisted(() => ({
    mockNavigateTo: vi.fn(),
    mockLoginApi: vi.fn()
}));

mockNuxtImport('navigateTo', () => mockNavigateTo);

// We want the real useAuth composable to execute, 
// so we mock only the API layer
// We can't easily mock authApi.login directly if it's not exposed via nuxt imports,
// but we *can* intercept the fetch or mock the composable.
// Since it's a functional test, let's mock the `api` utility.
// Wait, `authApi` is in `~/utils/api.ts`.
vi.mock('~/utils/api', () => {
    return {
        authApi: {
            login: (args: any) => mockLoginApi(args)
        }
    };
});

mockComponent('InputField', () => import('../../components/InputField.vue'));

describe('Auth Flow Functional', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        const store = useAuthStore();
        store.clearAuth();
    });

    it('successful login flow updates store and navigates to index', async () => {
        const store = useAuthStore();
        expect(store.isAuthenticated).toBe(false);

        mockLoginApi.mockResolvedValueOnce({
            email: 'success@dpay.com',
            role: 'admin',
            token: 'real-token'
        });

        const wrapper = await mountSuspended(LoginForm);

        // Fill form
        (wrapper.vm as any).email = 'success@dpay.com';
        (wrapper.vm as any).password = 'password123';

        // Submit
        await wrapper.find('form').trigger('submit');

        // Wait for all async operations to finish
        await flushPromises();

        // Verify side effects
        expect(mockLoginApi).toHaveBeenCalledWith({ email: 'success@dpay.com', password: 'password123' });
        expect(store.isAuthenticated).toBe(true);
        expect(store.user?.email).toBe('success@dpay.com');
        expect(mockNavigateTo).toHaveBeenCalledWith('/', { replace: true });
    });

    it('failed login shows error and does not navigate', async () => {
        const store = useAuthStore();

        // Mock API error
        mockLoginApi.mockRejectedValueOnce({
            data: { message: 'Invalid credentials' }
        });

        const wrapper = await mountSuspended(LoginForm);

        (wrapper.vm as any).email = 'fail@dpay.com';
        (wrapper.vm as any).password = 'wrong';
        await wrapper.find('form').trigger('submit');
        await flushPromises();

        expect(store.isAuthenticated).toBe(false);
        expect((wrapper.vm as any).apiError).toBe('Invalid credentials');
        expect(mockNavigateTo).not.toHaveBeenCalled();
    });
});
