export default defineNuxtRouteMiddleware((to) => {


  const tokenCookie = useCookie<string>('auth_token', {
    default: () => '',
    path: '/',
  });
  const userCookie = useCookie<string>('auth_user', {
    default: () => '',
    path: '/',
  });

  const token = tokenCookie.value;
  const userVal = userCookie.value;

  let isAuthenticated = false;
  if (token && userVal) {
    if (typeof userVal === 'object') {
      isAuthenticated = !!((userVal as any).email);
    } else if (typeof userVal === 'string' && userVal !== '') {
      try {
        const user = JSON.parse(userVal);
        isAuthenticated = !!(user?.email);
      } catch {
        isAuthenticated = false;
      }
    }
  }

  if (!isAuthenticated && to.path !== '/login') {
    return navigateTo('/login');
  }

  if (isAuthenticated && to.path === '/login') {
    return navigateTo('/');
  }
});
