export default defineNitroPlugin((nitroApp) => {
    nitroApp.hooks.hook('request', (event) => {
        // Check if we are in an E2E environment
        if (process.env.NUXT_E2E !== 'true') return;

        const path = event.path;

        // Simulate Auth Login Success
        if (path.includes('/dashboard/v1/auth/login')) {
            // Since it's a hook before the proxy, we can't easily hijack the response in the raw 'request' hook.
            // Nuxt 3 prefers 'beforeResponse' or injecting a Nitro handler overlay for mocking.
        }
    });
});
