<template>
  <div class="payment-table-card">
    <div class="payment-table-header">
      <h3 class="payment-table-title">Latest Payments</h3>
      <a href="#" class="payment-table-view-all">View All Transactions</a>
    </div>
    <div class="payment-table-wrapper">
      <table class="payment-table">
        <thead>
          <tr>
            <th>PAYMENT ID</th>
            <th>MERCHANT NAME</th>
            <th>DATE</th>
            <th>AMOUNT</th>
            <th>STATUS</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="payment in payments" :key="payment.id">
            <td class="payment-table-id">{{ payment.id }}</td>
            <td class="payment-table-merchant">{{ payment.merchant }}</td>
            <td class="payment-table-date">{{ new Date(payment.created_at).toLocaleDateString() }}</td>
            <td class="payment-table-amount">{{ payment.amount }}</td>
            <td>
              <span
                class="payment-table-status"
                :class="{
                  'payment-table-status-success': payment.status === 'completed',
                  'payment-table-status-failed': payment.status === 'failed',
                  'payment-table-status-processing': payment.status === 'processing'
                }"
              >
                <span class="payment-table-status-dot"></span>
                {{ payment.status }}
              </span>
            </td>
          </tr>
          <tr v-if="payments.length === 0">
            <td colspan="5" class="payment-table-empty">No payments found</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Payment {
  id: string;
  merchant: string;
  created_at: string;
  amount: string;
  status: string;
}

defineProps<{
  payments: Payment[];
}>();
</script>

<style scoped>
/* Add style for processing status if not exists, reusing existing styles logic */
.payment-table-status-processing {
  background-color: #eff6ff;
  color: #3b82f6;
}
.payment-table-status-processing .payment-table-status-dot {
  background-color: #3b82f6;
}

.payment-table-status-success {
    background-color: #ecfdf5;
    color: #10b981;
}
.payment-table-status-success .payment-table-status-dot {
    background-color: #10b981;
}

.payment-table-status-failed {
    background-color: #fef2f2;
    color: #ef4444;
}
.payment-table-status-failed .payment-table-status-dot {
    background-color: #ef4444;
}
.payment-table-empty {
    text-align: center;
    padding: 24px;
    color: #6b7280;
}
</style>
