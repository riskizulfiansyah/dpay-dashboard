// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";

const apiUrl = process.env.NUXT_PUBLIC_API_URL || 'http://localhost:8080';

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  modules: ['@pinia/nuxt'],
  css: ['~/assets/css/main.css'],
  runtimeConfig: {
    public: {
      apiBaseUrl: apiUrl,
    },
  },
  vite: {
    plugins: [tailwindcss() as any],
  },
  components: [
    {
      path: '~/components',
      pathPrefix: false,
    },
  ],
  nitro: {
    routeRules: {
      '/api/**': {
        proxy: `${apiUrl}/**`
      }
    }
  }
})
