import { describe, it, expect, vi } from 'vitest';
import { mockNuxtImport, mountSuspended, mockComponent } from '@nuxt/test-utils/runtime';
import LoginForm from '../../../components/LoginForm.vue';
import { ref } from 'vue';

// Mock the InputField component since it's used inside LoginForm
mockComponent('InputField', () => import('../../../components/InputField.vue'));

// Mock useAuth composable
const mockLogin = vi.fn();
const mockIsLoading = ref(false);

mockNuxtImport('useAuth', () => {
    return () => ({
        login: mockLogin,
        isLoading: mockIsLoading,
    });
});

describe('LoginForm.vue', () => {
    it('renders email and password inputs', async () => {
        const wrapper = await mountSuspended(LoginForm);

        // We expect 2 inputs (via InputField or natively)
        expect(wrapper.html()).toContain('type="email"');
        expect(wrapper.html()).toContain('type="password"');
    });

    it('shows validation errors when fields are empty', async () => {
        const wrapper = await mountSuspended(LoginForm);

        await wrapper.find('form').trigger('submit');

        // Check if error strings are assigned in component local state
        // We would see them rendered if InputField exposes :error prop
        expect((wrapper.vm as any).errors.email).toBe('Email is required');
    });

    it('calls login when form is valid', async () => {
        const wrapper = await mountSuspended(LoginForm);

        (wrapper.vm as any).email = 'test@dpay.com';
        (wrapper.vm as any).password = 'password123';

        await wrapper.find('form').trigger('submit');

        expect(mockLogin).toHaveBeenCalledWith({
            email: 'test@dpay.com',
            password: 'password123'
        });
    });
});
