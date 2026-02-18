<template>
  <div class="payment-table-card">
    <div v-if="title || showViewAll" class="payment-table-header">
      <h3 v-if="title" class="payment-table-title">{{ title }}</h3>
      <NuxtLink v-if="showViewAll" to="/payments" class="payment-table-view-all">View All Transactions</NuxtLink>
    </div>
    
    <BaseTable
      :columns="columns"
      :data="payments"
      :sort="sort"
      :pagination="pagination"
      @update:sort="handleSort"
      @page-change="handlePageChange"
    >
      <!-- Custom Id Cell -->
      <template #cell-id="{ value }">
        <span class="payment-table-id">{{ value }}</span>
      </template>

      <!-- Custom Merchant Cell -->
      <template #cell-merchant="{ value }">
        <span class="payment-table-merchant">{{ value }}</span>
      </template>

      <!-- Custom Date Cell -->
      <template #cell-created_at="{ value }">
        <span class="payment-table-date">{{ new Date(value).toLocaleDateString() }}</span>
      </template>

      <!-- Custom Amount Cell -->
      <template #cell-amount="{ value }">
        <span class="payment-table-amount">{{ value }}</span>
      </template>

      <!-- Custom Status Cell -->
      <template #cell-status="{ value }">
        <span
          class="payment-table-status"
          :class="{
            'payment-table-status-success': value === 'completed',
            'payment-table-status-failed': value === 'failed',
            'payment-table-status-processing': value === 'processing'
          }"
        >
          <span class="payment-table-status-dot"></span>
          {{ value }}
        </span>
      </template>

      <template #empty>
        No payments found
      </template>
    </BaseTable>
  </div>
</template>

<script setup lang="ts">
import BaseTable from './BaseTable.vue';

interface Payment {
  id: string;
  merchant: string;
  created_at: string;
  amount: string;
  status: string;
}

interface SortConfig {
  field: string;
  direction: 'asc' | 'desc';
}

interface PaginationConfig {
  page: number;
  totalPages: number;
  totalCount?: number;
}

const props = withDefaults(defineProps<{
  payments: Payment[];
  title?: string;
  showViewAll?: boolean;
  sort?: SortConfig;
  pagination?: PaginationConfig;
}>(), {
  title: 'Latest Payments',
  showViewAll: true,
  sort: undefined,
  pagination: undefined,
});

const emit = defineEmits<{
  (e: 'update:sort', sort: SortConfig): void;
  (e: 'page-change', page: number): void;
}>();

const columns = [
  { header: 'PAYMENT ID', key: 'id', sortable: true },
  { header: 'MERCHANT NAME', key: 'merchant', sortable: false },
  { header: 'DATE', key: 'created_at', sortable: true },
  { header: 'AMOUNT', key: 'amount', sortable: true },
  { header: 'STATUS', key: 'status', sortable: true },
];

const handleSort = (sort: SortConfig) => {
  emit('update:sort', sort);
};

const handlePageChange = (page: number) => {
  emit('page-change', page);
};
</script>

<style scoped>
/* Reuse existing styles for status badges */
.payment-table-status-processing {
  background-color: #eff6ff;
  color: #3b82f6;
}
.payment-table-status-processing .payment-table-status-dot {
  background-color: #3b82f6;
}

.payment-table-status-success {
    background-color: #ecfdf5;
    color: #10b981;
}
.payment-table-status-success .payment-table-status-dot {
    background-color: #10b981;
}

.payment-table-status-failed {
    background-color: #fef2f2;
    color: #ef4444;
}
.payment-table-status-failed .payment-table-status-dot {
    background-color: #ef4444;
}
</style>
