<script setup lang="ts">
import {
  BadgeCheck,
  Bot,
  BriefcaseBusiness,
  CalendarDays,
  LayoutDashboard,
  Sparkles,
} from "lucide-vue-next";

const features = [
  {
    title: "AI CV Assessment",
    description:
      "Analyze applicant CVs and extract skill signals in a structured, recruiter-friendly view.",
    icon: Sparkles,
  },
  {
    title: "Job Matching",
    description:
      "Compare candidate profiles against job requirements to surface role-fit insights quickly.",
    icon: BriefcaseBusiness,
  },
  {
    title: "Explainable Skill Score",
    description:
      "Show transparent scoring factors so teams can understand why a candidate matches.",
    icon: BadgeCheck,
  },
  {
    title: "Recruiter Dashboard",
    description:
      "Track applications, review assessments, and prioritize the strongest candidates.",
    icon: LayoutDashboard,
  },
  {
    title: "Chatbot Job Search",
    description:
      "Help applicants discover opportunities with a guided conversational job search experience.",
    icon: Bot,
  },
  {
    title: "Interview Scheduling",
    description:
      "Coordinate the next hiring step faster with clear scheduling handoff points.",
    icon: CalendarDays,
  },
];
</script>

<template>
  <section id="features" class="border-t border-slate-200/80 bg-slate-50/80">
    <div class="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8">
      <div class="max-w-2xl space-y-4">
        <Badge variant="outline" class="rounded-full px-3 py-1">
          Platform Features
        </Badge>
        <h2 class="text-3xl font-semibold tracking-tight text-slate-950 sm:text-4xl">
          Purpose-built tools for intelligent hiring workflows
        </h2>
        <p class="text-base leading-7 text-slate-600">
          From candidate discovery to interview planning, the platform keeps
          AI insights understandable and useful for both sides of the hiring
          process.
        </p>
      </div>

      <div class="mt-10 grid gap-5 md:grid-cols-2 xl:grid-cols-3">
        <Card
          v-for="feature in features"
          :key="feature.title"
          class="rounded-3xl border-slate-200 bg-white shadow-sm transition duration-200 hover:-translate-y-1 hover:shadow-lg"
        >
          <CardHeader class="space-y-4">
            <div
              class="flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br from-blue-600/10 to-violet-600/10 text-blue-700"
            >
              <component :is="feature.icon" class="h-6 w-6" />
            </div>
            <CardTitle class="text-xl text-slate-950">
              {{ feature.title }}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p class="text-sm leading-6 text-slate-600">
              {{ feature.description }}
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  </section>
</template>
