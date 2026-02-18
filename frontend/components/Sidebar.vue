<template>
  <aside
    class="sidebar"
    :class="[
      isExpanded ? 'sidebar-expanded' : 'sidebar-collapsed',
      { 'sidebar-overlay': isOverlay && !isExpanded }
    ]"
  >
    <div class="sidebar-header">
      <div class="sidebar-brand">
        <ShieldIcon :size="isExpanded ? 32 : 24" class="sidebar-brand-icon" />
        <span v-show="isExpanded" class="sidebar-brand-text">DPay</span>
      </div>
      <button 
        v-if="isExpanded" 
        class="sidebar-toggle sidebar-toggle-inside" 
        @click="toggleSidebar" 
        aria-label="Collapse sidebar"
      >
        <PanelLeftCloseIcon :size="20" />
      </button>
    </div>

    <nav class="sidebar-nav">
      <NuxtLink to="/" class="sidebar-nav-item sidebar-nav-item-active">
        <OverviewIcon :size="20" />
        <span v-show="isExpanded" class="sidebar-nav-text">Overview</span>
      </NuxtLink>
      <div class="sidebar-nav-item sidebar-nav-item-disabled">
        <PaymentIcon :size="20" />
        <span v-show="isExpanded" class="sidebar-nav-text">Payments</span>
      </div>
    </nav>

    <div class="sidebar-footer">
      <div class="sidebar-user">
        <div class="sidebar-user-avatar">
          <span>AR</span>
        </div>
        <div v-show="isExpanded" class="sidebar-user-info">
          <span class="sidebar-user-name">Alex Rivera</span>
          <span class="sidebar-user-role">Admin</span>
        </div>
      </div>
      <button class="sidebar-logout">
        <LogoutIcon :size="20" />
        <span v-show="isExpanded" class="sidebar-logout-text">Logout</span>
      </button>
    </div>
  </aside>

  <div
    v-if="isOverlay && !isExpanded"
    class="sidebar-backdrop"
    @click="closeSidebar"
  ></div>
</template>

<script setup lang="ts">
const { isExpanded, isOverlay, checkScreenSize, toggle: toggleSidebar, close: closeSidebar } = useSidebar();

onMounted(() => {
  checkScreenSize();
  window.addEventListener('resize', checkScreenSize);
});

onUnmounted(() => {
  window.removeEventListener('resize', checkScreenSize);
});
</script>
