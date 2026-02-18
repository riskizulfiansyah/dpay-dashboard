<template>
  <div class="dashboard-layout">
    <Sidebar />
    <main class="dashboard-main" :class="{ 'sidebar-collapsed': !isExpanded }">
      <div class="dashboard-content">
        <header class="dashboard-header">
          <div class="dashboard-header-left">
            <button 
              v-if="!isExpanded" 
              class="sidebar-toggle-btn" 
              @click="toggle"
              aria-label="Expand sidebar"
            >
              <PanelLeftOpenIcon :size="20" />
            </button>
            <h1 class="dashboard-header-title">Payments</h1>
          </div>
          <div class="dashboard-header-actions">
            <PaymentFilter v-model="filterStatus" />
          </div>
        </header>

        <section class="dashboard-table-section">
          <PaymentTable 
            :payments="payments" 
            title="" 
            :show-view-all="false"
            v-model:sort="sortConfig"
            :pagination="pagination"
            @page-change="changePage"
          />
        </section>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
// Arrow icons no longer needed here as they are in BaseTable
import { ref, computed } from 'vue';

definePageMeta({
  middleware: 'auth',
});

const { isExpanded, checkScreenSize, toggle } = useSidebar();
const filterStatus = ref('All');
const page = ref(1);
const limit = ref(10);
const sortConfig = ref<{ field: string; direction: 'asc' | 'desc' } | undefined>({
  field: 'created_at',
  direction: 'desc'
});

// Reset page when filter changes
watch(filterStatus, () => {
  page.value = 1;
});

const { data: paymentsData, refresh } = await useAsyncData('payments-list', () => {
  const query: any = {
    page: page.value,
    limit: limit.value,
    status: filterStatus.value === 'All' ? undefined : filterStatus.value,
  };

  if (sortConfig.value) {
    query.sort = `${sortConfig.value.direction === 'desc' ? '-' : ''}${sortConfig.value.field}`;
  }

  return api<PaymentListResponse>('/dashboard/v1/payments', { query });
}, {
  watch: [page, filterStatus, sortConfig]
});

const payments = computed(() => paymentsData.value?.payments || []);
const pagination = computed(() => ({
  page: page.value,
  totalPages: paymentsData.value?.pagination?.total_pages || 1,
  totalCount: paymentsData.value?.pagination?.total_count || 0,
}));

const changePage = (newPage: number) => {
  if (newPage >= 1 && newPage <= pagination.value.totalPages) {
    page.value = newPage;
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
};

onMounted(() => {
  checkScreenSize();
  window.addEventListener('resize', checkScreenSize);
});

onUnmounted(() => {
  window.removeEventListener('resize', checkScreenSize);
});
</script>
