import { createFileRoute } from "@tanstack/react-router"

import { CronJobsPage } from "@/components/cron/cron-jobs-page"

export const Route = createFileRoute("/cron-jobs")({
  component: CronJobsPage,
})
