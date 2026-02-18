<template>
  <div class="stat-card">
    <div class="stat-card-icon" :class="iconClass">
      <slot name="icon"></slot>
    </div>
    <h3 class="stat-card-title">{{ title }}</h3>
    <div class="stat-card-value-row">
      <span class="stat-card-value">{{ value }}</span>
      <span v-if="percentage" class="stat-card-percentage" :class="percentageClass">
        {{ percentage }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  title: string;
  value: string;
  percentage?: string;
  variant?: 'default' | 'success' | 'danger';
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
});

const iconClass = computed(() => {
  return `stat-card-icon-${props.variant}`;
});

const percentageClass = computed(() => {
  const baseClass = 'stat-card-percentage-';
  if (props.percentage?.startsWith('-')) {
    return baseClass + 'negative';
  }
  return baseClass + props.variant;
});
</script>
