import { createFileRoute } from "@tanstack/react-router"

import { SessionsPage } from "@/components/chat/sessions-page"

export const Route = createFileRoute("/sessions")({
  component: SessionsPage,
})
