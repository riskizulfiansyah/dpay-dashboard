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
          <span>{{ userInitials }}</span>
        </div>
        <div v-show="isExpanded" class="sidebar-user-info">
          <span class="sidebar-user-name">{{ userName }}</span>
          <span class="sidebar-user-role">{{ userRole }}</span>
        </div>
      </div>
      <button class="sidebar-logout" @click="handleLogout">
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
const { user, logout } = useAuth();

const userName = computed(() => user.value?.email || 'Guest');
const userRole = computed(() => {
  const role = user.value?.role || '';
  return role.charAt(0).toUpperCase() + role.slice(1);
});
const userInitials = computed(() => {
  if (!user.value?.email) return 'G';
  const name = user.value.email.split('@')[0];
  return name.slice(0, 2).toUpperCase();
});

const handleLogout = () => {
  logout();
};

onMounted(() => {
  checkScreenSize();
  window.addEventListener('resize', checkScreenSize);
});

onUnmounted(() => {
  window.removeEventListener('resize', checkScreenSize);
});
</script>
