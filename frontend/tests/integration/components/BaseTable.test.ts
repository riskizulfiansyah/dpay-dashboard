import { describe, it, expect } from 'vitest';
import { mountSuspended } from '@nuxt/test-utils/runtime';
import BaseTable from '../../../components/BaseTable.vue';

describe('BaseTable.vue', () => {
    const columns = [
        { header: 'ID', key: 'id', sortable: true },
        { header: 'Name', key: 'name', sortable: false },
    ];

    const data = [
        { id: 1, name: 'Alice' },
        { id: 2, name: 'Bob' },
    ];

    it('renders correct headers and rows', async () => {
        const wrapper = await mountSuspended(BaseTable, {
            props: { columns, data }
        });

        const headers = wrapper.findAll('th');
        expect(headers.length).toBe(2);
        expect(headers[0].text()).toContain('ID');
        expect(headers[1].text()).toContain('Name');

        const rows = wrapper.findAll('tbody tr');
        expect(rows.length).toBe(2);
        expect(rows[0].text()).toContain('Alice');
        expect(rows[1].text()).toContain('Bob');
    });

    it('shows empty state when data is empty', async () => {
        const wrapper = await mountSuspended(BaseTable, {
            props: { columns, data: [] }
        });

        expect(wrapper.text()).toContain('No data found');
    });

    it('emits update:sort on sortable column click', async () => {
        const wrapper = await mountSuspended(BaseTable, {
            props: {
                columns,
                data,
                sort: { field: 'id', direction: 'desc' }
            }
        });

        const headers = wrapper.findAll('th');
        await headers[0].trigger('click'); // ID is sortable

        expect(wrapper.emitted('update:sort')).toBeTruthy();
        expect(wrapper.emitted('update:sort')![0]).toEqual([{ field: 'id', direction: 'asc' }]);

        await headers[1].trigger('click'); // Name is not sortable
        expect(wrapper.emitted('update:sort')?.length).toBe(1); // No new emit
    });

    it('renders pagination and emits page-change', async () => {
        const wrapper = await mountSuspended(BaseTable, {
            props: {
                columns,
                data,
                pagination: { page: 2, totalPages: 5 }
            }
        });

        expect(wrapper.text()).toContain('Page 2 of 5');

        const buttons = wrapper.findAll('.pagination-btn');
        expect(buttons.length).toBe(2);

        await buttons[0].trigger('click'); // Previous
        expect(wrapper.emitted('page-change')![0]).toEqual([1]);

        await buttons[1].trigger('click'); // Next
        expect(wrapper.emitted('page-change')![1]).toEqual([3]);
    });
});
