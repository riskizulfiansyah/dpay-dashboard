<template>
  <div class="login-page">
    <div class="login-page-bg-shield">
      <ShieldIcon :size="300" class="shield-icon-bg" />
    </div>
    
    <div class="login-page-container">
      <div class="login-page-card">
        <div class="login-page-header">
          <div class="login-page-icon">
            <ShieldIcon :size="48" class="shield-icon-color" />
          </div>
          <h1 class="login-page-title">DPay System</h1>
          <p class="login-page-subtitle">
            Internal Payment Dashboard
          </p>
        </div>
        
        <LoginForm />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
});
</script>
