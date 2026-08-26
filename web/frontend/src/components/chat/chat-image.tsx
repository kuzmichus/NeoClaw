import type { ComponentProps } from "react"

import { cn } from "@/lib/utils"

type ChatImageProps = Omit<ComponentProps<"img">, "ref"> & {
  node?: unknown
}

export function ChatImage({ className, node, ...props }: ChatImageProps) {
  void node
  return (
    <img
      {...props}
      className={cn(
        "border-border/60 my-2 max-h-80 max-w-full rounded-lg border object-contain",
        className,
      )}
    />
  )
}
