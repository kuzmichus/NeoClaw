import type { TFunction } from "i18next"
import { toast } from "sonner"

import type { ChatAttachment } from "@/store/chat"

const CHAT_IMAGE_MIME_TYPES = [
  "image/jpeg",
  "image/png",
  "image/gif",
  "image/webp",
  "image/bmp",
] as const

const CHAT_IMAGE_MIME_TYPE_SET = new Set<string>(CHAT_IMAGE_MIME_TYPES)

const CHAT_IMAGE_MIME_BY_EXTENSION: Record<string, string> = {
  ".jpg": "image/jpeg",
  ".jpeg": "image/jpeg",
  ".png": "image/png",
  ".gif": "image/gif",
  ".webp": "image/webp",
  ".bmp": "image/bmp",
}

// Document attachments are detected primarily by extension because browsers
// report empty or generic MIME types for most code and text files.
const CHAT_DOCUMENT_MIME_BY_EXTENSION: Record<string, string> = {
  ".txt": "text/plain",
  ".log": "text/plain",
  ".ini": "text/plain",
  ".cfg": "text/plain",
  ".conf": "text/plain",
  ".toml": "text/plain",
  ".env": "text/plain",
  ".md": "text/markdown",
  ".markdown": "text/markdown",
  ".csv": "text/csv",
  ".tsv": "text/csv",
  ".json": "application/json",
  ".yaml": "application/x-yaml",
  ".yml": "application/x-yaml",
  ".xml": "application/xml",
  ".html": "text/html",
  ".htm": "text/html",
  ".pdf": "application/pdf",
  ".py": "text/x-python",
  ".js": "text/javascript",
  ".mjs": "text/javascript",
  ".ts": "text/typescript",
  ".tsx": "text/typescript",
  ".jsx": "text/javascript",
  ".go": "text/x-go",
  ".rs": "text/x-rust",
  ".java": "text/x-java",
  ".c": "text/x-c",
  ".h": "text/x-c",
  ".cpp": "text/x-c++",
  ".hpp": "text/x-c++",
  ".cs": "text/x-csharp",
  ".rb": "text/x-ruby",
  ".php": "text/x-php",
  ".swift": "text/x-swift",
  ".kt": "text/x-kotlin",
  ".sh": "text/x-shellscript",
  ".bash": "text/x-shellscript",
  ".zsh": "text/x-shellscript",
  ".sql": "text/x-sql",
  ".diff": "text/x-diff",
  ".patch": "text/x-diff",
}

const CHAT_DOCUMENT_MIME_TYPE_SET = new Set(
  Object.values(CHAT_DOCUMENT_MIME_BY_EXTENSION),
)

export const CHAT_IMAGE_ACCEPT = CHAT_IMAGE_MIME_TYPES.join(",")
export const CHAT_DOCUMENT_ACCEPT = Object.keys(
  CHAT_DOCUMENT_MIME_BY_EXTENSION,
).join(",")
export const CHAT_FILE_ACCEPT = `${CHAT_IMAGE_ACCEPT},${CHAT_DOCUMENT_ACCEPT}`

const MAX_CHAT_IMAGE_SIZE_BYTES = 7 * 1024 * 1024
const MAX_CHAT_IMAGE_SIZE_LABEL = "7 MB"
const MAX_CHAT_DOCUMENT_SIZE_BYTES = 10 * 1024 * 1024
const MAX_CHAT_DOCUMENT_SIZE_LABEL = "10 MB"

function readFileAsDataUrl(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => {
      if (typeof reader.result === "string") {
        resolve(reader.result)
        return
      }
      reject(new Error("Failed to read file"))
    }
    reader.onerror = () =>
      reject(reader.error || new Error("Failed to read file"))
    reader.readAsDataURL(file)
  })
}

function getFileExtension(fileName: string): string {
  const lastDotIndex = fileName.lastIndexOf(".")
  if (lastDotIndex === -1) {
    return ""
  }
  return fileName.slice(lastDotIndex).toLowerCase()
}

type SupportedFileKind = "image" | "document"

