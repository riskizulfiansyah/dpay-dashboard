<template>
  <div class="base-table-container">
    <div class="base-table-wrapper">
      <table class="base-table">
        <thead>
          <tr>
            <th
              v-for="column in columns"
              :key="column.key"
              :class="{ sortable: column.sortable }"
              @click="handleSort(column)"
            >
              <div class="th-content">
                {{ column.header }}
                <component
                  :is="getSortIcon(column.key)"
                  v-if="getSortIcon(column.key)"
                  :size="16"
                />
              </div>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(row, index) in data" :key="index">
            <td v-for="column in columns" :key="column.key" :class="column.tdClass">
              <slot :name="'cell-' + column.key" :row="row" :value="row[column.key]">
                {{ row[column.key] }}
              </slot>
            </td>
          </tr>
          <tr v-if="data.length === 0">
            <td :colspan="columns.length" class="base-table-empty">
              <slot name="empty">No data found</slot>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div v-if="pagination && pagination.totalPages > 1" class="pagination-controls">
      <button
        class="pagination-btn"
        :disabled="pagination.page === 1"
        @click="changePage(pagination.page - 1)"
      >
        <ArrowLeftIcon :size="16" />
        Previous
      </button>
      <span class="pagination-info">
        Page {{ pagination.page }} of {{ pagination.totalPages }}
      </span>
      <button
        class="pagination-btn"
        :disabled="pagination.page === pagination.totalPages"
        @click="changePage(pagination.page + 1)"
      >
        Next
        <ArrowRightIcon :size="16" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import ArrowUpIcon from './icons/ArrowUpIcon.vue';
import ArrowDownIcon from './icons/ArrowDownIcon.vue';
import ArrowLeftIcon from './icons/ArrowLeftIcon.vue';
import ArrowRightIcon from './icons/ArrowRightIcon.vue';

interface Column {
  header: string;
  key: string;
  sortable?: boolean;
  tdClass?: string;
}

interface SortConfig {
  field: string;
  direction: 'asc' | 'desc';
}

interface PaginationConfig {
  page: number;
  totalPages: number;
  totalCount?: number;
}

const props = withDefaults(defineProps<{
  columns: Column[];
  data: any[];
  sort?: SortConfig;
  pagination?: PaginationConfig;
}>(), {
  sort: undefined,
  pagination: undefined,
});

const emit = defineEmits<{
  (e: 'update:sort', sort: SortConfig): void;
  (e: 'page-change', page: number): void;
}>();

const handleSort = (column: Column) => {
  if (!column.sortable) return;

  const currentField = props.sort?.field;
  const currentDirection = props.sort?.direction;

  let newDirection: 'asc' | 'desc' = 'desc';

  if (currentField === column.key) {
    newDirection = currentDirection === 'desc' ? 'asc' : 'desc';
  }

  emit('update:sort', { field: column.key, direction: newDirection });
};

const getSortIcon = (field: string) => {
  if (!props.sort || props.sort.field !== field) return null;
  return props.sort.direction === 'asc' ? ArrowUpIcon : ArrowDownIcon;
};

const changePage = (newPage: number) => {
  if (
    props.pagination &&
    newPage >= 1 &&
    newPage <= props.pagination.totalPages
  ) {
    emit('page-change', newPage);
  }
};
</script>

<style scoped>
.base-table-container {
  background-color: white;
  border-radius: 8px;
  overflow: hidden;
}

.base-table-wrapper {
  overflow-x: auto;
}

.base-table {
  width: 100%;
  border-collapse: collapse;
}

.base-table th {
  padding: 12px 24px;
  text-align: left;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  color: #6b7280;
  border-bottom: 1px solid #e5e7eb;
}

.base-table td {
  padding: 16px 24px;
  font-size: 14px;
  color: #111827;
  border-bottom: 1px solid #f3f4f6;
}

.base-table tbody tr:last-child td {
  border-bottom: none;
}

.base-table-empty {
  text-align: center;
  padding: 24px;
  color: #6b7280;
}

.th-content {
  display: flex;
  align-items: center;
  gap: 4px;
}

th.sortable {
  cursor: pointer;
}

th.sortable:hover {
  background-color: #f9fafb;
}

/* Pagination Styles */
.pagination-controls {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  padding: 16px; 
  border-top: 1px solid #e5e7eb;
}

.pagination-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  background-color: white;
  color: #374151;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.pagination-btn:hover:not(:disabled) {
  background-color: #f9fafb;
  border-color: #d1d5db;
}

.pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background-color: #f3f4f6;
}

.pagination-info {
  font-size: 14px;
  color: #6b7280;
}
</style>
