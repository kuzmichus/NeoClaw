import { IconCheck, IconCode, IconCopy, IconEye } from "@tabler/icons-react"
import { useState } from "react"
import { useTranslation } from "react-i18next"

import { Button } from "@/components/ui/button"
import { useCopyToClipboard } from "@/hooks/use-copy-to-clipboard"

const CODE_LABEL_FONT_FAMILY =
  'ui-monospace, "SFMono-Regular", Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei UI", "Microsoft YaHei", monospace'

type SvgViewMode = "visual" | "code"

function buildSvgDataUri(svg: string): string {
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg.trim())}`
}

interface MessageSvgBlockProps {
  code: string
}

export function MessageSvgBlock({ code }: MessageSvgBlockProps) {
  const { t } = useTranslation()
  const { copy, isCopied } = useCopyToClipboard()
  const [viewMode, setViewMode] = useState<SvgViewMode>("visual")

  const copyLabel = isCopied ? t("chat.copiedLabel") : t("chat.copyCode")
  const toggleLabel =
    viewMode === "visual" ? t("chat.viewCode") : t("chat.viewVisual")

  return (
    <div
      data-picoclaw-code-block=""
      data-picoclaw-svg=""
      className="not-prose my-4 overflow-hidden rounded-lg border border-[#d0d7de] bg-[#f6f8fa] text-[#24292f] shadow-xs dark:border-[#30363d] dark:bg-[#0d1117] dark:text-[#c9d1d9]"
    >
      <div className="flex items-center justify-between gap-2 border-b border-[#d0d7de] bg-black/[0.03] px-3 py-2 dark:border-[#30363d] dark:bg-white/[0.03]">
        <span
          className="text-[11px] font-medium text-zinc-600 dark:text-zinc-400"
          style={{ fontFamily: CODE_LABEL_FONT_FAMILY }}
        >
          {t("chat.svgLabel")}
        </span>
        <div className="flex items-center gap-1">
          <Button
            type="button"
            variant="ghost"
            size="xs"
            className="h-7 text-zinc-600 hover:bg-zinc-300/70 hover:text-zinc-900 dark:text-zinc-400 dark:hover:bg-zinc-800 dark:hover:text-zinc-100"
            onClick={() => void copy(code)}
            aria-label={copyLabel}
            title={copyLabel}
          >
            {isCopied ? <IconCheck className="text-green-500" /> : <IconCopy />}
            <span className="hidden sm:inline">{copyLabel}</span>
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="xs"
            className="h-7 px-2 text-[11px] text-zinc-600 hover:bg-zinc-300/70 hover:text-zinc-900 dark:text-zinc-400 dark:hover:bg-zinc-800 dark:hover:text-zinc-100"
            onClick={() =>
              setViewMode((mode) => (mode === "visual" ? "code" : "visual"))
            }
            aria-pressed={viewMode === "code"}
            aria-label={toggleLabel}
            title={toggleLabel}
          >
            {viewMode === "visual" ? <IconCode /> : <IconEye />}
            <span className="hidden sm:inline">{toggleLabel}</span>
          </Button>
        </div>
      </div>

      <div className="overflow-x-auto bg-transparent px-4 py-3">
        {viewMode === "visual" ? (
          <img
            src={buildSvgDataUri(code)}
            alt={t("chat.svgLabel")}
            className="mx-auto flex max-h-80 max-w-full object-contain"
          />
        ) : (
          <pre className="m-0 overflow-x-auto font-mono text-[13px] leading-6">
            <code>{code}</code>
          </pre>
        )}
      </div>
    </div>
  )
}
