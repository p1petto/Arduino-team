<script setup lang="ts">
import type { ICell } from './types.ts';

interface ICellProps {
    data: ICell
}

const props = defineProps<ICellProps>();
</script>

<template>
    <div class="bg-slate-100 border hover:bg-slate-200">{{ props.data.color }}</div>
</template>