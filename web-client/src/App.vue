<script setup lang="ts">
import { ref } from "vue";
import Navbar from "@/components/UI/Navbar.vue";
import Cell from "@/components/Game/Cell/index.vue";
import type { ICell } from "@/components/Game/Cell/types.ts";

const cells = ref<ICell[]>([]);
// const currentColor = ref(Paint);

for (let i = 0; i < 24 * 12; i++) {
  cells.value.push({
    color: "white",
  });
}

const colorizeCell = (idx: number) => {
  cells.value[idx].color = "red";
};
</script>

<template>
  <div class="flex flex-col h-screen w-screen">
    <Navbar>
      <a href="#">Навбар</a>
    </Navbar>

    <div class="flex flex-grow">
      <!-- Панель с инструментами -->
      <div class="w-16"></div>

      <!-- Игровая сетка -->
      <div class="grid grid-rows-24 grid-cols-12 grid-flow-row h-full w-full">
        <Cell
          :data="cell"
          :key="idx"
          @click="colorizeCell(idx)"
          v-for="(cell, idx) in cells"
        />
      </div>
    </div>
  </div>
</template>

<style scoped></style>
