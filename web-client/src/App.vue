<script setup lang="ts">
import { Client } from "@/api/core/client";
import { Cell, type ICell } from "@/components/Game/Cell/index";
import { IconButton, Modal, Navbar, Tooltip } from "@/components/UI/index";

import {
  ArrowPathIcon,
  LinkIcon,
  PaintBrushIcon,
  UserIcon,
  ExclamationCircleIcon,
  InformationCircleIcon,
} from "@heroicons/vue/24/outline";
import { onMounted, onUnmounted, ref } from "vue";

const cells = ref<ICell[]>([]);
const pallete = ["red", "blue", "green", "white", "brown"];
const currentColor = ref(0);

for (let i = 0; i < 16 * 16; i++) {
  cells.value.push({});
}
type ToastStatus = "error" | "info" | "ok";

const colorizeCell = (idx: number) => {
  const color2rgb: Record<string, number[]> = {
    red: [255, 0, 0],
    green: [0, 255, 0],
    blue: [0, 0, 255],
    white: [255, 255, 255],
    brown: [150, 75, 0]
  };

  const color = pallete[currentColor.value];

  if (!gameSocket.value) {
    return toast("Для начала нужно подключиться к игре!", "error");
  }

  if (gameSocket.value.readyState === WebSocket.CLOSED) {
    return toast("Соединение с сервером потеряно", "error");
  }

  gameSocket.value!!.send(
    JSON.stringify({
      type: "Input",
      color: color2rgb[color],
      X: idx % 16,
      Y: Math.floor(idx / 16),
    })
  );

  // cells.value[idx].color = color;
};

const currentGameID = ref("-");

const selectColor = (idx: number) => {
  currentColor.value = idx;
  colorPickerVisible.value = false;
};

const colorPickerVisible = ref(false);
const authPanelVisible = ref(false);
const gamesModalVisible = ref(false);
const username = ref(localStorage.getItem("username") ?? "");
const token = ref(localStorage.getItem("token") ?? "");
const client = new Client("http://localhost:1090/", token.value);

async function logout() {
  toast("Вы вышли из учетной записи", "info");
  token.value = "";
  username.value = "";
  localStorage.clear();
}

async function createUser() {
  const resp = await client.post(`/login/${username.value}`);

  if (resp.status === 409) {
    toast(`Пользователь с именем ${username.value} уже существует!`, "error");
    return;
  }

  if (!resp.ok) {
    console.error("Failed to login with user");
    toast(`Не удалось войти под именем ${username.value}`, "error");
    return;
  }

  token.value = await resp.text();
  localStorage.setItem("token", token.value);
  localStorage.setItem("username", username.value);
  client.token = token.value;

  toast(`Вы вошли под пользователем ${username.value}`, "info");
  console.log(`Token for user ${username.value} created: ${token.value}`);
}

const currentGames = ref({});

async function updateGames() {
  try {
    const resp = await client.get("/rooms");
    currentGames.value = Object.values(await resp.json());
  } catch (e) {
    return toast("Не удалось обновить информацию о доступных играх", "error");
  }
}

updateGames();

async function createGame() {
  const resp = await fetch(client.url("/rooms"), {
    method: "POST",
    headers: client.headers(),
    body: new URLSearchParams({
      name: "new-room",
      IP: "127.0.0.1",
    }),
  });

  if (!resp.ok) {
    return toast("не удалось создать игру!", "error");
  }

  return toast("Новая игра успешно создана!", "ok");
}

function toast(message: string, type: ToastStatus) {
  toasts.value.push({ message, type });
}

const gameSocket = ref<WebSocket | null>(null);

async function connectGame(ID: string) {
  gameSocket.value = new WebSocket(
    `ws://localhost:1090/ws/${ID}?token=${token.value}`
  );

  currentGameID.value = ID;
  toast(`вы успешно подключились к игре ${ID}`, "ok");
  console.log(`Connected to room ${ID}`);

  gameSocket.value.onmessage = (ev: MessageEvent) => {
    const { type, message } = JSON.parse(ev.data);
    if (type === "Error") {
      return toast(message, "error");
    }
    if (type === "Output") {
      const colorCodes = (message as number[][][]).flat();

      colorCodes.forEach((elem, idx) => {
        switch (elem) {
          case [255, 0, 0]:
            cells.value[idx].color = "red";
            break;
          case [0, 255, 0]:
            cells.value[idx].color = "green";
            break;
          case [0, 0, 255]:
            cells.value[idx].color = "blue";
            break;

          case [150, 75, 0]:
            cells.value[idx].color = "brown";
            break;

          default:
            cells.value[idx].color = "black";
            break;
        }
      });
    }
  };
}

const timer = ref<number | null>(null);

onMounted(() => {
  timer.value = setInterval(() => {
    toasts.value.shift();
  }, 3000);
});

onUnmounted(() => {
  if (timer.value) clearInterval(timer.value);
});

const toasts = ref<{ message: string; type: ToastStatus }[]>([]);
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
        <div class="m-auto">{{ currentGameID }}</div>
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
        <div v-if="token">
          <div class="font-bold">{{ username }}</div>
          <button @click="logout()">Выйти</button>
        </div>
        <div v-else class="flex flex-row gap-2">
          <input type="text" v-model="username" />
          <button @click="createUser()">Войти</button>
        </div>
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

      <div
        class="absolute flex flex-col gap-4 right-4 bottom-4 rounded-lg shadow-lg bg-white z-10"
      >
        <div
          v-for="(toast, idx) in toasts"
          :key="idx"
          class="p-4 flex flex-row gap-2"
        >
          <ExclamationCircleIcon
            v-if="toast.type === 'error'"
            class="size-6 text-red-500"
          />
          <InformationCircleIcon
            v-if="toast.type === 'info'"
            class="size-6 text-blue-500"
          />
          {{ toast.message }}
        </div>
      </div>
    </div>
  </div>
</template>
