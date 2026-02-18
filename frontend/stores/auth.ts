import { defineStore } from 'pinia';
import type { LoginRequest, LoginResponse } from '~/utils/api';

interface User {
  email: string;
  role: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  error: string | null;
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    token: null,
    isLoading: false,
    error: null,
  }),

  getters: {
    isAuthenticated: (state) => !!(state.token && state.user?.email),
    userRole: (state) => state.user?.role || null,
    userEmail: (state) => state.user?.email || null,
  },

  actions: {
    initialize() {
      const tokenCookie = useCookie<string>('auth_token', {
        maxAge: 60 * 60 * 24,
        default: () => '',
        path: '/',
      });
      const userCookie = useCookie<string>('auth_user', {
        maxAge: 60 * 60 * 24,
        default: () => '',
        path: '/',
      });

      const token = tokenCookie.value;
      const userStr = userCookie.value;

      if (token && userStr) {
        if (typeof userStr === 'object') {
          this.token = token;
          this.user = userStr as unknown as User;
        } else if (typeof userStr === 'string' && token !== '' && userStr !== '') {
          try {
            this.token = token;
            this.user = JSON.parse(userStr) as User;
          } catch {
            this.token = null;
            this.user = null;
          }
        }
      }
    },

    setAuth(auth: LoginResponse) {
      const user: User = {
        email: auth.email,
        role: auth.role,
      };

      this.user = user;
      this.token = auth.token;
      this.error = null;

      const tokenCookie = useCookie<string>('auth_token', {
        maxAge: 60 * 60 * 24,
        default: () => '',
        path: '/',
      });
      const userCookie = useCookie<string>('auth_user', {
        maxAge: 60 * 60 * 24,
        default: () => '',
        path: '/',
      });

      tokenCookie.value = auth.token;
      userCookie.value = JSON.stringify(user);
    },

    clearAuth() {
      this.user = null;
      this.token = null;
      this.error = null;
      this.isLoading = false;

      const tokenCookie = useCookie<string>('auth_token', { path: '/' });
      const userCookie = useCookie<string>('auth_user', { path: '/' });

      tokenCookie.value = '';
      userCookie.value = '';
    },

    setLoading(loading: boolean) {
      this.isLoading = loading;
    },

    setError(error: string | null) {
      this.error = error;
    },

    clearError() {
      this.error = null;
    },
  },
});
