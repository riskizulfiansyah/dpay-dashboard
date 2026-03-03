import { describe, it, expect } from 'vitest';
import { mountSuspended } from '@nuxt/test-utils/runtime';
import StatCard from '../../../components/StatCard.vue';

describe('StatCard.vue', () => {
    it('renders title and value', async () => {
        const wrapper = await mountSuspended(StatCard, {
            props: {
                title: 'Total Revenue',
                value: '$5,000',
            }
        });

        expect(wrapper.text()).toContain('Total Revenue');
        expect(wrapper.text()).toContain('$5,000');
    });

    it('renders percentage with correct class', async () => {
        const wrapper = await mountSuspended(StatCard, {
            props: {
                title: 'Growth',
                value: '12%',
                percentage: '+2.5%',
                variant: 'success'
            }
        });

        expect(wrapper.text()).toContain('+2.5%');
        const percentageEl = wrapper.find('.stat-card-percentage');
        expect(percentageEl.exists()).toBe(true);
        expect(percentageEl.classes()).toContain('stat-card-percentage-success');
    });

    it('applies negative class when percentage starts with minus', async () => {
        const wrapper = await mountSuspended(StatCard, {
            props: {
                title: 'Drop',
                value: '10',
                percentage: '-1.5%',
                variant: 'danger'
            }
        });

        const percentageEl = wrapper.find('.stat-card-percentage');
        expect(percentageEl.classes()).toContain('stat-card-percentage-negative');
    });

    it('renders icon slot', async () => {
        const wrapper = await mountSuspended(StatCard, {
            props: { title: 'Test', value: '0' },
            slots: {
                icon: '<div class="test-icon">icon</div>'
            }
        });

        expect(wrapper.find('.test-icon').exists()).toBe(true);
    });
});