interface SupportedFile {
  kind: SupportedFileKind
  mimeType: string
}

function getSupportedImageMimeType(file: File): string | null {
  const normalizedType = file.type.trim().toLowerCase()
  if (normalizedType && CHAT_IMAGE_MIME_TYPE_SET.has(normalizedType)) {
    return normalizedType
  }

  const extension = getFileExtension(file.name)
  return CHAT_IMAGE_MIME_BY_EXTENSION[extension] ?? null
}

function getSupportedDocumentMimeByExtension(file: File): string | null {
  return CHAT_DOCUMENT_MIME_BY_EXTENSION[getFileExtension(file.name)] ?? null
}

function getSupportedFile(file: File): SupportedFile | null {
  const imageMime = getSupportedImageMimeType(file)
  if (imageMime) {
    return { kind: "image", mimeType: imageMime }
  }

  const extension = getFileExtension(file.name)
  const documentMimeByExtension = getSupportedDocumentMimeByExtension(file)
  if (documentMimeByExtension) {
    return { kind: "document", mimeType: documentMimeByExtension }
  }

  // Some browsers still report a precise MIME type for documents without
  // relying on the file name; accept it when it matches our allowlist.
  const normalizedType = file.type.trim().toLowerCase()
  if (
    !extension &&
    normalizedType &&
    CHAT_DOCUMENT_MIME_TYPE_SET.has(normalizedType)
  ) {
    return { kind: "document", mimeType: normalizedType }
  }

  return null
}

function getAttachmentFilename(
  file: File,
  supported: SupportedFile,
  index: number,
): string {
  const trimmedName = file.name.trim()
  if (trimmedName) {
    return trimmedName
  }

  if (supported.kind === "document") {
    return `document-${index + 1}`
  }
  const extension = CHAT_IMAGE_MIME_BY_EXTENSION[supported.mimeType] ?? ".png"
  return `image-${index + 1}${extension}`
}

function getTransferItemFiles(dataTransfer: DataTransfer | null): File[] {
  if (!dataTransfer) {
    return []
  }

  const files = Array.from(dataTransfer.files)
  if (files.length > 0) {
    return files
  }

  return Array.from(dataTransfer.items)
    .filter((item) => item.kind === "file")
    .map((item) => item.getAsFile())
    .filter((file): file is File => file !== null)
}

export function hasFileTransfer(dataTransfer: DataTransfer | null): boolean {
  if (!dataTransfer) {
    return false
  }

  if (dataTransfer.files.length > 0) {
    return true
  }

  return Array.from(dataTransfer.items).some((item) => item.kind === "file")
}

export function getTransferredFiles(dataTransfer: DataTransfer | null) {
  return getTransferItemFiles(dataTransfer)
}

export async function buildChatAttachments(
  files: readonly File[],
  t: TFunction,
): Promise<ChatAttachment[]> {
  const nextAttachments: ChatAttachment[] = []

  for (const [index, file] of files.entries()) {
    const supported = getSupportedFile(file)
    const filename = getAttachmentFilename(
      file,
      supported ?? { kind: "image", mimeType: "" },
      index,
    )

    if (!supported) {
      toast.error(t("chat.invalidAttachment", { name: filename }))
      continue
    }

    const isImage = supported.kind === "image"
    const maxSize = isImage
      ? MAX_CHAT_IMAGE_SIZE_BYTES
      : MAX_CHAT_DOCUMENT_SIZE_BYTES
    if (file.size > maxSize) {
      toast.error(
        t("chat.attachmentTooLarge", {
          name: filename,
          size: isImage
            ? MAX_CHAT_IMAGE_SIZE_LABEL
            : MAX_CHAT_DOCUMENT_SIZE_LABEL,
        }),
      )
      continue
    }

    try {
      nextAttachments.push({
        type: isImage ? "image" : "file",
        filename,
        url: await readFileAsDataUrl(file),
        contentType: supported.mimeType,
      })
    } catch {
      toast.error(t("chat.attachmentReadFailed", { name: filename }))
    }
  }

  return nextAttachments
}
