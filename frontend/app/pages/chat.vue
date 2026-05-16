<script setup lang="ts">
import type { ChatMessage } from "../composables/api/core/chat";
import { toast } from "vue-sonner";
import { Send, Bot, User, Loader2 } from "lucide-vue-next";

definePageMeta({ middleware: "auth" });

const chatAPI = useChatAPI();

const messages = ref<ChatMessage[]>([
  {
    role: "assistant",
    content:
      "Сайн байна уу! Би таны ажлын байр хайхад туслах AI туслагч. Цалин, байршил, чиглэлийн талаар хэлбэл тохирох ажлыг олоход тусална. Юу хайж байна вэ?",
  },
]);
const input = ref("");
const isLoading = ref(false);
const container = ref<HTMLElement | null>(null);

async function sendMessage() {
  const text = input.value.trim();
  if (!text || isLoading.value) return;

  messages.value.push({ role: "user", content: text });
  input.value = "";
  isLoading.value = true;

  await nextTick();
  scrollToBottom();

  try {
    // Send only the last 10 messages as history (to keep tokens low)
    const history = messages.value.slice(0, -1).slice(-10);
    const { reply } = await chatAPI.send(text, history);
    messages.value.push({ role: "assistant", content: reply });
  } catch {
    toast.error("Хариу авахад алдаа гарлаа. Дахин оролдоно уу.");
    messages.value.pop();
    input.value = text;
  } finally {
    isLoading.value = false;
    await nextTick();
    scrollToBottom();
  }
}

function scrollToBottom() {
  if (container.value) {
    container.value.scrollTop = container.value.scrollHeight;
  }
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === "Enter" && !e.shiftKey) {
    e.preventDefault();
    sendMessage();
  }
}
</script>

<template>
  <div
    class="flex h-[calc(100vh-8rem)] flex-col gap-0 overflow-hidden rounded-3xl border border-border bg-card"
  >
    <!-- Header -->
    <div class="flex items-center gap-3 border-b border-border px-6 py-4">
      <div
        class="flex h-10 w-10 items-center justify-center rounded-2xl bg-primary text-primary-foreground"
      >
        <Bot class="h-5 w-5" />
      </div>
      <div>
        <p class="font-semibold">Ажлын байр хайх туслагч</p>
        <p class="text-xs text-muted-foreground">
          AI-д суурилсан ажлын байрны зөвлөгөө
        </p>
      </div>
    </div>

    <!-- Messages -->
    <div ref="container" class="flex-1 space-y-4 overflow-y-auto px-6 py-4">
      <div
        v-for="(msg, i) in messages"
        :key="i"
        :class="[
          'flex gap-3',
          msg.role === 'user' ? 'flex-row-reverse' : 'flex-row',
        ]"
      >
        <div
          :class="[
            'flex h-8 w-8 shrink-0 items-center justify-center rounded-2xl text-xs font-medium',
            msg.role === 'user'
              ? 'bg-primary text-primary-foreground'
              : 'bg-muted text-muted-foreground',
          ]"
        >
          <User v-if="msg.role === 'user'" class="h-4 w-4" />
          <Bot v-else class="h-4 w-4" />
        </div>
        <div
          :class="[
            'max-w-[80%] rounded-2xl px-4 py-3 text-sm leading-6',
            msg.role === 'user'
              ? 'bg-primary text-primary-foreground'
              : 'bg-muted text-foreground',
          ]"
        >
          {{ msg.content }}
        </div>
      </div>

      <div v-if="isLoading" class="flex gap-3">
        <div
          class="flex h-8 w-8 shrink-0 items-center justify-center rounded-2xl bg-muted text-muted-foreground"
        >
          <Bot class="h-4 w-4" />
        </div>
        <div class="rounded-2xl bg-muted px-4 py-3">
          <Loader2 class="h-4 w-4 animate-spin text-muted-foreground" />
        </div>
      </div>
    </div>

    <!-- Input -->
    <div class="border-t border-border px-6 py-4">
      <div class="flex gap-3">
        <textarea
          v-model="input"
          rows="1"
          placeholder="Ажил хайлтынхаа талаар бичнэ үү..."
          class="flex-1 resize-none rounded-2xl border border-input bg-background px-4 py-2.5 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          @keydown="onKeydown"
        />
        <Button
          class="h-10 w-10 shrink-0 rounded-2xl p-0"
          :disabled="!input.trim() || isLoading"
          @click="sendMessage"
        >
          <Send class="h-4 w-4" />
        </Button>
      </div>
      <p class="mt-2 text-xs text-muted-foreground">
        Enter дарж илгээх · Shift+Enter шинэ мөр
      </p>
    </div>
  </div>
</template>
