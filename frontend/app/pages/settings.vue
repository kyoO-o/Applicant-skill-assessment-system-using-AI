<script setup lang="ts">
import {
  Mail,
  Lock,
  LogOut,
  Trash2,
  Check,
  X,
  ArrowLeft,
} from "lucide-vue-next";
import { toast } from "vue-sonner";

const router = useRouter();
const { user, me } = useAuth();
const authAPI = useAuthAPI();
const integrationsAPI = useIntegrationsAPI();
const route = useRoute();
const { disconnect: disconnectWS } = useNotifications();

// ── Google Calendar ────────────────────────────────────────────
const gcalConnected = ref(false);
const gcalLoading = ref(true);
const gcalWorking = ref(false);

const isLoading = useState<boolean>("auth-loading", () => false);
const initialized = useState<boolean>("auth-initialized", () => false);

onMounted(async () => {
  try {
    const status = await integrationsAPI.googleCalendarStatus();
    gcalConnected.value = status.connected;
  } catch {
    // not configured or error — treat as not connected
  } finally {
    gcalLoading.value = false;
  }

  if (route.query.gcal === "connected") {
    gcalConnected.value = true;
    toast.success("Google Calendar амжилттай холбогдлоо!");
  } else if (route.query.gcal === "error") {
    toast.error("Google Calendar холбоход алдаа гарлаа. Дахин оролдоно уу.");
  }
});

function connectGoogleCalendar() {
  window.location.href = integrationsAPI.googleCalendarConnectURL();
}

async function disconnectGoogleCalendar() {
  gcalWorking.value = true;
  try {
    await integrationsAPI.googleCalendarDisconnect();
    gcalConnected.value = false;
    toast.success("Google Calendar холболт цуцлагдлаа.");
  } catch {
    toast.error("Алдаа гарлаа. Дахин оролдоно уу.");
  } finally {
    gcalWorking.value = false;
  }
}

// ── Notifications (local UI state) ────────────────────────────
const notificationItems = reactive([
  {
    key: "result",
    title: "Ажилд орох хүсэлтийн хариу",
    enabled: true,
  },
  {
    key: "interview",
    title: "Ярилцлага товлогдох / өөрчлөгдөх",
    enabled: true,
  },
  {
    key: "task",
    title: "Шинэ даалгавар ирэх",
    enabled: true,
  },
]);

// ── Email change modal ─────────────────────────────────────────
const emailModalOpen = ref(false);
const emailStep = ref<"enter" | "verify" | "done">("enter");
const newEmail = ref("");
const emailPassword = ref("");
const emailCode = ref<string[]>(["", "", "", "", "", ""]);
const isSendingCode = ref(false);
const isVerifyingEmail = ref(false);
const emailError = ref("");
const resendIn = ref(0);
const codeInputs = ref<(HTMLInputElement | null)[]>([]);
let resendTimer: ReturnType<typeof setInterval> | null = null;

function openEmailModal() {
  emailModalOpen.value = true;
  emailStep.value = "enter";
  newEmail.value = "";
  emailPassword.value = "";
  emailCode.value = ["", "", "", "", "", ""];
  emailError.value = "";
  resendIn.value = 0;
}

function closeEmailModal() {
  emailModalOpen.value = false;
  if (resendTimer) {
    clearInterval(resendTimer);
    resendTimer = null;
  }
}

function isValidNewEmail(email: string) {
  return (
    /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email) && email !== user.value?.email
  );
}

async function sendEmailCode() {
  if (!isValidNewEmail(newEmail.value) || !emailPassword.value) return;
  isSendingCode.value = true;
  emailError.value = "";
  try {
    const res = await authAPI.initiateEmailChange(newEmail.value.trim());
    toast.success(res.message || "Код илгээгдлээ");
    emailStep.value = "verify";
    resendIn.value = 30;
    startResendTimer();
    await nextTick();
    codeInputs.value[0]?.focus();
  } catch (e: any) {
    emailError.value = e?.data?.message || "Код илгээхэд алдаа гарлаа";
  } finally {
    isSendingCode.value = false;
  }
}

function startResendTimer() {
  if (resendTimer) clearInterval(resendTimer);
  resendTimer = setInterval(() => {
    resendIn.value--;
    if (resendIn.value <= 0) {
      clearInterval(resendTimer!);
      resendTimer = null;
    }
  }, 1000);
}

