<template>
  <div class="relative">
    <LandingNavbar />

    <main>
      <LandingHero />
      <LandingFeatures />

      <section class="border-t border-slate-200/80 bg-white">
        <div class="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8">
          <div class="grid gap-6 lg:grid-cols-2">
            <Card class="rounded-3xl border-slate-200 shadow-sm">
              <CardHeader class="space-y-3">
                <Badge variant="secondary" class="w-fit rounded-full px-3 py-1">
                  For Applicants
                </Badge>
                <CardTitle class="text-2xl text-slate-950">
                  Showcase skills beyond the CV
                </CardTitle>
                <CardDescription class="text-base leading-7 text-slate-600">
                  Upload your profile, discover relevant roles, and understand
                  how your skills align with each opportunity through transparent
                  AI scoring.
                </CardDescription>
              </CardHeader>
              <CardContent class="space-y-4 text-sm text-slate-600">
                <div class="rounded-2xl bg-slate-50 p-4">
                  Personalized job matching with role-fit signals.
                </div>
                <div class="rounded-2xl bg-slate-50 p-4">
                  Explainable assessment feedback to guide improvement.
                </div>
                <Button as-child class="mt-2">
                  <NuxtLink to="/register">Create applicant account</NuxtLink>
                </Button>
              </CardContent>
            </Card>

            <Card class="rounded-3xl border-slate-200 shadow-sm">
              <CardHeader class="space-y-3">
                <Badge class="w-fit rounded-full px-3 py-1">
                  For Recruiters
                </Badge>
                <CardTitle class="text-2xl text-slate-950">
                  Screen candidates with confidence
                </CardTitle>
                <CardDescription class="text-base leading-7 text-slate-600">
                  Review AI-assisted candidate assessments, compare skill
                  alignment, and move from sourcing to scheduling inside one
                  streamlined workflow.
                </CardDescription>
              </CardHeader>
              <CardContent class="space-y-4 text-sm text-slate-600">
                <div class="rounded-2xl bg-slate-50 p-4">
                  Structured candidate scoring aligned to job requirements.
                </div>
                <div class="rounded-2xl bg-slate-50 p-4">
                  Recruiter-friendly dashboards for pipeline decisions.
                </div>
                <Button as-child variant="outline" class="mt-2">
                  <NuxtLink to="/register">Create recruiter account</NuxtLink>
                </Button>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      <section class="border-t border-slate-200/80 bg-slate-50/80">
        <div class="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8">
          <div class="max-w-2xl space-y-4">
            <Badge variant="outline" class="rounded-full px-3 py-1">
              Simple Workflow
            </Badge>
            <h2 class="text-3xl font-semibold tracking-tight text-slate-950">
              Designed to keep applicants and recruiters moving
            </h2>
            <p class="text-base leading-7 text-slate-600">
              A straightforward process for discovering jobs, submitting
              applications, and reviewing AI-supported skill assessments.
            </p>
          </div>

          <div class="mt-10 grid gap-4 md:grid-cols-2 xl:grid-cols-4">
            <Card
              v-for="(step, index) in workflowSteps"
              :key="step.title"
              class="rounded-3xl border-slate-200 bg-white shadow-sm"
            >
              <CardHeader class="space-y-4">
                <div
                  class="flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br from-blue-600 to-violet-600 text-sm font-semibold text-white"
                >
                  {{ index + 1 }}
                </div>
                <CardTitle class="text-xl text-slate-950">
                  {{ step.title }}
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p class="text-sm leading-6 text-slate-600">
                  {{ step.description }}
                </p>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>
    </main>

    <footer class="border-t border-slate-200 bg-white">
      <div
        class="mx-auto flex max-w-7xl flex-col gap-4 px-4 py-8 text-sm text-slate-500 sm:px-6 lg:flex-row lg:items-center lg:justify-between lg:px-8"
      >
        <div>
          <p class="font-semibold text-slate-900">SkillMatch AI</p>
          <p>AI-powered assessment for modern hiring teams.</p>
        </div>
        <div class="flex flex-wrap items-center gap-4">
          <NuxtLink class="transition hover:text-slate-900" to="/login">
            Login
          </NuxtLink>
          <NuxtLink class="transition hover:text-slate-900" to="/register">
            Get Started
          </NuxtLink>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
const workflowSteps = [
  {
    title: "Create account",
    description:
      "Sign up as an applicant or recruiter and enter the platform with the right role from day one.",
  },
  {
    title: "Browse or post jobs",
    description:
      "Applicants explore opportunities while recruiters publish roles and define hiring requirements.",
  },
  {
    title: "Upload CV or receive applications",
    description:
      "Candidates submit their CVs and recruiters collect incoming applications in one workflow.",
  },
  {
    title: "View AI assessment",
    description:
      "See skill alignment, matching insights, and explainable scoring to support better decisions.",
  },
];
</script>
