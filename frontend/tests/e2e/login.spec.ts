import { test, expect } from '@playwright/test';

test.describe('Login E2E', () => {
    test('unauthenticated user is redirected to /login', async ({ page }) => {
        // Navigate straight to dashboard
        await page.goto('/');

        // Should redirect
        await expect(page).toHaveURL(/.*\/login/);
        await expect(page.locator('h1')).toContainText('DPay System');
    });

    test('successful login redirects to dashboard', async ({ page }) => {
        // Mock ALL APIs needed during the redirect to the dashboard
        await Promise.all([
            page.route('**/api/dashboard/v1/auth/login', async route => {
                const json = { email: 'test@dpay.com', role: 'admin', token: 'fake-token' };
                await route.fulfill({ json });
            }),
            page.route('**/api/dashboard/v1/payments/summary', async route => {
                await route.fulfill({ json: { total: 0, status_counts: [] } });
            }),
            page.route('**/api/dashboard/v1/payments*', async route => {
                // Only mock if it's NOT the summary endpoint to avoid conflict
                if (!route.request().url().includes('/summary')) {
                    await route.fulfill({ json: { pagination: { page: 1, limit: 5, total_pages: 1, total_count: 0 }, payments: [] } });
                } else {
                    route.fallback();
                }
            })
        ]);

        await page.goto('/login');

        // Fill form using exact locators against the composed InputField 
        await page.locator('input[type="email"]').pressSequentially('test@dpay.com', { delay: 10 });
        await page.locator('input[type="password"]').pressSequentially('password123', { delay: 10 });

        // Submit
        await page.getByRole('button', { name: 'Sign In to Dashboard' }).click();

        // Wait for the full redirect and dashboard hydration
        await expect(page).toHaveURL(/\/?$/);
        await expect(page.locator('h1').filter({ hasText: 'Overview Dashboard' })).toBeVisible();
    });

    test('failed login shows error message', async ({ page }) => {
        await page.route('**/api/dashboard/v1/auth/login', async route => {
            // API utility maps standard errors automatically. 
            // Based on useAuth: err.data?.message || err.message
            const json = { message: 'Invalid credentials' };
            await route.fulfill({ status: 401, json });
        });

        await page.goto('/login');

        await page.locator('input[type="email"]').pressSequentially('test@dpay.com', { delay: 10 });
        await page.locator('input[type="password"]').pressSequentially('wrongpin', { delay: 10 });
        await page.getByRole('button', { name: 'Sign In to Dashboard' }).click();

        // Look for the specific error paragraph emitted by InputField
        await expect(page.locator('p.input-error-message').first()).toContainText('Invalid credentials');

        // URL shouldn't change
        await expect(page).toHaveURL(/.*\/login/);
    });
});
