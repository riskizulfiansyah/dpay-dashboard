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
      :error="errors.password || apiError"
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
const errors = reactive({
  email: '',
  password: '',
});
const apiError = ref('');

const { login, isLoading } = useAuth();

const handleSubmit = async () => {
  errors.email = '';
  errors.password = '';
  apiError.value = '';
  
  if (!email.value) {
    errors.email = 'Email is required';
    return;
  }
  
  if (!password.value) {
    errors.password = 'Password is required';
    return;
  }
  
  try {
    await login({ email: email.value, password: password.value });
  } catch (error) {
    apiError.value = error instanceof Error ? error.message : 'Login failed';
  }
};
</script>
