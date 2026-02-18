<template>
  <div class="payment-filter">
    <div class="filter-group">
      <button 
        v-for="status in statuses" 
        :key="status.value"
        class="filter-btn"
        :class="{ active: modelValue === status.value }"
        @click="$emit('update:modelValue', status.value)"
      >
        {{ status.label }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: string;
}>();

defineEmits<{
  (e: 'update:modelValue', value: string): void;
}>();

const statuses = [
  { label: 'All', value: 'All' },
  { label: 'Completed', value: 'completed' },
  { label: 'Processing', value: 'processing' },
  { label: 'Failed', value: 'failed' },
];
</script>

<style scoped>
.payment-filter {
  display: flex;
  align-items: center;
}

.filter-group {
  display: flex;
  background-color: #f3f4f6;
  padding: 4px;
  border-radius: 8px;
  gap: 4px;
}

.filter-btn {
  padding: 8px 16px;
  border: none;
  background: transparent;
  color: #6b7280;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.filter-btn:hover {
  color: #111827;
}

.filter-btn.active {
  background-color: white;
  color: #4f46e5;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}
</style>
