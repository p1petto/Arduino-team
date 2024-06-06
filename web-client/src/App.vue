<script setup lang="ts">
import { Cell, type ICell } from "@/components/Game/Cell/index";
import { IconButton, Navbar, Tooltip, Modal } from "@/components/UI/index";
import { get, post } from "@/api/core/index";

import {
  ArrowPathIcon,
  PaintBrushIcon,
  UserIcon,
  CheckCircleIcon,
  LinkIcon,
} from "@heroicons/vue/24/outline";
import { ref } from "vue";

const cells = ref<ICell[]>([]);
const pallete = ["red", "blue", "green", "white"];
const currentColor = ref(0);

for (let i = 0; i < 16 * 16; i++) {
  cells.value.push({});
}

const colorizeCell = (idx: number) => {
  cells.value[idx].color = pallete[currentColor.value];
};

const selectColor = (idx: number) => {
  currentColor.value = idx;
  colorPickerVisible.value = false;
};

const colorPickerVisible = ref(false);
const authPanelVisible = ref(false);
const gamesModalVisible = ref(false);
const username = ref("");
const token = ref("");

(async function () {
  const tokenCached = localStorage.getItem("token");

  if (tokenCached) {
    console.log(`Using token from data: ${tokenCached}`);
    token.value = tokenCached;
    return;
  }
})();

function getHeaders() {
  return {
    Authorization: `Bearer ${token.value}`,
    "Content-Type": "application/x-www-form-urlencoded",
  };
}

async function createUser() {
  const resp = await post(`/login/${username.value}`);

  if (resp.status !== 201) {
    console.error("Failed to login with user");
    return;
  }

  const newToken = await resp.text();
  localStorage.setItem("token", newToken);
  token.value = newToken;

  console.log(`Token for user ${username.value} created: ${newToken}`);
}

const currentGames = ref({});

async function updateGames() {
  const resp = await fetch("http://localhost:1090/rooms", {
    headers: getHeaders(),
  });
  currentGames.value = Object.values(await resp.json());
}

updateGames();

async function createGame() {
  const resp = await fetch("http://localhost:1090/rooms", {
    method: "POST",
    headers: getHeaders(),
    body: new URLSearchParams({
      name: "new-room",
      IP: "127.0.0.1",
    }),
  });
  console.log(resp);
}

async function connectGame(ID: string) {
  const _socket = new WebSocket(`ws://localhost:1090/ws/${ID}`);
  console.log(ID);
}
</script>

<template>
  <Modal v-if="gamesModalVisible">
    <h3 class="font-bold text-center">Подключиться к игре</h3>

    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Имя</th>
          <th>IPv4</th>
          <th>Статус</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="{ ID, name, IP, status } of currentGames">
          <td>{{ ID }}</td>
          <td>{{ name }}</td>
          <td>{{ IP }}</td>
          <td>{{ status }}</td>
          <td>
            <button @click="connectGame(ID)">Подключиться</button>
          </td>
        </tr>
      </tbody>
    </table>

    <div class="flex flex-row gap-2">
      <button @click="createGame()">Создать</button>
      <IconButton @click="updateGames()" :icon="ArrowPathIcon" />
    </div>
  </Modal>

  <div class="flex flex-col h-screen w-screen bg-[#fafafa]">
    <Navbar class="flex flex-row justify-between">
      <div class="m-auto font-bold flex gap-2">
        <div class="m-auto">Текущая игра</div>
        <IconButton
          @click="gamesModalVisible = !gamesModalVisible"
          :icon="LinkIcon"
        />
      </div>
      <IconButton
        @click="authPanelVisible = !authPanelVisible"
        :icon="UserIcon"
      />
      <Tooltip v-if="authPanelVisible" class="right-4 top-12 left-auto">
        <div class="flex flex-row gap-2">
          <input type="text" v-model="username" />
          <CheckCircleIcon v-if="token" class="size-6 m-auto text-green-500" />
        </div>
        <button @click="createUser()">Войти</button>
      </Tooltip>
    </Navbar>

    <div class="flex flex-grow">
      <!-- Панель с инструментами -->
      <div class="flex flex-col justify-between p-2 bg-white shadow">
        <div>
          <IconButton
            :style="{ borderBottom: `3px solid ${pallete[currentColor]}` }"
            :icon="PaintBrushIcon"
            @click="colorPickerVisible = !colorPickerVisible"
          />

          <Tooltip v-if="colorPickerVisible">
            <template #header>Выбрать цвет</template>
            <div class="grid grid-cols-4">
              <div
                class="rounded size-6 border transition-all hover:shadow hover:scale-110"
                :style="{ backgroundColor: color }"
                @click="selectColor(idx)"
                v-for="(color, idx) in pallete"
              ></div>
            </div>
          </Tooltip>
        </div>
      </div>

      <!-- Игровая сетка -->
      <div class="grid grid-cols-[repeat(16,_minmax(0,_1fr))] m-auto">
        <Cell
          :color="color"
          :key="idx"
          @click="colorizeCell(idx)"
          v-for="({ color }, idx) in cells"
        />
      </div>
    </div>
  </div>
</template>
