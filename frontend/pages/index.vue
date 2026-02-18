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
            <button class="dashboard-filter-btn">
              <FilterIcon :size="18" />
              <span>Filter</span>
            </button>
          </div>
        </header>

        <section class="dashboard-stats">
          <StatCard
            title="Total Payments"
            value="12,458"
            percentage="+12.5%"
            variant="default"
          >
            <template #icon>
              <BarChartIcon :size="24" />
            </template>
          </StatCard>

          <StatCard
            title="Success Rate"
            value="98.2%"
            percentage="+2.1%"
            variant="success"
          >
            <template #icon>
              <CheckCircleIcon :size="24" />
            </template>
          </StatCard>

          <StatCard
            title="Failed"
            value="223"
            percentage="-5.3%"
            variant="danger"
          >
            <template #icon>
              <AlertCircleIcon :size="24" />
            </template>
          </StatCard>
        </section>

        <section class="dashboard-table-section">
          <PaymentTable />
        </section>

        <section class="dashboard-chart-section">
          <VolumeChart />
        </section>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
const { isExpanded, checkScreenSize, toggle } = useSidebar();

onMounted(() => {
  checkScreenSize();
  window.addEventListener('resize', checkScreenSize);
});

onUnmounted(() => {
  window.removeEventListener('resize', checkScreenSize);
});
</script>
