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
            <h1 class="dashboard-header-title">Overview Dashboard</h1>
          </div>
          <div class="dashboard-header-actions">
            <PaymentFilter v-model="filterStatus" />
          </div>
        </header>

        <section class="dashboard-stats">
          <StatCard
            title="Total Payments"
            :value="formattedTotalPayments"
            variant="default"
          >
            <template #icon>
              <BarChartIcon :size="24" />
            </template>
          </StatCard>

          <StatCard
            title="Success"
            :value="formattedSuccessPayments"
            variant="success"
          >
            <template #icon>
              <CheckCircleIcon :size="24" />
            </template>
          </StatCard>

          <StatCard
            title="Failed"
            :value="formattedFailedPayments"
            variant="danger"
          >
            <template #icon>
              <AlertCircleIcon :size="24" />
            </template>
          </StatCard>
        </section>

        <section class="dashboard-table-section">
          <PaymentTable :payments="payments" />
        </section>

        <section class="dashboard-chart-section">
          <VolumeChart />
        </section>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
});

const { isExpanded, checkScreenSize, toggle } = useSidebar();
const authStore = useAuthStore();
const filterStatus = ref('All');

const { data: paymentsData, refresh: refreshPayments } = await useAsyncData('dashboard-payments', () => api<PaymentListResponse>('/dashboard/v1/payments', {
  query: {
    limit: 5,
    status: filterStatus.value === 'All' ? undefined : filterStatus.value
  }
}), {
  watch: [filterStatus]
});

const { data: summaryData, refresh: refreshSummary } = await useAsyncData('dashboard-summary', () => api<PaymentSummaryResponse>('/dashboard/v1/payments/summary'));


const payments = computed(() => {
  return paymentsData.value?.payments || [];
});

const formattedTotalPayments = computed(() => {
  return summaryData.value?.total?.toLocaleString() || '0';
});

const formattedFailedPayments = computed(() => {
  const failedCount = summaryData.value?.status_counts?.find(s => s.status === 'failed')?.count || 0;
  return failedCount.toLocaleString();
});

const formattedSuccessPayments = computed(() => {
  const completedCount = summaryData.value?.status_counts?.find(s => s.status === 'completed')?.count || 0;
  return completedCount.toLocaleString();
});

onMounted(() => {
  checkScreenSize();
  window.addEventListener('resize', checkScreenSize);
});

onUnmounted(() => {
  window.removeEventListener('resize', checkScreenSize);
});
</script>
