<script setup lang="ts">
import { toast } from "vue-sonner";
import {
  ArrowRight,
  Bot,
  Calendar,
  ChartNoAxesColumnIncreasing,
  BriefcaseBusiness,
  ClipboardList,
  FileText,
  User,
  LayoutGrid,
  Settings,
  LogOut,
  PanelLeft,
} from "lucide-vue-next";

const { user, logout, isLoading } = useAuth();
const router = useRouter();
const route = useRoute();

const menuUser = [
  { title: "Нүүр хуудас", icon: ChartNoAxesColumnIncreasing, to: "/" },
  { title: "Ажлын байрууд", icon: BriefcaseBusiness, to: "/jobs" },
  { title: "Миний анкетууд", icon: FileText, to: "/applications" },
  { title: "Даалгаврууд", icon: ClipboardList, to: "/tasks" },
  { title: "AI Чатбот", icon: Bot, to: "/chat" },
  { title: "Профайл", icon: User, to: "/profile" },
  { title: "Тохиргоо", icon: Settings, to: "/settings" },
];

const menuRecruiter = [
  { title: "Нүүр хуудас", icon: ChartNoAxesColumnIncreasing, to: "/" },
  { title: "Ажлын байрууд", icon: BriefcaseBusiness, to: "/jobs" },
  { title: "Даалгаврууд", icon: ClipboardList, to: "/tasks" },
  { title: "Ярилцлага", icon: Calendar, to: "/interviews" },
  { title: "Профайл", icon: User, to: "/profile" },
  { title: "Тохиргоо", icon: Settings, to: "/settings" },
];

const menu = computed(() => {
  if (user?.value?.role === "recruiter") return menuRecruiter;
  return menuUser;
});

const currentPage = computed(() => {
  return (
    menu.value.find((item) => item.to === route.path)?.title || "Ажлын орчин"
  );
});

const currentSubtitle = computed(() => {
  if (user.value?.role === "recruiter") {
    return "Ажилд авах урсгалаа хянаж, зарыг удирдаж, ярилцлагыг зохион байгуулна уу.";
  }

  return "Анкетуудаа хянаж, тохирох ажлыг олж, дараагийн алхамдаа бэлэн байгаарай.";
});

const recruiterNeedsCompany = computed(
  () => user.value?.role === "recruiter" && !user.value?.company_id,
);

function handleMenuNavigation(path: string) {
  if (path === "/jobs" && recruiterNeedsCompany.value) {
    toast.warning("Ажлын байр нээхийн өмнө компанийн мэдээллээ нэмнэ үү.");
    router.push("/profile");
    return;
  }

  router.push(path);
}

function openJobsWorkspace() {
  if (recruiterNeedsCompany.value) {
    toast.warning(
      "Ажлын байрны менежерийг нээхийн өмнө компанийн мэдээллээ нэмнэ үү.",
    );
    router.push("/profile");
    return;
  }

  router.push("/jobs");
}

async function handleLogout() {
  await logout();
  await router.push("/login");
}
</script>
<template>
  <div class="min-h-screen bg-background text-foreground">
    <SidebarProvider>
      <Sidebar class="border-r border-border bg-card">
        <SidebarHeader class="border-b border-border px-4 py-5">
          <div class="flex items-center gap-3">
            <div class="h-10 w-10 items-center justify-center rounded-2xl">
              <img src="/icons/logo.svg" alt="" />
            </div>

            <div>
              <h1 class="text-base font-semibold leading-none">MatchHire</h1>
              <p class="mt-1 text-xs text-muted-foreground">
                Ур чадварын үнэлгээний систем
              </p>
            </div>
          </div>
        </SidebarHeader>

        <SidebarContent class="px-1 py-2">
          <SidebarGroup>
            <SidebarMenu class="gap-0">
              <SidebarMenuItem v-for="item in menu" :key="item.title">
                <SidebarMenuButton
                  @click="handleMenuNavigation(item.to)"
                  class="group m-0 flex w-full items-center gap-3 rounded-2xl px-3 py-5 text-sm font-medium text-muted-foreground transition hover:bg-muted hover:text-foreground"
                  :data-active="$route.path === item.to"
                >
                  <component :is="item.icon" class="h-5 w-5 shrink-0" />

                  <span class="truncate">
                    {{ item.title }}
                  </span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroup>
        </SidebarContent>

        <SidebarFooter class="mt-auto border-t px-4 py-4">
          <div class="rounded-2xl border border-border bg-muted/40 px-3 py-3">
            <p class="text-xs font-medium text-muted-foreground">
              Нэвтэрсэн хэрэглэгч
            </p>
            <p class="mt-1 text-sm font-semibold">{{ user?.name }}</p>
            <p class="mt-1 text-xs capitalize text-muted-foreground">
              {{ user?.role === "recruiter" ? "Ажил олгогч" : "Ажил горилогч" }}
            </p>
          </div>

          <Button
            variant="outline"
            class="mt-3 w-full justify-start gap-2"
            :disabled="isLoading"
            @click="handleLogout"
          >
            <LogOut class="h-4 w-4" />
            {{ isLoading ? "Гарж байна..." : "Гарах" }}
          </Button>
        </SidebarFooter>
      </Sidebar>

      <div class="flex min-h-screen flex-1 flex-col bg-background">
        <header
          class="border-b border-border bg-card/80 px-5 py-4 backdrop-blur"
        >
          <div
            class="mx-auto flex w-full max-w-7xl flex-col gap-4 md:flex-row md:items-center md:justify-between"
          >
            <div class="flex items-start gap-3">
              <div class="rounded-2xl border border-border bg-muted/40 p-2.5">
                <PanelLeft class="h-4 w-4" />
              </div>
              <div>
                <p
                  class="text-xs font-medium uppercase tracking-[0.18em] text-muted-foreground"
                >
                  {{
                    user?.role === "recruiter"
                      ? "Ажил олгогчийн ажлын орчин"
                      : "Ажил горилогчийн ажлын орчин"
                  }}
                </p>
                <h1 class="mt-1 text-2xl font-semibold tracking-tight">
                  {{ currentPage }}
                </h1>
                <p class="mt-1 text-sm text-muted-foreground">
                  {{ currentSubtitle }}
                </p>
              </div>
            </div>

            <button
              type="button"
              class="inline-flex items-center gap-2 self-start rounded-full border border-border bg-background px-4 py-2 text-sm font-medium transition hover:bg-muted md:self-auto"
              :class="recruiterNeedsCompany ? 'opacity-70' : ''"
              @click="openJobsWorkspace"
            >
              {{
                user?.role === "recruiter"
                  ? "Ажлын байр удирдах"
                  : "Ажлын байр харах"
              }}
              <ArrowRight class="h-4 w-4" />
            </button>
          </div>
        </header>

        <main class="mx-auto flex w-full max-w-7xl flex-1 flex-col px-5 py-6">
          <NuxtPage />
        </main>
      </div>
    </SidebarProvider>
  </div>
</template>
