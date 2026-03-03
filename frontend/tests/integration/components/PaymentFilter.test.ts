import { describe, it, expect } from 'vitest';
import { mountSuspended } from '@nuxt/test-utils/runtime';
import PaymentFilter from '../../../components/PaymentFilter.vue';

describe('PaymentFilter.vue', () => {
    it('renders all filter options', async () => {
        const wrapper = await mountSuspended(PaymentFilter, {
            props: { modelValue: 'All' }
        });

        const buttons = wrapper.findAll('.filter-btn');
        expect(buttons.length).toBe(4);
        expect(buttons[0].text()).toBe('All');
        expect(buttons[1].text()).toBe('Completed');
        expect(buttons[2].text()).toBe('Processing');
        expect(buttons[3].text()).toBe('Failed');
    });

    it('applies active class to selected filter', async () => {
        const wrapper = await mountSuspended(PaymentFilter, {
            props: { modelValue: 'processing' }
        });

        const buttons = wrapper.findAll('.filter-btn');
        expect(buttons[0].classes()).not.toContain('active');
        expect(buttons[2].classes()).toContain('active'); // Processing
    });

    it('emits update:modelValue on click', async () => {
        const wrapper = await mountSuspended(PaymentFilter, {
            props: { modelValue: 'All' }
        });

        const buttons = wrapper.findAll('.filter-btn');
        await buttons[1].trigger('click'); // Click Completed

        expect(wrapper.emitted('update:modelValue')).toBeTruthy();
        expect(wrapper.emitted('update:modelValue')![0]).toEqual(['completed']);
    });
});
