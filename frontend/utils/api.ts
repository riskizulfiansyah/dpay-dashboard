interface ApiConfig {
  baseURL: string;
}

const config: ApiConfig = {
  baseURL: '/api',
};

export const api = $fetch.create({
  baseURL: config.baseURL,
  credentials: 'include',
  onRequest({ options }) {
    const token = useCookie<string>('auth_token', {
      default: () => '',
      path: '/',
    });
    if (token.value && token.value !== '') {
      options.headers = new Headers(options.headers);
      options.headers.set('Authorization', `Bearer ${token.value}`);
    }
  },
  onResponseError({ response }) {
    if (response.status === 401) {
      const token = useCookie<string>('auth_token', { path: '/' });
      token.value = '';
      const router = useRouter();
      if (router.currentRoute.value.path !== '/login') {
        router.push('/login');
      }
    }
  },
});

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  email: string;
  role: string;
  token: string;
}

export const authApi = {
  login: (data: LoginRequest) => api<LoginResponse>('/auth/login', {
    method: 'POST',
    body: data,
  }).then((response) => {
    if (typeof response === 'string') {
      return JSON.parse(response) as LoginResponse;
    }
    return response;
  }),
};
