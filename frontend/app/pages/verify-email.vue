<script setup lang="ts">
import { toast } from "vue-sonner";

definePageMeta({ layout: false, middleware: "guest" });

const authAPI = useAuthAPI();
const { setUserFromAuthResponse } = useAuth();
const router = useRouter();
const route = useRoute();

const email = computed(() => (route.query.email as string) || "");
const code = ref("");
const isVerifying = ref(false);
const isResending = ref(false);
const errorMessage = ref("");
const resendMessage = ref("");

onMounted(async () => {
  if (email.value) {
    await resend();
  }
});

async function submitCode() {
  if (code.value.trim().length !== 6) {
    errorMessage.value = "6 оронтой кодыг оруулна уу";
    return;
  }
  errorMessage.value = "";
  resendMessage.value = "";
  isVerifying.value = true;
  try {
    const response = await authAPI.verifyEmail(email.value, code.value.trim());
    await setUserFromAuthResponse(response);
    toast.success("И-мэйл амжилттай баталгаажлаа!");
    await router.push("/");
  } catch (e: any) {
    const msg = e?.data?.message || e?.message || "Баталгаажуулахад алдаа гарлаа";
    errorMessage.value = msg;
    toast.error(msg);
  } finally {
    isVerifying.value = false;
  }
}

async function resend() {
  if (!email.value) {
    toast.error("И-мэйл хаяг олдсонгүй");
    return;
  }
  resendMessage.value = "";
  isResending.value = true;
  try {
    const res = await authAPI.resendVerification(email.value);
    resendMessage.value = res.message;
    toast.success(res.message);
  } catch (e: any) {
    const msg = e?.data?.message || "Код дахин илгээхэд алдаа гарлаа";
    resendMessage.value = msg;
    toast.error(msg);
  } finally {
    isResending.value = false;
  }
}
</script>

<template>
  <AuthShell
    title="И-мэйл баталгаажуулалт"
    description="Бүртгэлийн аюулгүй байдлыг хангахын тулд и-мэйл хаягаа баталгаажуулна уу."
    :highlights="[
      { label: 'Аюулгүй', value: 'И-мэйл баталгаажуулалт' },
      { label: 'Хурдан', value: '15 минутын дотор' },
      { label: 'Хялбар', value: '6 оронтой код' },
    ]"
  >
    <AuthCard
      title="И-мэйл баталгаажуулах"
      description="И-мэйл хаягт илгээсэн 6 оронтой кодыг оруулна уу."
      footer-text="Бүртгэлтэй юу?"
      footer-link-text="Нэвтрэх"
      footer-link-to="/login"
    >
      <form class="space-y-5" @submit.prevent="submitCode">
        <div
          v-if="email"
          class="rounded-2xl border border-border bg-muted/40 px-4 py-3 text-sm text-muted-foreground"
        >
          Код илгээгдсэн хаяг: <span class="font-medium text-foreground">{{ email }}</span>
        </div>

        <div
          v-if="errorMessage"
          class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700"
        >
          {{ errorMessage }}
        </div>

        <div
          v-if="resendMessage && !errorMessage"
          class="rounded-2xl border border-green-200 bg-green-50 px-4 py-3 text-sm text-green-700"
        >
          {{ resendMessage }}
        </div>

        <div class="space-y-2">
          <Label for="code">Баталгаажуулах код</Label>
          <Input
            id="code"
            v-model="code"
            type="text"
            inputmode="numeric"
            maxlength="6"
            placeholder="000000"
            class="text-center text-2xl tracking-[0.5em] font-bold"
          />
          <p class="text-xs text-muted-foreground text-center">
            Код бүртгүүлэх үед и-мэйл хаягт илгээгдсэн байна
          </p>
        </div>

        <Button class="w-full" type="submit" :disabled="isVerifying || !code">
          {{ isVerifying ? "Баталгаажуулж байна..." : "Баталгаажуулах" }}
        </Button>

        <div class="text-center text-sm text-muted-foreground">
          Код ирээгүй юу?
          <button
            type="button"
            class="ml-1 font-medium text-foreground hover:text-muted-foreground transition disabled:opacity-50"
            :disabled="isResending"
            @click="resend"
          >
            {{ isResending ? "Илгээж байна..." : "Дахин илгээх" }}
          </button>
        </div>
      </form>
    </AuthCard>
  </AuthShell>
</template>
