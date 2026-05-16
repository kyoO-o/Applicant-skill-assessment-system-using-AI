<script setup lang="ts">
const integrationsAPI = useIntegrationsAPI();
const route = useRoute();

const gcalConnected = ref(false);
const gcalLoading = ref(true);
const gcalWorking = ref(false);
const gcalMessage = ref("");
const gcalError = ref("");

onMounted(async () => {
  try {
    const status = await integrationsAPI.googleCalendarStatus();
    gcalConnected.value = status.connected;
  } catch {
    // not configured or error - treat as not connected
  } finally {
    gcalLoading.value = false;
  }

  if (route.query.gcal === "connected") {
    gcalConnected.value = true;
    gcalMessage.value = "Google Calendar амжилттай холбогдлоо!";
  } else if (route.query.gcal === "error") {
    gcalError.value =
      "Google Calendar холбоход алдаа гарлаа. Дахин оролдоно уу.";
  }
});

function connectGoogleCalendar() {
  window.location.href = integrationsAPI.googleCalendarConnectURL();
}

async function disconnectGoogleCalendar() {
  gcalWorking.value = true;
  gcalMessage.value = "";
  gcalError.value = "";
  try {
    await integrationsAPI.googleCalendarDisconnect();
    gcalConnected.value = false;
    gcalMessage.value = "Google Calendar холболт цуцлагдлаа.";
  } catch {
    gcalError.value = "Алдаа гарлаа. Дахин оролдоно уу.";
  } finally {
    gcalWorking.value = false;
  }
}
</script>

<template>
  <div class="space-y-6">
    <section class="rounded-3xl border border-border bg-card px-6 py-6">
      <p class="text-sm font-medium text-muted-foreground">Тохиргоо</p>
      <h1 class="mt-2 text-3xl font-semibold tracking-tight">
        Ажлын орчны тохиргоо
      </h1>
      <p class="mt-2 max-w-2xl text-sm leading-6 text-muted-foreground">
        Бүртгэлийн болон дэлгэцийн тохиргоог энд хийнэ үү.
      </p>
    </section>

    <section class="grid gap-4 lg:grid-cols-2">
      <!-- Google Calendar integration -->
      <Card class="rounded-3xl border-border shadow-none">
        <CardHeader>
          <div class="flex items-center gap-3">
            <div
              class="flex h-9 w-9 items-center justify-center rounded-xl bg-blue-50"
            >
              <svg
                class="h-5 w-5"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <rect
                  x="3"
                  y="4"
                  width="18"
                  height="17"
                  rx="2"
                  stroke="#4285F4"
                  stroke-width="1.8"
                  fill="white"
                />
                <path d="M3 9h18" stroke="#4285F4" stroke-width="1.8" />
                <path
                  d="M8 2v4M16 2v4"
                  stroke="#4285F4"
                  stroke-width="1.8"
                  stroke-linecap="round"
                />
                <rect
                  x="7"
                  y="12"
                  width="4"
                  height="3"
                  rx="0.5"
                  fill="#4285F4"
                />
              </svg>
            </div>
            <div>
              <CardTitle>Google Calendar</CardTitle>
              <CardDescription class="mt-0.5">
                Ярилцлагийн товыг автоматаар календарт нэмэх
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent class="space-y-3">
          <div
            v-if="gcalMessage"
            class="rounded-2xl border border-green-200 bg-green-50 px-4 py-3 text-sm text-green-700"
          >
            {{ gcalMessage }}
          </div>
          <div
            v-if="gcalError"
            class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700"
          >
            {{ gcalError }}
          </div>

          <div
            v-if="gcalLoading"
            class="flex items-center gap-2 text-sm text-muted-foreground"
          >
            <Spinner class="h-4 w-4" />
            <span>Шалгаж байна...</span>
          </div>

          <template v-else>
            <div
              class="flex items-center justify-between rounded-2xl border border-border bg-muted/30 px-4 py-3"
            >
              <div class="flex items-center gap-2">
                <div
                  class="h-2 w-2 rounded-full"
                  :class="gcalConnected ? 'bg-green-500' : 'bg-zinc-400'"
                />
                <span class="text-sm font-medium">
                  {{ gcalConnected ? "Холбогдсон" : "Холбогдоогүй" }}
                </span>
              </div>
              <Button
                size="sm"
                :variant="gcalConnected ? 'outline' : 'default'"
                :disabled="gcalWorking"
                @click="
                  gcalConnected
                    ? disconnectGoogleCalendar()
                    : connectGoogleCalendar()
                "
              >
                {{
                  gcalWorking
                    ? "..."
                    : gcalConnected
                      ? "Холболт цуцлах"
                      : "Холбох"
                }}
              </Button>
            </div>

            <p class="text-xs text-muted-foreground">
              Холбосны дараа ярилцлага товлох үед Google Calendar-т автоматаар
              нэмэгдэнэ.
            </p>
          </template>
        </CardContent>
      </Card>

      <!-- Notifications placeholder -->
      <Card class="rounded-3xl border-border shadow-none">
        <CardHeader>
          <CardTitle>Мэдэгдэл</CardTitle>
          <CardDescription
            >Ярилцлага, шинэ горилогч болон сануулгын тохиргоо.</CardDescription
          >
        </CardHeader>
        <CardContent class="space-y-3">
          <div
            class="rounded-2xl border border-border bg-muted/30 p-4 text-sm text-muted-foreground"
          >
            Мэдэгдлийн тохиргоо удахгүй нэмэгдэнэ.
          </div>
        </CardContent>
      </Card>

      <!-- Appearance placeholder -->
      <Card class="rounded-3xl border-border shadow-none">
        <CardHeader>
          <CardTitle>Харагдац</CardTitle>
          <CardDescription
            >Дэлгэцийн загвар болон өнгөний тохиргоо.</CardDescription
          >
        </CardHeader>
        <CardContent class="space-y-3">
          <div
            class="rounded-2xl border border-border bg-muted/30 p-4 text-sm text-muted-foreground"
          >
            Харагдацын тохиргоо удахгүй нэмэгдэнэ.
          </div>
        </CardContent>
      </Card>
    </section>
  </div>
</template>
