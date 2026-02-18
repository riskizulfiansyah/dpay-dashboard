<template>
  <div class="volume-chart">
    <h3 class="volume-chart-title">Volume by Merchant</h3>
    <div class="volume-chart-list">
      <div v-for="merchant in merchants" :key="merchant.name" class="volume-chart-item">
        <div class="volume-chart-header">
          <span class="volume-chart-name">{{ merchant.name }}</span>
          <span class="volume-chart-value">{{ merchant.value }}</span>
        </div>
        <div class="volume-chart-bar">
          <div
            class="volume-chart-fill"
            :style="{ width: merchant.percentage + '%' }"
          ></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Merchant {
  name: string;
  value: string;
  percentage: number;
}

const merchants: Merchant[] = [
  { name: 'TechStore Inc', value: '$45,200', percentage: 85 },
  { name: 'Fashion Hub', value: '$32,500', percentage: 65 },
  { name: 'FoodMart', value: '$28,900', percentage: 55 },
  { name: 'BookWorld', value: '$21,300', percentage: 42 },
  { name: 'Gadget Zone', value: '$18,750', percentage: 35 },
];
</script>
