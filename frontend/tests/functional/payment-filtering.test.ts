import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mountSuspended, mockComponent } from '@nuxt/test-utils/runtime';
import { ref } from 'vue';
import PaymentFilter from '../../components/PaymentFilter.vue';

// For this flow we would ideally test the page (`index.vue` or `payments/index.vue`), 
// but Nuxt page testing with AsyncData can be complex to mock.
// Instead, let's create a functional test simulating the page's behavior
// wrapper component to ensure the v-model binds properly and resets things correctly.

const mockApi = vi.fn();
vi.mock('../../utils/api', () => ({
    api: (url: string, options: any) => mockApi(url, options)
}));

describe('Payment Filtering Flow', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('filter changes trigger reactive state updates as in the pages', async () => {
        // Replicating a simplified version of `payments/index.vue` logic
        const filterStatus = ref('All');
        const page = ref(1);

        // Simulated watch behavior from page
        const watchFilter = (newVal: string) => {
            filterStatus.value = newVal;
            page.value = 1; // reset page
        };

        const wrapper = await mountSuspended(PaymentFilter, {
            props: {
                modelValue: filterStatus.value,
                'onUpdate:modelValue': (e: string) => watchFilter(e)
            }
        });

        // Initial state
        expect(filterStatus.value).toBe('All');
        expect(page.value).toBe(1);

        // Change page manually (e.g. user clicked next)
        page.value = 3;

        // Click 'Failed' filter
        const buttons = wrapper.findAll('.filter-btn');
        await buttons[3].trigger('click'); // Failed

        // Check if the reactive state was updated correctly
        expect(filterStatus.value).toBe('failed');
        expect(page.value).toBe(1); // Page must be reset to 1
    });
});
