import { test, expect } from '@playwright/test';

// To avoid SSR issues where `page.goto('/')` bypasses `page.route` mocks
// because Nuxt calls the API serverside (Nitro), we will load the dashboard
// by navigating through the client-side login flow.

test.describe('Dashboard E2E', () => {
    test.beforeEach(async ({ page }) => {
        // 1. Setup mock routes for all endpoints
        await page.route('**/api/dashboard/v1/auth/login', async route => {
            const json = { email: 'admin@dpay.com', role: 'admin', token: 'mock-jwt-token' };
            await route.fulfill({ json });
        });

        await page.route('**/api/dashboard/v1/payments/summary', async route => {
            const json = {
                total: 100,
                status_counts: [
                    { status: 'completed', count: 80 },
                    { status: 'failed', count: 10 },
                    { status: 'processing', count: 10 }
                ]
            };
            await route.fulfill({ json });
        });

        await page.route('**/api/dashboard/v1/payments*', async route => {
            const url = new URL(route.request().url());
            const statusParam = url.searchParams.get('status');

            let payments = [
                { id: 'PAY-1', merchant: 'Test Merchant', amount: '$50.00', status: 'completed', created_at: '2023-01-01T00:00:00Z' },
                { id: 'PAY-2', merchant: 'Test Merchant 2', amount: '$20.00', status: 'processing', created_at: '2023-01-02T00:00:00Z' }
            ];

            if (statusParam === 'processing') {
                payments = [payments[1]]; // Only processing
            }

            await route.fulfill({
                json: {
                    pagination: { page: 1, limit: 5, total_pages: 5, total_count: 25 },
                    payments
                }
            });
        });

        // Setup auth state by performing an actual UI login.
        // This is safe now because NUXT_E2E=true disables SSR for the dashboard,
        // and disables the Nitro proxy, allowing Playwright to intercept everything.
        await page.goto('/login');

        await page.locator('input[type="email"]').pressSequentially('admin@dpay.com', { delay: 10 });
        await page.locator('input[type="password"]').pressSequentially('password123', { delay: 10 });
        await page.getByRole('button', { name: 'Sign In to Dashboard' }).click();

        // Wait for page to load
        await expect(page.locator('h1').filter({ hasText: 'Overview Dashboard' })).toBeVisible();
    });

    test('dashboard shows stat cards and data', async ({ page }) => {
        // Verify stats from mocked API
        await expect(page.locator('text=100').first()).toBeVisible(); // Total
        await expect(page.locator('text=80').first()).toBeVisible();  // Success
        await expect(page.locator('text=10').nth(0)).toBeVisible();   // Failed
    });

    test('dashboard table renders mocked data', async ({ page }) => {
        // Table elements
        await expect(page.locator('text=PAY-1')).toBeVisible();
        await expect(page.locator('text=PAY-2')).toBeVisible();

        // Status badges
        await expect(page.locator('.payment-table-status-success')).toContainText('completed');
    });

    test('payment filter passes correct param', async ({ page }) => {
        const processingBtn = page.getByRole('button', { name: 'Processing' });

        // First, verify the button is visible before trying to click it
        await expect(processingBtn).toBeVisible();

        // Use promise.all to wait for API request triggered by click
        const [request] = await Promise.all([
            page.waitForRequest(req => req.url().includes('/api/dashboard/v1/payments') && !req.url().includes('/summary') && req.url().includes('status=processing')),
            processingBtn.click()
        ]);

        expect(request.url()).toContain('status=processing');

        // UI should update to show only PAY-2 based on our dynamic route mock
        await expect(page.locator('text=PAY-1')).not.toBeVisible();
        await expect(page.locator('text=PAY-2')).toBeVisible();
    });
});
