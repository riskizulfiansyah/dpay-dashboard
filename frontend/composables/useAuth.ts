import { authApi, type LoginRequest } from '~/utils/api';
import { useAuthStore } from '~/stores/auth';

export function useAuth() {
  const authStore = useAuthStore();

  const login = async (credentials: LoginRequest) => {
    authStore.clearError();
    authStore.setLoading(true);

    try {
      const response = await authApi.login(credentials);
      authStore.setAuth(response);
      await navigateTo('/', { replace: true })
    } catch (error: unknown) {
      const err = error as { data?: { message?: string } };
      const message = err.data?.message || 'Login failed. Please try again.';
      authStore.setError(message);
      throw new Error(message);
    } finally {
      authStore.setLoading(false);
    }
  };

  const logout = () => {
    authStore.clearAuth();
    navigateTo('/login');
  };

  return {
    login,
    logout,
    user: computed(() => authStore.user),
    token: computed(() => authStore.token),
    isAuthenticated: computed(() => authStore.isAuthenticated),
    isLoading: computed(() => authStore.isLoading),
    error: computed(() => authStore.error),
    userRole: computed(() => authStore.userRole),
    clearError: authStore.clearError,
  };
}
