<template>
  <div class="input-wrapper">
    <label v-if="label" :for="id" class="input-label">
      {{ label }}
    </label>
    <div class="relative">
      <div v-if="showIcon" class="input-icon">
        <slot name="icon"></slot>
      </div>
      <input
        :id="id"
        v-model="modelValue"
        :type="type"
        :placeholder="placeholder"
        :disabled="disabled"
        :class="[
          'input-field',
          showIcon ? 'input-field-with-icon' : '',
          error ? 'input-field-error' : ''
        ]"
        @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      />
    </div>
    <p v-if="error" class="input-error-message">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
interface Props {
  modelValue?: string;
  type?: string;
  label?: string;
  placeholder?: string;
  disabled?: boolean;
  error?: string;
  showIcon?: boolean;
  id?: string;
}

withDefaults(defineProps<Props>(), {
  modelValue: '',
  type: 'text',
  placeholder: '',
  disabled: false,
  showIcon: false,
  id: () => `input-${Math.random().toString(36).substr(2, 9)}`,
});

const modelValue = defineModel<string>();

defineEmits<{
  'update:modelValue': [value: string];
}>();
</script>
