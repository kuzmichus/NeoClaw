import { launcherFetch } from "@/api/http"

export interface CronSchedule {
  kind: "at" | "every" | "cron"
  atMs?: number
  everyMs?: number
  expr?: string
  tz?: string
}

export interface CronPayload {
  kind: string
  message: string
  command?: string
  channel?: string
  to?: string
}

export interface CronJobState {
  nextRunAtMs?: number
  lastRunAtMs?: number
  lastStatus?: string
  lastError?: string
}

export interface CronJob {
  id: string
  name: string
  enabled: boolean
  schedule: CronSchedule
  payload: CronPayload
  state: CronJobState
  createdAtMs: number
  updatedAtMs: number
  deleteAfterRun: boolean
}

export interface CronJobInput {
  name: string
  enabled: boolean
  schedule: CronSchedule
  message: string
  command?: string
  channel?: string
  to?: string
}

export interface SessionOption {
  id: string
  title: string
  channel?: string
}

export async function getCronJobs(): Promise<CronJob[]> {
  const res = await launcherFetch("/api/cron/jobs")
  if (!res.ok) {
    throw new Error(`Failed to fetch cron jobs: ${res.status}`)
  }
  return res.json()
}

export async function createCronJob(input: CronJobInput): Promise<CronJob> {
  const res = await launcherFetch("/api/cron/jobs", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  })
  if (!res.ok) {
    const detail = await safeError(res)
    throw new Error(detail || `Failed to create cron job: ${res.status}`)
  }
  return res.json()
}

export async function updateCronJob(
  id: string,
  input: CronJobInput,
): Promise<CronJob> {
  const res = await launcherFetch(`/api/cron/jobs/${encodeURIComponent(id)}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  })
  if (!res.ok) {
    const detail = await safeError(res)
    throw new Error(detail || `Failed to update cron job: ${res.status}`)
  }
  return res.json()
}

export async function deleteCronJob(id: string): Promise<void> {
  const res = await launcherFetch(`/api/cron/jobs/${encodeURIComponent(id)}`, {
    method: "DELETE",
  })
  if (!res.ok) {
    throw new Error(`Failed to delete cron job: ${res.status}`)
  }
}

async function safeError(res: Response): Promise<string | null> {
  try {
    const data = await res.json()
    if (data && typeof data.error === "string") {
      return data.error
    }
  } catch {
    // ignore
  }
  return null
}