function setCodeDigit(i: number, value: string) {
  const cleaned = value.replace(/\D/g, "").slice(0, 1);
  emailCode.value[i] = cleaned;
  if (cleaned && i < 5) {
    nextTick(() => codeInputs.value[i + 1]?.focus());
  }
}

function handleCodeKeydown(i: number, e: KeyboardEvent) {
  if (e.key === "Backspace" && !emailCode.value[i] && i > 0)
    codeInputs.value[i - 1]?.focus();
  if (e.key === "ArrowLeft" && i > 0) codeInputs.value[i - 1]?.focus();
  if (e.key === "ArrowRight" && i < 5) codeInputs.value[i + 1]?.focus();
}

function handleCodePaste(e: ClipboardEvent) {
  const text = (e.clipboardData?.getData("text") || "")
    .replace(/\D/g, "")
    .slice(0, 6);
  if (!text) return;
  e.preventDefault();
  const next = ["", "", "", "", "", ""];
  for (let i = 0; i < text.length; i++) next[i] = text.charAt(i);
  emailCode.value = next;
  nextTick(() => codeInputs.value[Math.min(text.length, 5)]?.focus());
}

async function verifyEmailCode() {
  const joined = emailCode.value.join("");
  if (joined.length !== 6) return;
  isVerifyingEmail.value = true;
  emailError.value = "";
  try {
    await authAPI.verifyEmailChange(joined);
    await me();
    emailStep.value = "done";
    setTimeout(() => closeEmailModal(), 1500);
  } catch (e: any) {
    emailError.value = e?.data?.message || "Баталгаажуулахад алдаа гарлаа";
  } finally {
    isVerifyingEmail.value = false;
  }
}

// ── Password change modal ──────────────────────────────────────
const passwordModalOpen = ref(false);
const passwordForm = reactive({ current: "", newPwd: "", confirm: "" });
const showPasswords = ref(false);
const isChangingPassword = ref(false);
const passwordDone = ref(false);
const passwordError = ref("");

const passwordChecks = computed(() => [
  { ok: passwordForm.newPwd.length >= 8, label: "Наад зах нь 8 тэмдэгт" },
  {
    ok: /[A-Z]/.test(passwordForm.newPwd) && /[a-z]/.test(passwordForm.newPwd),
    label: "Том, жижиг үсэг",
  },
  { ok: /\d/.test(passwordForm.newPwd), label: "Тоо агуулсан" },
  { ok: /[^A-Za-z0-9]/.test(passwordForm.newPwd), label: "Тусгай тэмдэгт" },
]);

const passwordScore = computed(
  () => passwordChecks.value.filter((c) => c.ok).length,
);
const passwordStrength = computed(
  () => ["Хэт сул", "Сул", "Дунд", "Хүчтэй", "Маш хүчтэй"][passwordScore.value],
);
const passwordStrengthClass = computed(() =>
  passwordScore.value >= 3
    ? "text-emerald-600 dark:text-emerald-400"
    : passwordScore.value === 2
      ? "text-yellow-600 dark:text-yellow-400"
      : "text-red-500",
);
const passwordBarClass = computed(() =>
  passwordScore.value >= 3
    ? "bg-emerald-500"
    : passwordScore.value === 2
      ? "bg-yellow-500"
      : "bg-red-500",
);
const passwordMatch = computed(
  () =>
    passwordForm.newPwd &&
    passwordForm.confirm &&
    passwordForm.newPwd === passwordForm.confirm,
);
const canSubmitPassword = computed(
  () =>
    passwordForm.current &&
    passwordScore.value >= 3 &&
    passwordMatch.value &&
    !isChangingPassword.value,
);

function openPasswordModal() {
  passwordModalOpen.value = true;
  passwordForm.current = "";
  passwordForm.newPwd = "";
  passwordForm.confirm = "";
  showPasswords.value = false;
  passwordDone.value = false;
  passwordError.value = "";
}

function closePasswordModal() {
  passwordModalOpen.value = false;
}

async function changePassword() {
  if (!canSubmitPassword.value) return;
  isChangingPassword.value = true;
  passwordError.value = "";
  try {
    const res = await authAPI.changePassword(
      passwordForm.current,
      passwordForm.newPwd,
    );
    toast.success(res.message || "Нууц үг шинэчлэгдлээ");
    passwordDone.value = true;
    setTimeout(() => closePasswordModal(), 1300);
  } catch (e: any) {
    passwordError.value = e?.data?.message || "Нууц үг солиход алдаа гарлаа";
  } finally {
    isChangingPassword.value = false;
  }
}

