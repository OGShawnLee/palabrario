import type { IPalabrarioAPI } from "./contract"
import { resolveAdapterKind } from "./config"
import { MockAdapter } from "./mock/mock.adapter"
import { WailsAdapter } from "./live/wails/wails.adapter"

export type { APIMode } from "./config"

function createAPIInstance(): IPalabrarioAPI {
  return resolveAdapterKind() === "wails" ? new WailsAdapter() : new MockAdapter()
}

export const api: IPalabrarioAPI = createAPIInstance()
export * from "./contract"
