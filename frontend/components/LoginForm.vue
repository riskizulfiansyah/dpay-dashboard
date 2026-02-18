<template>
  <form class="login-form" @submit.prevent="handleSubmit">
    <InputField
      v-model="email"
      type="email"
      label="Email"
      placeholder="Enter your email"
      :error="errors.email"
      :show-icon="true"
    >
      <template #icon>
        <MailIcon :size="20" />
      </template>
    </InputField>
    
    <InputField
      v-model="password"
      type="password"
      label="Password"
      placeholder="Enter your password"
      :error="errors.password"
      :show-icon="true"
    >
      <template #icon>
        <LockIcon :size="20" />
      </template>
    </InputField>
    
    <button
      type="submit"
      class="login-form-button"
      :disabled="isLoading"
    >
      <span>{{ isLoading ? 'Signing in...' : 'Sign In to Dashboard' }}</span>
      <ArrowRightIcon v-if="!isLoading" :size="20" />
    </button>
  </form>
</template>

<script setup lang="ts">
const email = ref('');
const password = ref('');
const isLoading = ref(false);
const errors = reactive({
  email: '',
  password: '',
});

const handleSubmit = () => {
  errors.email = '';
  errors.password = '';
  
  if (!email.value) {
    errors.email = 'Email is required';
    return;
  }
  
  if (!password.value) {
    errors.password = 'Password is required';
    return;
  }
  
  isLoading.value = true;
  
  console.log('Login submitted:', { email: email.value, password: password.value });
};
</script>
