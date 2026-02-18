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
            <td class="payment-table-merchant">{{ payment.merchantName }}</td>
            <td class="payment-table-date">{{ payment.date }}</td>
            <td class="payment-table-amount">{{ payment.amount }}</td>
            <td>
              <span
                class="payment-table-status"
                :class="payment.status === 'Success' ? 'payment-table-status-success' : 'payment-table-status-failed'"
              >
                <span class="payment-table-status-dot"></span>
                {{ payment.status }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Payment {
  id: string;
  merchantName: string;
  date: string;
  amount: string;
  status: 'Success' | 'Failed';
}

const payments: Payment[] = [
  { id: 'PAY-001', merchantName: 'TechStore Inc', date: 'Feb 18, 2026', amount: '$1,250.00', status: 'Success' },
  { id: 'PAY-002', merchantName: 'Fashion Hub', date: 'Feb 18, 2026', amount: '$890.50', status: 'Success' },
  { id: 'PAY-003', merchantName: 'FoodMart', date: 'Feb 17, 2026', amount: '$425.00', status: 'Failed' },
  { id: 'PAY-004', merchantName: 'BookWorld', date: 'Feb 17, 2026', amount: '$1,100.00', status: 'Success' },
  { id: 'PAY-005', merchantName: 'Gadget Zone', date: 'Feb 16, 2026', amount: '$2,340.00', status: 'Success' },
  { id: 'PAY-006', merchantName: 'Home Decor', date: 'Feb 16, 2026', amount: '$675.00', status: 'Failed' },
];
</script>