async function logout() {
  isLoading.value = true;
  try {
    await authAPI.logout();
    user.value = null;
  } finally {
    initialized.value = true;
    isLoading.value = false;
  }
}

async function handleLogout() {
  disconnectWS();
  await logout();
  await router.push("/login");
}
</script>

<template>
  <div>
    <!-- Page header -->
    <div class="mb-6">
      <h1 class="text-2xl font-semibold tracking-tight">Тохиргоо</h1>
      <p class="mt-1.5 text-sm text-muted-foreground">
        Холболт болон бүртгэлийн тохиргоог энд удирдана уу.
      </p>
    </div>

    <div class="flex max-w-2xl flex-col gap-4">
      <!-- Connected accounts -->
      <Card class="rounded-3xl border-border shadow-none gap-4">
        <CardHeader>
          <CardTitle>Холбогдсон бүртгэлүүд</CardTitle>
        </CardHeader>
        <CardContent>
          <div
            class="flex items-center gap-3.5 rounded-2xl border border-border bg-muted/20 p-4"
          >
            <!-- Google multicolor icon -->
            <div
              class="flex h-11 w-11 flex-shrink-0 items-center justify-center rounded-xl border border-border bg-white"
            >
              <svg class="h-5 w-5" viewBox="0 0 24 24" aria-hidden="true">
                <path
                  fill="#4285F4"
                  d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                />
                <path
                  fill="#34A853"
                  d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                />
                <path
                  fill="#FBBC05"
                  d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                />
                <path
                  fill="#EA4335"
                  d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                />
              </svg>
            </div>

            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <span class="text-sm font-semibold">Google Calendar</span>
                <span
                  v-if="gcalConnected"
                  class="inline-flex items-center gap-1.5 rounded-full bg-emerald-50 px-2 py-0.5 text-xs font-semibold text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-400"
                >
                  <span class="h-1.5 w-1.5 rounded-full bg-emerald-500" />
                  Холбогдсон
                </span>
              </div>
              <p class="mt-0.5 text-xs text-muted-foreground">
                {{
                  gcalConnected
                    ? "Ярилцлагийн товыг автоматаар нэмнэ"
                    : "Холбож ярилцлагийн товыг автоматаар нэмэх"
                }}
              </p>
            </div>

            <div v-if="gcalLoading">
              <Spinner class="h-4 w-4 text-muted-foreground" />
            </div>
            <template v-else>
              <Button
                v-if="gcalConnected"
                variant="outline"
                size="sm"
                :disabled="gcalWorking"
                class="border-destructive/30 text-destructive hover:bg-destructive/5 hover:text-destructive"
                @click="disconnectGoogleCalendar"
              >
                {{ gcalWorking ? "..." : "Салгах" }}
              </Button>
              <Button v-else size="sm" @click="connectGoogleCalendar">
                Холбох
              </Button>
            </template>
          </div>
        </CardContent>
      </Card>

      <!-- Notifications -->
      <!-- <Card class="rounded-3xl border-border shadow-none gap-4">
        <CardHeader>
          <CardTitle>Мэдэгдэл</CardTitle>
          <CardDescription
            >Ямар мэдэгдэл хүлээн авахыг тохируулна уу.</CardDescription
          >
        </CardHeader>
        <CardContent class="flex flex-col gap-2">
          <div
            v-for="notif in notificationItems"
            :key="notif.key"
            class="flex items-center gap-3 rounded-xl bg-muted/30 p-3"
          >
            <div class="flex-1">
              <p class="text-sm font-medium">{{ notif.title }}</p>
            </div>
            <Switch
              :checked="notif.enabled"
              @update:checked="(v: boolean) => (notif.enabled = v)"
            />
          </div>
        </CardContent>
      </Card> -->

      <!-- Account -->
      <Card class="rounded-3xl border-border shadow-none gap-4">
        <CardHeader>
          <CardTitle>Бүртгэл</CardTitle>
        </CardHeader>
        <CardContent class="flex flex-col gap-2">
          <!-- Email -->
          <div class="flex items-center gap-3 rounded-xl bg-muted/30 p-3">
            <Mail class="h-4 w-4 flex-shrink-0 text-muted-foreground" />
            <div class="min-w-0 flex-1">
              <p class="text-sm font-medium">Э-мэйл</p>
              <p class="truncate text-xs text-muted-foreground">
                {{ user?.email }}
              </p>
            </div>
            <Button variant="outline" size="sm" @click="openEmailModal"
              >Шинэчлэх</Button
            >
          </div>

          <!-- Password -->
          <div class="flex items-center gap-3 rounded-xl bg-muted/30 p-3">
            <Lock class="h-4 w-4 flex-shrink-0 text-muted-foreground" />
            <div class="flex-1">
              <p class="text-sm font-medium">Нууц үг</p>
              <p class="text-xs tracking-widest text-muted-foreground">
                ••••••••••
              </p>
            </div>
            <Button variant="outline" size="sm" @click="openPasswordModal"
              >Шинэчлэх</Button
            >
          </div>

          <!-- Sign out everywhere -->
          <div class="flex items-center gap-3 rounded-xl bg-muted/30 p-3">
            <LogOut class="h-4 w-4 flex-shrink-0 text-destructive" />
            <div class="flex-1">
              <p class="text-sm font-medium text-destructive">Гарах</p>
            </div>
            <Button
              variant="outline"
              size="sm"
              class="border-destructive/30 text-destructive hover:bg-destructive/5 hover:text-destructive"
              @click="handleLogout"
              >Гарах</Button
            >
          </div>

          <!-- Delete account -->
          <div class="flex items-center gap-3 rounded-xl bg-muted/30 p-3">
            <Trash2 class="h-4 w-4 flex-shrink-0 text-destructive" />
            <div class="flex-1">
              <p class="text-sm font-medium text-destructive">Бүртгэл устгах</p>
              <p class="text-xs text-muted-foreground">
                Энэ үйлдлийг буцаах боломжгүй
              </p>
            </div>
            <Button
              variant="outline"
              size="sm"
              class="border-destructive/30 text-destructive hover:bg-destructive/5 hover:text-destructive"
              >Устгах</Button
            >
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- ── Email change modal ────────────────────────────────── -->
    <Dialog v-model:open="emailModalOpen">
      <DialogContent class="max-w-[460px]">
        <DialogHeader>
          <DialogTitle>
            {{ emailStep === "done" ? "Э-мэйл шинэчлэгдлээ" : "Э-мэйл солих" }}
          </DialogTitle>
          <DialogDescription>
            <template v-if="emailStep === 'enter'">
              Шинэ хаягт баталгаажуулах 6 оронтой код илгээнэ.
            </template>
            <template v-else-if="emailStep === 'verify'">
              {{ newEmail }} хаягт код илгээгдлээ.
            </template>
            <template v-else>
              Таны нэвтрэх и-мэйл амжилттай шинэчлэгдлээ.
            </template>
          </DialogDescription>
        </DialogHeader>

        <!-- Step: enter -->
        <template v-if="emailStep === 'enter'">
          <div class="flex flex-col gap-4">
            <div class="space-y-1.5">
              <Label>Одоогийн и-мэйл</Label>
              <div
                class="rounded-xl border border-border bg-muted/30 px-3 py-2.5 text-sm text-muted-foreground"
              >
                {{ user?.email }}
              </div>
            </div>
            <div class="space-y-1.5">
              <Label for="new-email-input">Шинэ и-мэйл</Label>
              <Input
                id="new-email-input"
                v-model="newEmail"
                type="email"
                placeholder="шинэ@хаяг.com"
                autocomplete="email"
              />
            </div>
            <div class="space-y-1.5">
              <div class="flex items-center justify-between">
                <Label for="email-modal-pw">Одоогийн нууц үг</Label>
                <span class="text-xs text-muted-foreground"
                  >Аюулгүйн үүднээс</span
                >
              </div>
              <Input
                id="email-modal-pw"
                v-model="emailPassword"
                type="password"
                placeholder="••••••••"
                autocomplete="current-password"
              />
            </div>
            <div
              v-if="emailError"
              class="flex items-center gap-2 rounded-xl border border-destructive/20 bg-destructive/10 px-3 py-2 text-sm text-destructive"
            >
              <X class="h-3.5 w-3.5 flex-shrink-0" />
              {{ emailError }}
            </div>
          </div>
          <DialogFooter>
            <Button variant="ghost" @click="closeEmailModal">Цуцлах</Button>
            <Button
              :disabled="
                !isValidNewEmail(newEmail) || !emailPassword || isSendingCode
              "
              @click="sendEmailCode"
            >
              {{ isSendingCode ? "Илгээж байна..." : "Код илгээх" }}
            </Button>
          </DialogFooter>
        </template>

        <!-- Step: verify -->
        <template v-else-if="emailStep === 'verify'">
          <div class="flex flex-col gap-4">
            <div
              class="flex items-start gap-2.5 rounded-xl border border-blue-200 bg-blue-50 px-3.5 py-3 text-sm text-blue-700 dark:border-blue-800 dark:bg-blue-950/30 dark:text-blue-300"
            >
              <Mail class="mt-0.5 h-4 w-4 flex-shrink-0" />
              <span>
                <strong>{{ newEmail }}</strong> хаягт 6 оронтой код илгээгдлээ.
                10 минутын дотор хүчинтэй.
              </span>
            </div>

            <div class="space-y-2">
              <Label>Баталгаажуулах код</Label>
              <div class="flex gap-2" @paste="handleCodePaste">
                <input
                  v-for="(_, i) in emailCode"
                  :key="i"
                  :ref="
                    (el: any) => {
                      if (el) codeInputs[i] = el;
                    }
                  "
                  :value="emailCode[i]"
                  inputmode="numeric"
                  maxlength="1"
                  class="h-14 w-12 rounded-xl border border-border bg-background text-center font-mono text-xl font-semibold outline-none transition focus:border-primary focus:ring-1 focus:ring-primary"
                  :class="emailCode[i] ? 'border-primary' : ''"
                  @input="
                    setCodeDigit(i, ($event.target as HTMLInputElement).value)
                  "
                  @keydown="handleCodeKeydown(i, $event)"
                />
              </div>
            </div>

            <div
              v-if="emailError"
              class="flex items-center gap-2 rounded-xl border border-destructive/20 bg-destructive/10 px-3 py-2 text-sm text-destructive"
            >
              <X class="h-3.5 w-3.5 flex-shrink-0" />
              {{ emailError }}
            </div>

            <div class="flex items-center justify-between text-xs">
              <button
                class="flex items-center gap-1 text-muted-foreground transition hover:text-foreground"
                @click="
                  () => {
                    emailStep = 'enter';
                    emailCode = ['', '', '', '', '', ''];
                    emailError = '';
                  }
                "
              >
                <ArrowLeft class="h-3.5 w-3.5" />
                Э-мэйл хаяг өөрчлөх
              </button>
              <button
                :disabled="resendIn > 0"
                class="font-medium transition"
                :class="
                  resendIn > 0
                    ? 'cursor-default text-muted-foreground'
                    : 'text-primary hover:opacity-80'
                "
                @click="resendIn === 0 && sendEmailCode()"
              >
                {{
                  resendIn > 0
                    ? `${resendIn}с-ийн дараа дахин илгээх`
                    : "Дахин илгээх"
                }}
              </button>
            </div>
          </div>

          <DialogFooter>
            <Button variant="ghost" @click="closeEmailModal">Цуцлах</Button>
            <Button
              :disabled="emailCode.join('').length !== 6 || isVerifyingEmail"
              @click="verifyEmailCode"
            >
              {{
                isVerifyingEmail ? "Баталгаажуулж байна..." : "Баталгаажуулах"
              }}
            </Button>
          </DialogFooter>
        </template>

        <!-- Step: done -->
        <template v-else>
          <div class="flex flex-col items-center gap-3 py-6 text-center">
            <div
              class="flex h-16 w-16 items-center justify-center rounded-full bg-emerald-50 text-emerald-600 dark:bg-emerald-950/40 dark:text-emerald-400"
            >
              <Check class="h-7 w-7" />
            </div>
            <div>
              <p class="font-semibold">{{ newEmail }}</p>
              <p class="mt-1 text-sm text-muted-foreground">
                Одооноос эхлэн энэ хаягаар нэвтэрнэ үү.
              </p>
            </div>
          </div>
        </template>
      </DialogContent>
    </Dialog>

    <!-- ── Password change modal ─────────────────────────────── -->
    <Dialog v-model:open="passwordModalOpen">
      <DialogContent class="max-w-[460px]">
        <DialogHeader>
          <DialogTitle>
            {{ passwordDone ? "Нууц үг шинэчлэгдлээ" : "Нууц үг солих" }}
          </DialogTitle>
          <DialogDescription>
            {{
              passwordDone
                ? "Бусад хэрэгсэлээс аюулгүйн үүднээс гарсан байна."
                : "Одоо байгаа нууц үгийг оруулж баталгаажуулна уу."
            }}
          </DialogDescription>
        </DialogHeader>

        <!-- Done state -->
        <template v-if="passwordDone">
          <div class="flex flex-col items-center gap-3 py-6 text-center">
            <div
              class="flex h-16 w-16 items-center justify-center rounded-full bg-emerald-50 text-emerald-600 dark:bg-emerald-950/40 dark:text-emerald-400"
            >
              <Check class="h-7 w-7" />
            </div>
            <p class="text-sm text-muted-foreground">
              Бусад хэрэгсэлээс аюулгүйн үүднээс гарсан байна.
            </p>
          </div>
        </template>

        <!-- Form -->
        <template v-else>
          <div class="flex flex-col gap-4">
            <div class="space-y-1.5">
              <Label for="pw-current">Одоогийн нууц үг</Label>
              <Input
                id="pw-current"
                :type="showPasswords ? 'text' : 'password'"
                v-model="passwordForm.current"
                placeholder="••••••••"
                autocomplete="current-password"
              />
            </div>

            <div class="space-y-1.5">
              <Label for="pw-new">Шинэ нууц үг</Label>
              <Input
                id="pw-new"
                :type="showPasswords ? 'text' : 'password'"
                v-model="passwordForm.newPwd"
                placeholder="Наад зах нь 8 тэмдэгт"
                autocomplete="new-password"
              />
              <template v-if="passwordForm.newPwd">
                <div class="mt-2 flex gap-1">
                  <div
                    v-for="i in 4"
                    :key="i"
                    class="h-1 flex-1 rounded-full transition-colors"
                    :class="i <= passwordScore ? passwordBarClass : 'bg-border'"
                  />
                </div>
                <p
                  class="mt-1 text-xs font-semibold"
                  :class="passwordStrengthClass"
                >
                  {{ passwordStrength }}
                </p>
                <div class="mt-2 grid grid-cols-2 gap-1">
                  <div
                    v-for="check in passwordChecks"
                    :key="check.label"
                    class="flex items-center gap-1.5 text-xs"
                    :class="
                      check.ok
                        ? 'text-emerald-600 dark:text-emerald-400'
                        : 'text-muted-foreground'
                    "
                  >
                    <Check v-if="check.ok" class="h-3 w-3" />
                    <span
                      v-else
                      class="inline-block h-2 w-2 rounded-full border border-muted-foreground/40"
                    />
                    {{ check.label }}
                  </div>
                </div>
              </template>
            </div>

            <div class="space-y-1.5">
              <Label for="pw-confirm">Нууц үг давтах</Label>
              <Input
                id="pw-confirm"
                :type="showPasswords ? 'text' : 'password'"
                v-model="passwordForm.confirm"
                placeholder="Нууц үгийг дахин оруулна уу"
                autocomplete="new-password"
              />
              <div
                v-if="passwordForm.confirm && !passwordMatch"
                class="flex items-center gap-1.5 text-xs text-red-500"
              >
                <X class="h-3 w-3" />
                Нууц үг таарахгүй байна
              </div>
              <div
                v-else-if="passwordMatch"
                class="flex items-center gap-1.5 text-xs text-emerald-600 dark:text-emerald-400"
              >
                <Check class="h-3 w-3" />
                Нууц үг таарч байна
              </div>
            </div>

            <label
              class="flex cursor-pointer select-none items-center gap-2 text-xs text-muted-foreground"
            >
              <input
                type="checkbox"
                v-model="showPasswords"
                class="accent-primary"
              />
              Нууц үг харуулах
            </label>

            <div
              v-if="passwordError"
              class="flex items-center gap-2 rounded-xl border border-destructive/20 bg-destructive/10 px-3 py-2 text-sm text-destructive"
            >
              <X class="h-3.5 w-3.5 flex-shrink-0" />
              {{ passwordError }}
            </div>
          </div>

          <DialogFooter>
            <Button variant="ghost" @click="closePasswordModal">Цуцлах</Button>
            <Button :disabled="!canSubmitPassword" @click="changePassword">
              {{
                isChangingPassword ? "Шинэчлэж байна..." : "Нууц үг шинэчлэх"
              }}
            </Button>
          </DialogFooter>
        </template>
      </DialogContent>
    </Dialog>
  </div>
</template>
