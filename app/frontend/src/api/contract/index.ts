import type { IArenaAPI } from "./arena.contract"
import type { IAuthAPI } from "./auth.contract"
import type { IDictionaryAPI } from "./dictionary.contract"
import type { IStudioAPI } from "./studio.contract"

export type { IArenaAPI } from "./arena.contract"
export type { IAuthAPI } from "./auth.contract"
export type { IDictionaryAPI } from "./dictionary.contract"
export type { IStudioAPI } from "./studio.contract"

export interface IPalabrarioAPI {
  readonly arena: IArenaAPI
  readonly auth: IAuthAPI
  readonly dictionary: IDictionaryAPI
  readonly studio: IStudioAPI
}
